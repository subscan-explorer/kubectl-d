package args

import (
	"github.com/AlecAivazis/survey/v2"
)

var (
	capabilities = map[string]string{
		"SYS_MODULE":      "Load and unload kernel modules.",
		"SYS_RAWIO":       "Perform I/O port operations (iopl(2) and ioperm(2)).",
		"SYS_PACCT":       "Use acct(2), switch process accounting on or off.",
		"SYS_ADMIN":       "Perform a range of system administration operations.",
		"SYS_NICE":        "Raise process nice value (nice(2), setpriority(2)) and change the nice value for arbitrary processes.",
		"SYS_RESOURCE":    "Override resource Limits.",
		"SYS_TIME":        "Set system clock (settimeofday(2), stime(2), adjtimex(2)); set real-time (hardware) clock.",
		"SYS_TTY_CONFIG":  "Use vhangup(2); employ various privileged ioctl(2) operations on virtual terminals.",
		"AUDIT_CONTROL":   "Enable and disable kernel auditing; change auditing filter rules; retrieve auditing status and filtering rules.",
		"MAC_OVERRIDE":    "Allow MAC configuration or state changes. Implemented for the Smack LSM.",
		"MAC_ADMIN":       "Override Mandatory Access Control (MAC). Implemented for the Smack Linux Security Module (LSM).",
		"NET_ADMIN":       "Perform various network-related operations.",
		"SYSLOG":          "Perform privileged syslog(2) operations.",
		"DAC_READ_SEARCH": "Bypass file read permission checks and directory read and execute permission checks.",
		"LINUX_IMMUTABLE": "Set the FS_APPEND_FL and FS_IMMUTABLE_FL i-node flags.",
		"NET_BROADCAST":   "Make socket broadcasts, and listen to multicasts.",
		"IPC_LOCK":        "Lock memory (mlock(2), mlockall(2), mmap(2), shmctl(2)).",
		"IPC_OWNER":       "Bypass permission checks for operations on System V IPC objects.",
		"SYS_PTRACE":      "Trace arbitrary processes using ptrace(2).",
		"SYS_BOOT":        "Use reboot(2) and kexec_load(2), reboot and load a new kernel for later execution.",
		"LEASE":           "Establish leases on arbitrary files (see fcntl(2)).",
		"WAKE_ALARM":      "Trigger something that will wake up the system.",
		"BLOCK_SUSPEND":   "Employ features that can block system suspend.",
	}
)

type CapabilitiesOptions struct {
	question []*survey.Question
}

func NewCapabilitiesOptions(messages ...string) *CapabilitiesOptions {
	po := &CapabilitiesOptions{}
	var msg = "Which capabilities should be add to the ephemeral container?"
	if len(messages) != 0 && messages[0] != "" {
		msg = messages[0]
	}
	var keys []string
	for key := range capabilities {
		keys = append(keys, key)
	}
	po.question = []*survey.Question{
		{
			Name: "capabilities",
			Prompt: &survey.MultiSelect{
				Message: msg,
				Options: keys,
				Help:    "choose capabilities",
				Description: func(value string, index int) string {
					return capabilities[value]
				},
				Default: "SYS_PTRACE",
			},
			Validate: survey.Required,
		},
	}
	return po
}

func (p *CapabilitiesOptions) Ask() (result interface{}, err error) {
	var capabilities []string
	err = survey.Ask(p.question, &capabilities)
	if err != nil {
		return nil, err
	}
	return capabilities, nil
}
