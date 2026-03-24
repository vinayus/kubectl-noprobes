package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	namespace := flag.String("n", "default", "namespace to check (use \"\" for all namespaces)")
	flag.Parse()

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}

	pods, err := clientset.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing pods: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAMESPACE\tPOD\tCONTAINER\tMISSING")
	found := false

	for _, pod := range pods.Items {
		for _, c := range pod.Spec.Containers {
			missing := missingProbes(c)
			if missing != "" {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", pod.Namespace, pod.Name, c.Name, missing)
				found = true
			}
		}
	}

	w.Flush()

	if !found {
		fmt.Println("All containers have liveness and readiness probes configured.")
	}
}

func missingProbes(c corev1.Container) string {
	noLiveness := c.LivenessProbe == nil
	noReadiness := c.ReadinessProbe == nil

	switch {
	case noLiveness && noReadiness:
		return "liveness, readiness"
	case noLiveness:
		return "liveness"
	case noReadiness:
		return "readiness"
	default:
		return ""
	}
}
