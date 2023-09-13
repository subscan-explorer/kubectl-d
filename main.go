package main

import (
	"flag"
	"github.com/spf13/cast"
	"k8s.io/client-go/util/homedir"
	"kubectl-debug/pkg"
	"kubectl-debug/pkg/args"
	"path/filepath"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Params struct {
	ns            string
	pod           string
	containerName string
	capabilities  arrayFlags
	image         string
	kubeconfig    string
}

var (
	param = new(Params)
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func interface2String(val interface{}, err error) (string, error) {
	return cast.ToString(val), err
}

func interface2StringArray(val interface{}, err error) ([]string, error) {
	return val.([]string), err
}

func main() {
	flag.StringVar(&param.ns, "n", "default", "namespace")
	flag.StringVar(&param.pod, "p", "", "pod name")
	flag.StringVar(&param.containerName, "c", "", "container name, if not specified, the first container will be selected")
	flag.StringVar(&param.image, "image", "nicolaka/netshoot", "use image to create ephemeral container")
	flag.Var(&param.capabilities, "capabilities", "add capabilities to the container, default is SYS_PTRACE")
	flag.StringVar(&param.kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kubeconfig file path")
	flag.Parse()
	if len(param.capabilities) == 0 {
		param.capabilities = append(param.capabilities, "SYS_PTRACE")
	}
	client := pkg.InitK8sClient(param.kubeconfig)
	if client == nil {
		return
	}

	if param.pod == "" {
		var err error
		param.ns, err = interface2String(args.NewNsOptions(client).Ask())
		Panic(err)
		param.pod, err = interface2String(args.NewPodOptions(param.ns, client).Ask())
		Panic(err)
		param.containerName, err = interface2String(args.NewContainersOptions(param.ns, client, cast.ToString(param.pod)).Ask())
		Panic(err)
		param.capabilities, err = interface2StringArray(args.NewCapabilitiesOptions().Ask())
		Panic(err)
		param.image, err = interface2String(args.NewImageOptions().Ask())
		Panic(err)
	}
	Panic(pkg.NewEphemeralContainer(param.pod, param.containerName, param.ns, param.image, param.capabilities).Do(client))
}
