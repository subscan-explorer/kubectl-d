package args

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type PodOptions struct {
	Ns        string
	question  []*survey.Question
	pods      []string
	k8sClient kubernetes.Interface
}

func NewPodOptions(ns string, k8sClient kubernetes.Interface, messages ...string) *PodOptions {
	po := &PodOptions{
		Ns:        ns,
		k8sClient: k8sClient,
	}
	var msg = "Which pod should be create ephemeral container?"
	if len(messages) != 0 && messages[0] != "" {
		msg = messages[0]
	}
	po.question = []*survey.Question{
		{
			Name: "pod",
			Prompt: &survey.Input{
				Message: msg,
				Suggest: po.suggest,
			},
			Validate: survey.Required,
		},
	}
	return po
}

func (p *PodOptions) suggest(toComplete string) []string {
	if len(p.pods) == 0 {
		pods, err := p.k8sClient.CoreV1().Pods(p.Ns).List(
			context.Background(),
			metav1.ListOptions{},
		)
		if err != nil {
			panic(err)
		}
		for _, pod := range pods.Items {
			p.pods = append(p.pods, pod.Name)
		}
	}
	if toComplete == "" {
		return p.pods
	}
	var pods []string
	for _, pod := range p.pods {
		if strings.HasPrefix(pod, toComplete) {
			pods = append(pods, pod)
		}
	}
	return pods
}

func (p *PodOptions) Ask() (result interface{}, err error) {
	var pod string
	err = survey.Ask(p.question, &pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}
