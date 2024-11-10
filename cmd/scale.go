package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var replicas int32

var scaleCmd = &cobra.Command{
	Use:   "scale [service-name]",
	Short: "Scale a service in Minikube",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]

		deployment, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting deployment: %v", err)
		}

		deployment.Spec.Replicas = &replicas

		_, err = clientset.AppsV1().Deployments("default").Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error updating deployment: %v", err)
		}

		fmt.Printf("Service %s scaled to %d replicas\n", serviceName, replicas)
		return nil
	},
}

func init() {
	scaleCmd.Flags().Int32VarP(&replicas, "replicas", "r", 1, "Number of replicas")
	rootCmd.AddCommand(scaleCmd)
}
