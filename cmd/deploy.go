package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [service-name]",
	Short: "Deploy a service on Minikube",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceName,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": serviceName,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": serviceName,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  serviceName,
								Image: "nginx:latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 80,
									},
								},
							},
						},
					},
				},
			},
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceName,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": serviceName,
				},
				Ports: []corev1.ServicePort{
					{
						Port:       80,
						TargetPort: intstr.FromInt(80),
					},
				},
				Type: corev1.ServiceTypeClusterIP,
			},
		}

		_, err := clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating deployment: %v", err)
		}

		_, err = clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating service: %v", err)
		}

		fmt.Printf("Service %s successfully deployed\n", serviceName)
		return nil
	},
}

func int32Ptr(i int32) *int32 {
	return &i
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
