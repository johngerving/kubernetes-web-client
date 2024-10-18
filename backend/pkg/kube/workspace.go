package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *KubeClient) ListPods(ctx context.Context) ([]v1.Pod, error) {
	pods, err := k.clientset.CoreV1().Pods(k.Namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
