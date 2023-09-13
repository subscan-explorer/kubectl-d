`kubectl-debug` is a tool that simplifies the use of `kubectl debug`. With just one command, you can enter the container without having to look up information such as the name of the Pod, the name of the container, and the namespace. In addition, it provides configuration for `securityContext`.

## Installation

- Use the following command to install:

```
go install github.com/subscan-explorer/kubectl-debug@latest
```

- Download the binary file from the [release](https://github.com/subscan-explorer/kubectl-debug/releases) page and add it to the `PATH` environment variable.

## Usage

- Run `kubectl-debug` in the terminal to enter interactive command line mode. Follow the prompts to operate.
[![asciicast](https://asciinema.org/a/607671.svg)](https://asciinema.org/a/607671)
- Use the following command directly to create a temporary container and enter it:
```
kubectl-debug -n <namespace> <pod-name> -c <container-name> -capabilities <capabilities> -capabilities <capabilities> -image <debug-image>
```
[![asciicast](https://asciinema.org/a/607672.svg)](https://asciinema.org/a/607672)
