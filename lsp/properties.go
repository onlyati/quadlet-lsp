package lsp

type propertyMapItem struct {
	label      string
	hover      []string
	parameters []string
}

func propertiesMap() map[string][]propertyMapItem {
	return map[string][]propertyMapItem{
		"Container": {
			{
				label: "AddCapability",
				hover: []string{
					"Add these capabilities, in addition to the default Podman capability set, to the container.",
					"This is a space separated list of capabilities. This key can be listed multiple times.",
					"",
					"For example:",
					"```systemd",
					"AddCapability=CAP_DAC_OVERRIDE CAP_IPC_OWNER",
					"```",
				},
			},
			{
				label: "AddDevice",
				hover: []string{
					"Adds a device node from the host into the container. The format of this is HOST-DEVICE[:CONTAINER-DEVICE][:PERMISSIONS], where HOST-DEVICE is the path of the device node on the host, CONTAINER-DEVICE is the path of the device node in the container, and PERMISSIONS is a list of permissions combining ‘r’ for read, ‘w’ for write, and ‘m’ for mknod(2). The - prefix tells Quadlet to add the device only if it exists on the host.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				label: "AddHost",
				hover: []string{
					"Add host-to-IP mapping to /etc/hosts. The format is hostname:ip.",
					"",
					"Equivalent to the Podman --add-host option. This key can be listed multiple times.",
				},
			},
			{
				label: "Annotation",
				hover: []string{
					"Set one or more OCI annotations on the container. The format is a list of key=value items, similar to Environment.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				label: "AutoUpdate",
				hover: []string{
					"Indicates whether the container will be auto-updated ([podman-auto-update(1)](https://docs.podman.io/en/latest/markdown/podman-auto-update.1.html)). The following values are supported:",
					"",
					"- `registry`: Requires a fully-qualified image reference (e.g., quay.io/podman/stable:latest) to be used to create the container. This enforcement is necessary to know which image to actually check and pull. If an image ID was used, Podman does not know which image to check/pull anymore.",
					"- `local`: Tells Podman to compare the image a container is using to the image with its raw name in local storage. If an image is updated locally, Podman simply restarts the systemd unit executing the container.",
				},
				parameters: []string{"registry", "local"},
			},
			{
				label: "CgroupsMode",
				hover: []string{
					"The cgroups mode of the Podman container. Equivalent to the Podman `--cgroups` option.",
					"",
					"By default, the cgroups mode of the container created by Quadlet is `split`, which differs from the default (`enabled`) used by the Podman CLI.",
					"",
					"If the container joins a pod (i.e. `Pod=` is specified), you may want to change this to `no-conmon` or `enabled` so that pod level cgroup resource limits can take effect.",
				},
			},
			{
				label: "ContainerName",
				hover: []string{
					"The (optional) name of the Podman container. If this is not specified, the default value of `systemd-%N` is used, which is the same as the service name but with a `systemd-` prefix to avoid conflicts with user-managed containers.",
				},
			},
			{
				label: "ContainersConfModule",
				hover: []string{},
			},
			{
				label:      "DNS",
				hover:      []string{},
				parameters: []string{"1.1.1.1", "8.8.8.8"},
			},
			{
				label: "DNSOption",
				hover: []string{},
			},
			{
				label: "DNSSearch",
				hover: []string{},
			},
			{
				label: "DropCapability",
				hover: []string{},
			},
			{
				label: "Entrypoint",
				hover: []string{},
			},
			{
				label: "Environment",
				hover: []string{},
			},
			{
				label: "EnvironmentFile",
				hover: []string{},
			},
			{
				label: "EnvironmentHost",
				hover: []string{},
			},
			{
				label: "Exec",
				hover: []string{},
			},
			{
				label: "ExposeHostPort",
				hover: []string{},
			},
			{
				label: "GIDMap",
				hover: []string{},
			},
			{
				label: "GlobalArgs",
				hover: []string{},
			},
			{
				label: "Group",
				hover: []string{},
			},
			{
				label: "GroupAdd",
				hover: []string{},
			},
			{
				label: "HealthCmd",
				hover: []string{},
			},
			{
				label: "HealthInterval",
				hover: []string{},
			},
			{
				label: "HealthLogDestination",
				hover: []string{},
			},
			{
				label: "HealthMaxLogCount",
				hover: []string{},
			},
			{
				label: "HealthMaxLogSize",
				hover: []string{},
			},
			{
				label: "HealthOnFailure",
				hover: []string{},
			},
			{
				label: "HealthRetries",
				hover: []string{},
			},
			{
				label: "HealthStartPeriod",
				hover: []string{},
			},
			{
				label: "HealthStartupCmd",
				hover: []string{},
			},
			{
				label: "HealthStartupInterval",
				hover: []string{},
			},
			{
				label: "HealthStartupRetries",
				hover: []string{},
			},
			{
				label: "HealthStartupSuccess",
				hover: []string{},
			},
			{
				label: "HealthStartupTimeout",
				hover: []string{},
			},
			{
				label: "HealthTimeout",
				hover: []string{},
			},
			{
				label: "HostName",
				hover: []string{},
			},
			{
				label: "Image",
				hover: []string{},
			},
			{
				label: "IP",
				hover: []string{},
			},
			{
				label: "IP6",
				hover: []string{},
			},
			{
				label: "Label",
				hover: []string{},
			},
			{
				label: "LogDriver",
				hover: []string{},
			},
			{
				label: "LogOpt",
				hover: []string{},
			},
			{
				label: "Mask",
				hover: []string{},
			},
			{
				label: "Memory",
				hover: []string{},
			},
			{
				label: "Mount",
				hover: []string{},
			},
			{
				label: "Network",
				hover: []string{},
			},
			{
				label: "NetworkAlias",
				hover: []string{},
			},
			{
				label: "NoNewPrivileges",
				hover: []string{},
			},
			{
				label: "Notify",
				hover: []string{},
			},
			{
				label: "PidsLimit",
				hover: []string{},
			},
			{
				label: "Pod",
				hover: []string{},
			},
			{
				label: "PodmanArgs",
				hover: []string{},
			},
			{
				label: "PublishPort",
				hover: []string{},
			},
			{
				label: "Pull",
				hover: []string{},
			},
			{
				label: "ReadOnly",
				hover: []string{},
			},
			{
				label: "ReadOnlyTmpfs",
				hover: []string{},
			},
			{
				label: "ReloadCmd",
				hover: []string{},
			},
			{
				label: "ReloadSignal",
				hover: []string{},
			},
			{
				label: "Retry",
				hover: []string{},
			},
			{
				label: "RetryDelay",
				hover: []string{},
			},
			{
				label: "Rootfs",
				hover: []string{},
			},
			{
				label: "RunInit",
				hover: []string{},
			},
			{
				label: "SeccompProfile",
				hover: []string{},
			},
			{
				label: "Secret",
				hover: []string{},
			},
			{
				label: "SecurityLabelDisable",
				hover: []string{},
			},
			{
				label: "SecurityLabelFileType",
				hover: []string{},
			},
			{
				label: "SecurityLabelLevel",
				hover: []string{},
			},
			{
				label: "SecurityLabelNested",
				hover: []string{},
			},
			{
				label: "SecurityLabelType",
				hover: []string{},
			},
			{
				label: "ShmSize",
				hover: []string{},
			},
			{
				label: "StartWithPod",
				hover: []string{},
			},
			{
				label: "StopSignal",
				hover: []string{},
			},
			{
				label: "StopTimeout",
				hover: []string{},
			},
			{
				label: "SubGIDMap",
				hover: []string{},
			},
			{
				label: "SubUIDMap",
				hover: []string{},
			},
			{
				label: "Sysctl",
				hover: []string{},
			},
			{
				label: "Timezone",
				hover: []string{},
			},
			{
				label: "Tmpfs",
				hover: []string{},
			},
			{
				label: "UIDMap",
				hover: []string{},
			},
			{
				label: "Ulimit",
				hover: []string{},
			},
			{
				label: "Unmask",
				hover: []string{},
			},
			{
				label: "User",
				hover: []string{},
			},
			{
				label: "UserNS",
				hover: []string{},
			},
			{
				label: "Volume",
				hover: []string{},
			},
			{
				label: "WorkingDir",
				hover: []string{},
			},
		},
		"Pod": {
			{
				label: "AddHost",
				hover: []string{},
			},
			{
				label: "ContainersConfModule",
				hover: []string{},
			},
			{
				label:      "DNS",
				hover:      []string{},
				parameters: []string{"1.1.1.1", "8.8.8.8"},
			},
			{
				label: "DNSOption",
				hover: []string{},
			},
			{
				label: "DNSSearch",
				hover: []string{},
			},
			{
				label: "ExitPolicy",
				hover: []string{},
			},
			{
				label: "GIDMap",
				hover: []string{},
			},
			{
				label: "GlobalArgs",
				hover: []string{},
			},
			{
				label: "HostName",
				hover: []string{},
			},
			{
				label: "IP",
				hover: []string{},
			},
			{
				label: "IP6",
				hover: []string{},
			},
			{
				label: "Label",
				hover: []string{},
			},
			{
				label: "Network",
				hover: []string{},
			},
			{
				label: "NetworkAlias",
				hover: []string{},
			},
			{
				label: "PodmanArgs",
				hover: []string{},
			},
			{
				label: "PodName",
				hover: []string{},
			},
			{
				label: "PublishPort",
				hover: []string{},
			},
			{
				label: "ServiceName",
				hover: []string{},
			},
			{
				label: "ShmSize",
				hover: []string{},
			},
			{
				label: "SubGIDMap",
				hover: []string{},
			},
			{
				label: "SubUIDMap",
				hover: []string{},
			},
			{
				label: "UIDMap",
				hover: []string{},
			},
			{
				label: "UserNS",
				hover: []string{},
			},
			{
				label: "Volume",
				hover: []string{},
			},
		},
		"Kube": {
			{
				label:      "AutoUpdate",
				hover:      []string{},
				parameters: []string{"registry", "local"},
			},
			{
				label: "ConfigMap",
				hover: []string{},
			},
			{
				label: "ContainersConfModule",
				hover: []string{},
			},
			{
				label: "ExitCodePropagation",
				hover: []string{},
			},
			{
				label: "GlobalArgs",
				hover: []string{},
			},
			{
				label: "KubeDownForce",
				hover: []string{},
			},
			{
				label: "LogDriver",
				hover: []string{},
			},
			{
				label: "Network",
				hover: []string{},
			},
			{
				label: "PodmanArgs",
				hover: []string{},
			},
			{
				label: "PublishPort",
				hover: []string{},
			},
			{
				label: "SetWorkingDirectory",
				hover: []string{},
			},
			{
				label: "UserNS",
				hover: []string{},
			},
			{
				label: "Yaml",
				hover: []string{},
			},
		},
		"Network": {
			{
				label: "ContainersConfModule",
				hover: []string{},
			},
			{
				label: "DisableDNS",
				hover: []string{},
			},
			{
				label:      "DNS",
				hover:      []string{},
				parameters: []string{"1.1.1.1", "8.8.8.8"},
			},
			{
				label: "Driver",
				hover: []string{},
			},
			{
				label: "Gateway",
				hover: []string{},
			},
			{
				label: "GlobalArgs",
				hover: []string{},
			},
			{
				label: "InterfaceName",
				hover: []string{},
			},
			{
				label: "Internal",
				hover: []string{},
			},
			{
				label: "IPAMDriver",
				hover: []string{},
			},
			{
				label: "IPRange",
				hover: []string{},
			},
			{
				label: "IPv6",
				hover: []string{},
			},
			{
				label: "Label",
				hover: []string{},
			},
			{
				label: "NetworkDeleteOnStop",
				hover: []string{},
			},
			{
				label: "NetworkName",
				hover: []string{},
			},
			{
				label: "Options",
				hover: []string{},
			},
			{
				label: "PodmanArgs",
				hover: []string{},
			},
			{
				label: "Subnet",
				hover: []string{},
			},
		},
		"Volume": {
			{
				label: "ContainersConfModule",
				hover: []string{},
			},
			{
				label: "Copy",
				hover: []string{},
			},
			{
				label: "Device",
				hover: []string{},
			},
			{
				label: "Driver",
				hover: []string{},
			},
			{
				label: "GlobalArgs",
				hover: []string{},
			},
			{
				label: "Group",
				hover: []string{},
			},
			{
				label: "Image",
				hover: []string{},
			},
			{
				label: "Label",
				hover: []string{},
			},
			{
				label: "Options",
				hover: []string{},
			},
			{
				label: "PodmanArgs",
				hover: []string{},
			},
			{
				label: "Type",
				hover: []string{},
			},
			{
				label: "User",
				hover: []string{},
			},
			{
				label: "VolumeName",
				hover: []string{},
			},
		},
	}
}
