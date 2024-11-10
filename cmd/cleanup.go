package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up resources",
	RunE: func(_ *cobra.Command, _ []string) error {
		// Remove deployments
		err := clientset.AppsV1().Deployments("default").DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error removing deployments: %v", err)
		}

		// Remove services
		services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error listing services: %v", err)
		}

		for _, svc := range services.Items {
			if svc.Name != "kubernetes" {
				err := clientset.CoreV1().Services("default").Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("error removing service %s: %v", svc.Name, err)
				}
			}
		}

		fmt.Println("Environment cleaned successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}
