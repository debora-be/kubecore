package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of services in Minikube",
	RunE: func(cmd *cobra.Command, args []string) error {
		deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error listing deployments: %v", err)
		}

		fmt.Println("\nDeployment Status:")
		fmt.Printf("%-20s %-10s %-10s %-10s\n", "NAME", "DESIRED", "CURRENT", "AVAILABLE")
		for _, d := range deployments.Items {
			fmt.Printf("%-20s %-10d %-10d %-10d\n",
				d.Name,
				*d.Spec.Replicas,
				d.Status.ReadyReplicas,
				d.Status.AvailableReplicas,
			)
		}

		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error listing pods: %v", err)
		}

		fmt.Println("\nPod Status:")
		fmt.Printf("%-30s %-15s %-15s\n", "NAME", "STATUS", "IP")
		for _, pod := range pods.Items {
			fmt.Printf("%-30s %-15s %-15s\n",
				pod.Name,
				string(pod.Status.Phase),
				pod.Status.PodIP,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}