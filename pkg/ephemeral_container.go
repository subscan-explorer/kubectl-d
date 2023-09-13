package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/exec"
	"strings"
	"time"
)

type EphemeralContainer struct {
	podName       string
	containerName string
	ns            string
	image         string
	capabilities  []string
}

func NewEphemeralContainer(podName, containerName, ns, image string, capabilities []string) *EphemeralContainer {
	return &EphemeralContainer{
		podName:       podName,
		containerName: containerName,
		ns:            ns,
		image:         image,
		capabilities:  capabilities,
	}
}

// Do
func (p *EphemeralContainer) Do(client kubernetes.Interface) error {
	// 获取当前pod中是否存在已知的ephemeral container
	// 如果存在，则询问是否进入该container
	// 如果不存在，则创建一个新的ephemeral container
	pod, err := client.CoreV1().Pods(p.ns).Get(context.Background(), p.podName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "get pod error")
	}
	// 判断是否存在已知的ephemeral container
	if len(pod.Spec.EphemeralContainers) != 0 {
		var ephemeralContainers []string
		for _, ephemeralContainer := range pod.Spec.EphemeralContainers {
			ephemeralContainers = append(ephemeralContainers, ephemeralContainer.Name)
		}
		// 询问是否进入已知的ephemeral container
		var IsEnterEphemeralContainer string
		err = survey.AskOne(&survey.Input{
			Message: "The ephemeral container already exists, is it entering?(Y/N)",
			Default: "Y",
		}, &IsEnterEphemeralContainer)
		if strings.ToUpper(IsEnterEphemeralContainer) == "Y" {
			// 进入已知的ephemeral container
			var ephemeralContainerName string
			err = survey.AskOne(&survey.Select{
				Message: "Which ephemeral container should be enter?",
				Options: ephemeralContainers,
			}, &ephemeralContainerName)
			if err != nil {
				return errors.Wrap(err, "select ephemeral container error")
			}
			return p.startShell(ephemeralContainerName)
		}
	}
	podJs, _ := json.Marshal(pod)
	// 为 pod 添加 EphemeralContainers
	var capabilities []v1.Capability
	for _, v := range p.capabilities {
		capabilities = append(capabilities, v1.Capability(v))
	}
	var ephemeralContainerName = fmt.Sprintf("debug-%s", time.Now().Format("20060102150405"))
	ec := v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:  ephemeralContainerName,
			Image: p.image,
			Command: []string{
				"sh",
				"-c",
				"clear; (bash || ash || sh)",
			},
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "Always",
			SecurityContext: &v1.SecurityContext{
				Capabilities: &v1.Capabilities{
					Add: capabilities,
				},
			},
			Stdin:     true,
			StdinOnce: true,
			TTY:       true,
		},
	}
	if p.containerName != "" {
		ec.TargetContainerName = p.containerName
	}
	pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, ec)
	debugPodJs, _ := json.Marshal(pod)
	data, err := strategicpatch.CreateTwoWayMergePatch(podJs, debugPodJs, pod)
	if err != nil {
		panic(err)
	}
	pods := client.CoreV1().Pods(p.ns)
	_, err = pods.Patch(context.Background(), pod.Name, types.StrategicMergePatchType,
		data, metav1.PatchOptions{}, "ephemeralcontainers")
	if err != nil {
		return errors.Wrap(err, "patch pod error")
	}
	// 等待ephemeral container创建完成
	for {
		fmt.Println("waiting for ephemeral container create...")
		pod, err = pods.Get(context.Background(), p.podName, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "get pod error")
		}
		if len(pod.Status.EphemeralContainerStatuses) != 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return p.startShell(ephemeralContainerName)
}

func (p *EphemeralContainer) startShell(ephemeralContainerName string) error {
	cmd := exec.Command("kubectl", "exec", "-it", "-n", p.ns, p.podName, "-c",
		ephemeralContainerName, "--", "sh", "-c", "(bash || ash || sh)")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.String())
	return cmd.Run()
}
