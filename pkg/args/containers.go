package args

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ContainersOptions struct {
	containerName string
}

func NewContainersOptions(ns string, k8sClient kubernetes.Interface, podName string, messages ...string) *ContainersOptions {
	po := &ContainersOptions{}
	var msg = "Which container should be create ephemeral container?"
	if len(messages) != 0 && messages[0] != "" {
		msg = messages[0]
	}
	pod, err := k8sClient.CoreV1().Pods(ns).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	if len(pod.Spec.Containers) == 0 {
		panic("no containers in pod")
	}
	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	if len(containers) == 1 {
		po.containerName = containers[0]
		return po
	}

	question := []*survey.Question{
		{
			Name: "container",
			Prompt: &survey.Select{
				Message: msg,
				Options: containers,
				Help:    "choose a container",
			},
			Validate: survey.Required,
		},
	}
	err = survey.Ask(question, &po.containerName)
	if err != nil {
		panic(err)
	}
	return po
}

func (p *ContainersOptions) Ask() (result interface{}, err error) {
	return p.containerName, nil
}
