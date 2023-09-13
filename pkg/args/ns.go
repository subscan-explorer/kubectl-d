package args

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type NsOptions struct {
	client kubernetes.Interface
}

func NewNsOptions(client kubernetes.Interface) *NsOptions {
	return &NsOptions{client: client}
}

func (n *NsOptions) Ask() (string, error) {
	var ns string
	err := survey.Ask([]*survey.Question{
		{
			Name: "ns",
			Prompt: &survey.Input{
				Message: "Which namespace should be create ephemeral container?",
				Suggest: n.suggest,
			},
			Validate: survey.Required,
		},
	}, &ns)
	return ns, err
}

func (n *NsOptions) suggest(toComplete string) []string {
	nsList, err := n.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var ns []string
	for _, n := range nsList.Items {
		if strings.HasPrefix(n.Name, toComplete) {
			ns = append(ns, n.Name)
		}
	}
	return ns
}
