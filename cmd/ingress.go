package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	hostname string
	path     string
)

var ingressCmd = &cobra.Command{
	Use:   "create-ingress [service-name]",
	Short: "Creates an Ingress for the service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createIngress(args[0], clientset)
	},
}

func createIngress(serviceName string, clientset kubernetes.Interface) error {
	pathType := networkingv1.PathTypePrefix
	ingressClassName := "nginx"

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-ingress", serviceName),
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.NetworkingV1().Ingresses("default").Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating ingress: %v", err)
	}

	fmt.Printf("Ingress created successfully for service %s\n", serviceName)
	fmt.Printf("You can access it through: http://%s%s\n", hostname, path)
	return nil
}

func init() {
	ingressCmd.Flags().StringVar(&hostname, "host", "example.local", "Hostname for the Ingress")
	ingressCmd.Flags().StringVar(&path, "path", "/", "Path for the service")
	rootCmd.AddCommand(ingressCmd)
}