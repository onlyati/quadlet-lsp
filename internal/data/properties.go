package data

import "github.com/onlyati/quadlet-lsp/internal/utils"

type PropertyMapItem struct {
	Label      string
	Hover      []string
	Parameters []string
	Macro      string
	MinVersion utils.PodmanVersion
}

type CategoryPropertyItem struct {
	InsertText *string
	Details    *string
}

var (
	CategoryProperty = map[string]CategoryPropertyItem{
		"newContainer": {
			Details: utils.ReturnAsStringPtr("define a new container"),
			InsertText: utils.ReturnAsStringPtr(`[Unit]
Description=${1:description}

[Container]
Image=${2:image}
AutoUpdate=registry
$0

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

[Install]
WantedBy=default.target
`),
		},
		"newVolume": {
			Details: utils.ReturnAsStringPtr("define new volume"),
			InsertText: utils.ReturnAsStringPtr(`[Unit]
Description=${1:description}

[Volume]
$0
`),
		},
		"newPod": {
			Details: utils.ReturnAsStringPtr("define new pod"),
			InsertText: utils.ReturnAsStringPtr(`[Unit]
Description=${1:description}

[Pod]
$0
`),
		},
		"newNetwork": {
			Details: utils.ReturnAsStringPtr("define new network"),
			InsertText: utils.ReturnAsStringPtr(`[Unit]
Description=${1:description}

[Network]
$0
`),
		},
		"newImage": {
			Details: utils.ReturnAsStringPtr("define new image"),
			InsertText: utils.ReturnAsStringPtr(`[Unit]
Description=${1:description}

[Image]
Image=${2:image}
$0
`),
		},
	}

	PropertiesMap = map[string][]PropertyMapItem{
		"Container": {
			{
				Label: "AddCapability",
				Hover: []string{
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
				Label: "AddDevice",
				Hover: []string{
					"Adds a device node from the host into the container. The format of this is HOST-DEVICE[:CONTAINER-DEVICE][:PERMISSIONS], where HOST-DEVICE is the path of the device node on the host, CONTAINER-DEVICE is the path of the device node in the container, and PERMISSIONS is a list of permissions combining 'r' for read, 'w' for write, and 'm' for mknod(2). The - prefix tells Quadlet to add the device only if it exists on the host.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "AddHost",
				Hover: []string{
					"Add host-to-IP mapping to /etc/hosts. The format is hostname:ip.",
					"",
					"Equivalent to the Podman --add-host option. This key can be listed multiple times.",
				},
				Macro: "AddHost=${1:hostname}:${2:ip}\n$0",
			},
			{
				Label: "Annotation",
				Hover: []string{
					"Set one or more OCI annotations on the container. The format is a list of key=value items, similar to Environment.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "Annotation=\"${1:key}=${2:value}\"\n$0",
			},
			{
				Label: "AutoUpdate",
				Hover: []string{
					"Indicates whether the container will be auto-updated ([podman-auto-update(1)](https://docs.podman.io/en/latest/markdown/podman-auto-update.1.html)). The following values are supported:",
					"",
					"- `registry`: Requires a fully-qualified image reference (e.g., quay.io/podman/stable:latest) to be used to create the container. This enforcement is necessary to know which image to actually check and pull. If an image ID was used, Podman does not know which image to check/pull anymore.",
					"- `local`: Tells Podman to compare the image a container is using to the image with its raw name in local storage. If an image is updated locally, Podman simply restarts the systemd unit executing the container.",
				},
				Parameters: []string{"registry", "local"},
			},
			{
				Label: "CgroupsMode",
				Hover: []string{
					"The cgroups mode of the Podman container. Equivalent to the Podman `--cgroups` option.",
					"",
					"By default, the cgroups mode of the container created by Quadlet is `split`, which differs from the default (`enabled`) used by the Podman CLI.",
					"",
					"If the container joins a pod (i.e. `Pod=` is specified), you may want to change this to `no-conmon` or `enabled` so that pod level cgroup resource limits can take effect.",
				},
			},
			{
				Label: "ContainerName",
				Hover: []string{
					"The (optional) name of the Podman container. If this is not specified, the default value of `systemd-%N` is used, which is the same as the service name but with a `systemd-` prefix to avoid conflicts with user-managed containers.",
				},
			},
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman --module option.",
					"",
					"This key can be listed multiple times",
				},
			},
			{
				Label: "DNS",
				Hover: []string{
					"Set network-scoped DNS resolver/nameserver for containers in this network.",
					"",
					"This key can be listed multiple times.",
					"",
					"For example:",
					"```systemd",
					"DNS=1.1.1.1",
					"DNS=1.0.0.1",
					"```",
				},
				Parameters: []string{
					"1.1.1.1",
					"1.0.0.1",
					"8.8.8.8",
					"8.8.4.4",
					"9.9.9.9",
					"149.112.112.112",
				},
			},
			{
				Label: "DNSOption",
				Hover: []string{
					"Set custom DNS options.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "DNSSearch",
				Hover: []string{
					"Set custom DNS search domains. Use `DNSSearch=`. to remove the search domain.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "DropCapability",
				Hover: []string{
					"Drop these capabilities from the default podman capability set, or `all` to drop all capabilities.",
					"",
					"This is a space separated list of capabilities. This key can be listed multiple times.",
					"",
					"For example:",
					"```systemd",
					"DropCapability=CAP_DAC_OVERRIDE CAP_IPC_OWNER",
					"```",
				},
			},
			{
				Label: "Entrypoint",
				Hover: []string{
					"Override the default ENTRYPOINT from the image. Equivalent to the Podman --entrypoint option. Specify multi option commands in the form of a JSON string.",
				},
			},
			{
				Label: "Environment",
				Hover: []string{
					"Set an environment variable in the container. This uses the same format as services in systemd and can be listed multiple times.",
					"",
					"For example:",
					"```systemd",
					"Environment=APP_USERNAME=appuser",
					"```",
				},
				Macro: "Environment=\"${1:name}=${2:value}\"\n$0",
			},
			{
				Label: "EnvironmentFile",
				Hover: []string{
					"Use a line-delimited file to set environment variables in the container. The path may be absolute or relative to the location of the unit file. This key may be used multiple times, and the order persists when passed to `podman run`.",
				},
			},
			{
				Label: "EnvironmentHost",
				Hover: []string{
					"Use the host environment inside of the container.",
				},
			},
			{
				Label: "Exec",
				Hover: []string{
					"Additional arguments for the container; this has exactly the same effect as passing more arguments after a `podman run <image> <arguments>` invocation.",
					"",
					"The format is the same as for [systemd command lines](https://www.freedesktop.org/software/systemd/man/systemd.service.html#Command%20lines), However, unlike the usage scenario for similarly-named systemd `ExecStart=` verb which operates on the ambient root filesystem, it is very common for container images to have their own `ENTRYPOINT` or `CMD` metadata which this interacts with.",
					"",
					"The default expectation for many images is that the image will include an `ENTRYPOINT` with a default binary, and this field will add arguments to that entrypoint.",
					"",
					"Another way to describe this is that it works the same way as the [args field in a Kubernetes pod](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell).",
				},
			},
			{
				Label: "ExposeHostPort",
				Hover: []string{
					"Exposes a port, or a range of ports (e.g. `50-59`), from the host to the container. Equivalent to the Podman `--expose` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "GIDMap",
				Hover: []string{
					"Run the container in a new user namespace using the supplied GID mapping. Equivalent to the Podman `--gidmap` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "GlobalArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `run` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Group",
				Hover: []string{
					"The (numeric) GID to run as inside the container. This does not need to match the GID on the host, which can be modified with UserNS, but if that is not specified, this GID is also used on the host.",
				},
			},
			{
				Label: "GroupAdd",
				Hover: []string{
					"Assign additional groups to the primary user running within the container process. Also supports the keep-groups special flag. Equivalent to the Podman --group-add option.",
				},
			},
			{
				Label: "HealthCmd",
				Hover: []string{
					"Set or alter a healthcheck command for a container. A value of none disables existing healthchecks. Equivalent to the Podman `--health-cmd option`.",
				},
			},
			{
				Label: "HealthInterval",
				Hover: []string{
					"Set an interval for the healthchecks. An interval of disable results in no automatic timer setup. Equivalent to the Podman `--health-interval` option.",
				},
			},
			{
				Label: "HealthLogDestination",
				Hover: []string{
					"Set the destination of the HealthCheck log. Directory path, local or events_logger (local use container state file) (Default: local) Equivalent to the Podman `--health-log-destination` option.",
					"",
					"- `local`: (default) HealthCheck logs are stored in overlay containers. (For example: `$runroot/healthcheck.log`)",
					"- `directory`: creates a log file named `<container-ID>-healthcheck.log` with HealthCheck logs in the specified directory.",
					"- `events_logger`: The log will be written with logging mechanism set by events_logger. It also saves the log to a default directory, for performance on a system with a large number of logs.",
				},
			},
			{
				Label: "HealthMaxLogCount",
				Hover: []string{
					"Set maximum number of attempts in the HealthCheck log file. (‘0’ value means an infinite number of attempts in the log file) (Default: 5 attempts) Equivalent to the Podman `--Health-max-log-count` option.",
				},
			},
			{
				Label: "HealthMaxLogSize",
				Hover: []string{
					"Set maximum length in characters of stored HealthCheck log. (“0” value means an infinite log length) (Default: 500 characters) Equivalent to the Podman `--Health-max-log-size` option.",
				},
			},
			{
				Label: "HealthOnFailure",
				Hover: []string{
					"Action to take once the container transitions to an unhealthy state. The “kill” action in combination integrates best with systemd. Once the container turns unhealthy, it gets killed, and systemd restarts the service. Equivalent to the Podman `--health-on-failure` option.",
				},
			},
			{
				Label: "HealthRetries",
				Hover: []string{
					"The number of retries allowed before a healthcheck is considered to be unhealthy. Equivalent to the Podman `--health-retries` option.",
				},
			},
			{
				Label: "HealthStartPeriod",
				Hover: []string{
					"The initialization time needed for a container to bootstrap. Equivalent to the Podman `--health-start-period` option.",
				},
			},
			{
				Label: "HealthStartupCmd",
				Hover: []string{
					"Set a startup healthcheck command for a container. Equivalent to the Podman `--health-startup-cmd` option.",
				},
			},
			{
				Label: "HealthStartupInterval",
				Hover: []string{
					"Set an interval for the startup healthcheck. An interval of disable results in no automatic timer setup. Equivalent to the Podman `--health-startup-interval` option.",
				},
			},
			{
				Label: "HealthStartupRetries",
				Hover: []string{
					"The number of attempts allowed before the startup healthcheck restarts the container. Equivalent to the Podman `--health-startup-retries` option.",
				},
			},
			{
				Label: "HealthStartupSuccess",
				Hover: []string{
					"The number of successful runs required before the startup healthcheck succeeds and the regular healthcheck begins. Equivalent to the Podman `--health-startup-success` option.",
				},
			},
			{
				Label: "HealthStartupTimeout",
				Hover: []string{
					"The maximum time a startup healthcheck command has to complete before it is marked as failed. Equivalent to the Podman `--health-startup-timeout` option.",
				},
			},
			{
				Label: "HealthTimeout",
				Hover: []string{
					"The maximum time allowed to complete the healthcheck before an interval is considered failed. Equivalent to the Podman `--health-timeout` option.",
				},
			},
			{
				Label: "HostName",
				Hover: []string{
					"Sets the host name that is available inside the container. Equivalent to the Podman --hostname option.",
				},
			},
			{
				Label: "Image",
				Hover: []string{
					"The image to run in the container. It is recommended to use a fully qualified image name rather than a short name, both for performance and robustness reasons.",
					"",
					"The format of the name is the same as when passed to `podman pull`. So, it supports using `:tag` or digests to guarantee the specific image version.",
					"",
					"Special cases:",
					"- If the name of the image ends with `.image`, Quadlet will use the image pulled by the corresponding `.image` file, and the generated systemd service contains a dependency on the `$name-image.service` (or the service name set in the .image file). Note that the corresponding `.image` file must exist.",
					"- If the name of the image ends with `.build`, Quadlet will use the image built by the corresponding `.build` file, and the generated systemd service contains a dependency on the `$name-build.service`. Note: the corresponding `.build` file must exist.",
				},
			},
			{
				Label: "IP",
				Hover: []string{
					"Specify a static IPv4 address for the container, for example **10.88.64.128**. Equivalent to the Podman `--ip` option.",
				},
			},
			{
				Label: "IP6",
				Hover: []string{
					"Specify a static IPv6 address for the container, for example **fd46:db93:aa76:ac37::10**. Equivalent to the Podman `--ip6` option.",
				},
			},
			{
				Label: "Label",
				Hover: []string{
					"Set one or more OCI labels on the container. The format is a list of `key=value` items, similar to `Environment`.",
					"",
					"This key can be listed multiple times.",
					"",
					"For example:",
					"```systemd",
					"Label=app=myapp",
					"```",
				},
				Macro: "Label=\"${1:key}=${2:value}\"\n$0",
			},
			{
				Label: "LogDriver",
				Hover: []string{
					"Set the log-driver used by Podman when running the container. Equivalent to the Podman `--log-driver` option.",
				},
			},
			{
				Label: "LogOpt",
				Hover: []string{
					"Set the log-opt (logging options) used by Podman when running the container. Equivalent to the Podman `--log-opt` option. This key can be listed multiple times.",
				},
			},
			{
				Label: "Mask",
				Hover: []string{
					"Specify the paths to mask separated by a colon. `Mask=/path/1:/path/2`. A masked path cannot be accessed inside the container.",
				},
			},
			{
				Label: "Memory",
				Hover: []string{
					"Specify the amount of memory for the container.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "Mount",
				Hover: []string{
					"Attach a filesystem mount to the container. This is equivalent to the Podman `--mount` option, and generally has the form `type=TYPE,TYPE-SPECIFIC-OPTION[,...]`.",
					"",
					"Special cases:",
					"- For `type=volume`, if source ends with `.volume`, the Podman named volume generated by the corresponding `.volume` file is used.",
					"- For `type=image`, if source ends with `.image`, the image generated by the corresponding `.image` file is used.",
					"",
					"In both cases, the generated systemd service will contain a dependency on the service generated for the corresponding unit. Note: the corresponding `.volume` or `.image` file must exist.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Network",
				Hover: []string{
					"Specify a custom network for the container. This has the same format as the `--network` option to `podman run`. For example, use `host` to use the host network in the container, or `none` to not set up networking in the container.",
					"",
					"Special cases:",
					"",
					"- If the `name` of the network ends with `.network`, a Podman network called `systemd-$name` is used, and the generated systemd service contains a dependency on the `$name-network.service`. Such a network can be automatically created by using a `$name.network` Quadlet file. Note: the corresponding `.network` file must exist.",
					"- If the `name` ends with `.container`, the container will reuse the network stack of another container created by `$name.container`. The generated systemd service contains a dependency on `$name.service`. Note: the corresponding `.container` file must exist.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "NetworkAlias",
				Hover: []string{
					"Add a network-scoped alias for the container. This has the same format as the `--network-alias` option to `podman run`. Aliases can be used to group containers together in DNS resolution: for example, setting `NetworkAlias=web` on multiple containers will make a DNS query for `web` resolve to all the containers with that alias.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "NoNewPrivileges",
				Hover: []string{
					"If enabled, this disables the container processes from gaining additional privileges via things like setuid and file capabilities. Defaults to false.",
				},
			},
			{
				Label: "Notify",
				Hover: []string{
					"By default, Podman is run in such a way that the systemd startup notify command is handled by the container runtime. In other words, the service is deemed started when the container runtime starts the child in the container. However, if the container application supports [sd_notify](https://www.freedesktop.org/software/systemd/man/sd_notify.html), then setting `Notify` to true passes the notification details to the container allowing it to notify of startup on its own.",
					"",
					"In addition, setting `Notify` to `healthy` will postpone startup notifications until such time as the container is marked healthy, as determined by Podman healthchecks. Note that this requires setting up a container healthcheck, see the `HealthCmd` option for more.",
					"",
					"Defaults to false.",
				},
			},
			{
				Label: "PidsLimit",
				Hover: []string{
					"Tune the container's pids limit. This is equivalent to the Podman `--pids-limit` option.",
				},
			},
			{
				Label: "Pod",
				Hover: []string{
					"Specify a Quadlet `.pod` unit to link the container to. The value must take the form of `<name>.pod` and the `.pod` unit must exist.",
					"",
					"Quadlet will add all the necessary parameters to link between the container and the pod and between their corresponding services.",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman run` command in the generated file (right before the image name in the command line). It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "PublishPort",
				Hover: []string{
					"Exposes a port, or a range of ports (e.g. `50-59`), from the container to the host. Equivalent to the Podman `--publish` option. The format is similar to the Podman options, which is of the form `ip:hostPort:containerPort`, `ip::containerPort`, `hostPort:containerPort` or `containerPort`, where the number of host and container ports must be the same (in the case of a range).",
					"",
					"If the IP is set to 0.0.0.0 or not set at all, the port is bound on all IPv4 addresses on the host; use [::] for IPv6.",
					"",
					"Note that not listing a host port means that Podman automatically selects one, and it may be different for each invocation of service. This makes that a less useful option. The allocated port can be found with the `podman port` command.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "PublishPort=${1:interface}:${2:exposed}:${3:source}\n$0",
			},
			{
				Label: "Pull",
				Hover: []string{
					"Set the image pull policy. This is equivalent to the Podman `--pull` option",
				},
				Parameters: []string{
					"always",
					"missing",
					"never",
					"newer",
				},
			},
			{
				Label: "ReadOnly",
				Hover: []string{
					"If enabled, makes the image read-only. Defaults to false.",
				},
			},
			{
				Label: "ReadOnlyTmpfs",
				Hover: []string{
					"If ReadOnly is set to `true`, mount a read-write tmpfs on /dev, /dev/shm, /run, /tmp, and /var/tmp. Defaults to false.",
				},
			},
			{
				Label: "ReloadCmd",
				Hover: []string{
					"Add `ExecReload` line to the `Service` that runs ` podman exec` with this command in this container.",
					"",
					"In order to execute the reload run `systemctl reload <Service>`",
					"",
					"Mutually exclusive with `ReloadSignal`",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "ReloadSignal",
				Hover: []string{
					"Add `ExecReload` line to the `Service` that runs `podman kill` with this signal which sends the signal to the main container process.",
					"",
					"In order to execute the reload run `systemctl reload <Service>`",
					"",
					"Mutually exclusive with `ReloadCmd`",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "Retry",
				Hover: []string{
					"Number of times to retry the image pull when a HTTP error occurs. Equivalent to the Podman `--retry` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "RetryDelay",
				Hover: []string{
					"Delay between retries. Equivalent to the Podman `--retry-delay` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "Rootfs",
				Hover: []string{
					"The rootfs to use for the container. Rootfs points to a directory on the system that contains the content to be run within the container. This option conflicts with the `Image` option.",
					"",
					"The format of the rootfs is the same as when passed to `podman run --rootfs`, so it supports overlay mounts as well.",
					"",
					"Note: On SELinux systems, the rootfs needs the correct label, which is by default unconfined_u:object_r:container_file_t:s0.",
				},
			},
			{
				Label: "RunInit",
				Hover: []string{
					"If enabled, the container has a minimal init process inside the container that forwards signals and reaps processes.",
				},
			},
			{
				Label: "SeccompProfile",
				Hover: []string{
					"Set the seccomp profile to use in the container. If unset, the default podman profile is used. Set to either the pathname of a JSON file, or `unconfined` to disable the seccomp filters.",
				},
			},
			{
				Label: "Secret",
				Hover: []string{
					"Use a Podman secret in the container either as a file or an environment variable. This is equivalent to the Podman `--secret` option and generally has the form `secret[,opt=opt ...]`",
				},
				Macro: "Secret=${1:secret},type=${2:type},target=${3:target}\n$0",
			},
			{
				Label: "SecurityLabelDisable",
				Hover: []string{
					"Turn off label separation for the container.",
				},
			},
			{
				Label: "SecurityLabelFileType",
				Hover: []string{
					"Set the label file type for the container files.",
				},
			},
			{
				Label: "SecurityLabelLevel",
				Hover: []string{
					"Set the label process level for the container processes.",
				},
			},
			{
				Label: "SecurityLabelNested",
				Hover: []string{
					"Allow SecurityLabels to function within the container. This allows separation of containers created within the container.",
				},
			},
			{
				Label: "SecurityLabelType",
				Hover: []string{
					"Set the label process type for the container processes.",
				},
			},
			{
				Label: "ShmSize",
				Hover: []string{
					"Size of /dev/shm.",
					"",
					"This is equivalent to the Podman `--shm-size` option and generally has the form `number[unit]`",
				},
			},
			{
				Label: "StartWithPod",
				Hover: []string{
					"Start the container after the associated pod is created. Default to **true**.",
					"",
					"If `true`, container will be started/stopped/restarted alongside the pod.",
					"",
					"If `false`, the container will not be started when the pod starts. The container will be stopped with the pod. Restarting the pod will also restart the container as long as the container was also running before.",
					"",
					"Note, the container can still be started manually or through a target by configuring the `[Install]` section. The pod will be started as needed in any case.",
				},
			},
			{
				Label: "StopSignal",
				Hover: []string{
					"Signal to stop a container. Default is **SIGTERM**.",
					"",
					"This is equivalent to the Podman `--stop-signal` option",
				},
				Parameters: []string{
					"SIGTERM",
					"SIGKILL",
				},
			},
			{
				Label: "StopTimeout",
				Hover: []string{
					"Seconds to wait before forcibly stopping the container.",
					"",
					"Note, this value should be lower than the actual systemd unit timeout to make sure the podman rm command is not killed by systemd.",
					"",
					"This is equivalent to the Podman `--stop-timeout` option",
				},
			},
			{
				Label: "SubGIDMap",
				Hover: []string{
					"Run the container in a new user namespace using the map with name in the /etc/subgid file. Equivalent to the Podman `--subgidname` option.",
				},
			},
			{
				Label: "SubUIDMap",
				Hover: []string{
					"Run the container in a new user namespace using the map with name in the /etc/subuid file. Equivalent to the Podman `--subuidname` option.",
				},
			},
			{
				Label: "Sysctl",
				Hover: []string{
					"Configures namespaced kernel parameters for the container. The format is `Sysctl=name=value`.",
					"",
					"This is a space separated list of kernel parameters. This key can be listed multiple times.",
					"",
					"For example:",
					"```",
					"Sysctl=net.ipv6.conf.all.disable_ipv6=1 net.ipv6.conf.all.use_tempaddr=1",
					"```",
				},
			},
			{
				Label: "Timezone",
				Hover: []string{
					"The timezone to run the container in.",
				},
			},
			{
				Label: "Tmpfs",
				Hover: []string{
					"Mount a tmpfs in the container. This is equivalent to the Podman `--tmpfs` option, and generally has the form `CONTAINER-DIR[:OPTIONS]`.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "UIDMap",
				Hover: []string{
					"Run the container in a new user namespace using the supplied UID mapping. Equivalent to the Podman `--uidmap` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Ulimit",
				Hover: []string{
					"Ulimit options. Sets the ulimits values inside of the container.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Unmask",
				Hover: []string{
					"Specify the paths to unmask separated by a colon. unmask=ALL or /path/1:/path/2, or shell expanded paths (/proc/*):",
					"",
					"If set to `ALL`, Podman will unmask all the paths that are masked or made read-only by default.",
					"",
					"The default masked paths are /proc/acpi, /proc/kcore, /proc/keys, /proc/latency_stats, /proc/sched_debug, /proc/scsi, /proc/timer_list, /proc/timer_stats, /sys/firmware, and /sys/fs/selinux.",
					"",
					"The default paths that are read-only are /proc/asound, /proc/bus, /proc/fs, /proc/irq, /proc/sys, /proc/sysrq-trigger, /sys/fs/cgroup.",
				},
			},
			{
				Label: "User",
				Hover: []string{
					"The (numeric) UID to run as inside the container. This does not need to match the UID on the host, which can be modified with `UserNS`, but if that is not specified, this UID is also used on the host.",
				},
			},
			{
				Label: "UserNS",
				Hover: []string{
					"Set the user namespace mode for the container. This is equivalent to the Podman `--userns` option and generally has the form `MODE[:OPTIONS,...]`.",
				},
				Parameters: []string{
					"auto",
					"host",
					"keep-id",
					"nomap",
				},
			},
			{
				Label: "Volume",
				Hover: []string{
					"Mount a volume in the container. This is equivalent to the Podman `--volume` option, and generally has the form `[[SOURCE-VOLUME|HOST-DIR:]CONTAINER-DIR[:OPTIONS]]`.",
					"",
					"If `SOURCE-VOLUME` starts with `.`, Quadlet resolves the path relative to the location of the unit file.",
					"",
					"Special case:",
					"- If `SOURCE-VOLUME` ends with `.volume`, a Podman named volume called `systemd-$name` is used as the source, and the generated systemd service contains a dependency on the `$name-volume.service`. Note that the corresponding `.volume` file must exist.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "Volume=${1:destination}:${2:source}\n$0",
			},
			{
				Label: "WorkingDir",
				Hover: []string{
					"Working directory inside the container.",
					"",
					"The default working directory for running binaries within a container is the root directory (/). The image developer can set a different default with the WORKDIR instruction. This option overrides the working directory by using the -w option.",
				},
			},
		},
		"Pod": {
			{
				Label: "AddHost",
				Hover: []string{
					"Add  host-to-IP mapping to /etc/hosts. The format is `hostname:ip`.",
					"",
					"Equivalent to the Podman `--add-host` option. This key can be listed multiple times.",
				},
			},
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman `--module` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "DNS",
				Hover: []string{
					"Set network-scoped DNS resolver/nameserver for containers in this pod.",
					"",
					"This key can be listed multiple times.",
				},
				Parameters: []string{
					"1.1.1.1",
					"1.0.0.1",
					"8.8.8.8",
					"8.8.4.4",
					"9.9.9.9",
					"149.112.112.112",
				},
			},
			{
				Label: "DNSOption",
				Hover: []string{
					"Set custom DNS options.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "DNSSearch",
				Hover: []string{
					"Set custom DNS search domains. Use **DNSSearch=.** to remove the search domain.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "ExitPolicy",
				Hover: []string{
					"Set the exit policy of the pod when the last container exits. Default for quadlets is **stop**.",
					"",
					"To keep the pod active, set `ExitPolicy=continue`.",
				},
				Parameters: []string{
					"stop",
					"continue",
				},
				MinVersion: utils.BuildPodmanVersion(5, 6, 0),
			},
			{
				Label: "GIDMap",
				Hover: []string{
					"Create the pod in a new user namespace using the supplied GID mapping. Equivalent to the Podman `--gidmap` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "GlobalArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `pod` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "HostName",
				Hover: []string{
					"Set the pod’s hostname inside all containers.",
					"",
					"The given hostname is also added to the /etc/hosts file using the container’s primary IP address (also see the `--add-host` option).",
					"",
					"Equivalent to the Podman `--hostname` option. This key can be listed multiple times.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "IP",
				Hover: []string{
					"Specify a static IPv4 address for the pod, for example **10.88.64.128**. Equivalent to the Podman `--ip` option.",
				},
			},
			{
				Label: "IP6",
				Hover: []string{
					"Specify a static IPv6 address for the pod, for example **fd46:db93:aa76:ac37::10**. Equivalent to the Podman `--ip6` option.",
				},
			},
			{
				Label: "Label",
				Hover: []string{
					"Set one or more OCI labels on the pod. The format is a list of `key=value` items, similar to `Environment`.",
					"",
					"This key can be listed multiple times.",
				},
				Macro:      "Label=\"${1:key}=${2:value}\"\n$0",
				MinVersion: utils.BuildPodmanVersion(5, 6, 0),
			},
			{
				Label: "Network",
				Hover: []string{
					"Specify a custom network for the pod. This has the same format as the `--network` option to `podman pod create`. For example, use `host` to use the host network in the pod, or `none` to not set up networking in the pod.",
					"",
					"Special case:",
					"- If the `name` of the network ends with `.network`, Quadlet will look for the corresponding `.network` Quadlet unit. If found, Quadlet will use the name of the Network set in the Unit, otherwise, `systemd-$name` is used.",
					"",
					"The generated systemd service contains a dependency on the service unit generated for that `.network` unit. Note: the corresponding `.network` file must exist.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "NetworkAlias",
				Hover: []string{
					"Add a network-scoped alias for the pod. This has the same format as the `--network-alias` option to `podman pod create`. Aliases can be used to group containers together in DNS resolution: for example, setting `NetworkAlias=web` on multiple containers will make a DNS query for `web` resolve to all the containers with that alias.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman pod create` command in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "PodName",
				Hover: []string{
					"The (optional) name of the Podman pod. If this is not specified, the default value is the same name as the unit, but with a `systemd-` prefix, i.e. a `$name.pod` file creates a `systemd-$name` Podman pod to avoid conflicts with user-managed pods.",
					"",
					"Please note that pods and containers cannot have the same name. So, if PodName is set, it must not conflict with any container.",
				},
			},
			{
				Label: "PublishPort",
				Hover: []string{
					"Exposes a port, or a range of ports (e.g. `50-59`), from the pod to the host. Equivalent to the Podman `--publish` option. The format is similar to the Podman options, which is of the form `ip:hostPort:containerPort`, `ip::containerPort`, `hostPort:containerPort` or `containerPort`, where the number of host and container ports must be the same (in the case of a range).",
					"",
					"If the IP is set to 0.0.0.0 or not set at all, the port is bound on all IPv4 addresses on the host; use [::] for IPv6.",
					"",
					"Note that not listing a host port means that Podman automatically selects one, and it may be different for each invocation of service. This makes that a less useful option. The allocated port can be found with the `podman port` command.",
					"",
					"When using `host` networking via `Network=host`, the `PublishPort=` option cannot be used.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "PublishPort=${1:interface}:${2:exposed}:${3:source}\n$0",
			},
			{
				Label: "ServiceName",
				Hover: []string{
					"By default, Quadlet will name the systemd service unit by appending `-pod` to the name of the Quadlet. Setting this key overrides this behavior by instructing Quadlet to use the provided name.",
					"",
					"Note, the name should not include the `.service` file extension",
				},
			},
			{
				Label: "ShmSize",
				Hover: []string{
					"Size of /dev/shm.",
					"",
					"This is equivalent to the Podman `--shm-size` option and generally has the form `number[unit]`",
				},
				MinVersion: utils.BuildPodmanVersion(5, 4, 0),
			},
			{
				Label: "SubGIDMap",
				Hover: []string{
					"Create the pod in a new user namespace using the map with name in the /etc/subgid file. Equivalent to the Podman `--subgidname` option.",
				},
			},
			{
				Label: "SubUIDMap",
				Hover: []string{
					"Create the pod in a new user namespace using the map with name in the /etc/subuid file. Equivalent to the Podman `--subuidname` option.",
				},
			},
			{
				Label: "UIDMap",
				Hover: []string{
					"Create the pod in a new user namespace using the supplied UID mapping. Equivalent to the Podman `--uidmap` option.",
				},
			},
			{
				Label: "UserNS",
				Hover: []string{
					"Set the user namespace mode for the pod. This is equivalent to the Podman `--userns` option and generally has the form `MODE[:OPTIONS,...]`.",
				},
				Parameters: []string{
					"auto",
					"host",
					"keep-id",
					"nomap",
				},
			},
			{
				Label: "Volume",
				Hover: []string{
					"Mount a volume in the pod. This is equivalent to the Podman `--volume` option, and generally has the form `[[SOURCE-VOLUME|HOST-DIR:]CONTAINER-DIR[:OPTIONS]]`.",
					"",
					"If `SOURCE-VOLUME` starts with `.`, Quadlet resolves the path relative to the location of the unit file.",
					"",
					"Special case:",
					"- If `SOURCE-VOLUME` ends with `.volume`, Quadlet will look for the corresponding `.volume` Quadlet unit. If found, Quadlet will use the name of the Volume set in the Unit, otherwise, `systemd-$name` is used. Note: the corresponding `.volume` file must exist.",
					"",
					"The generated systemd service contains a dependency on the service unit generated for that `.volume` unit, or on `$name-volume.service` if the `.volume` unit is not found.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "Volume=${1:destination}:${2:source}\n$0",
			},
		},
		"Kube": {
			{
				Label: "AutoUpdate",
				Hover: []string{
					"Indicates whether containers will be auto-updated ([podman-auto-update(1)](podman-auto-update.1.md)). AutoUpdate can be specified multiple times. The following values are supported:",
					"- `registry`: Requires a fully-qualified image reference (e.g., quay.io/podman/stable:latest) to be used to create the container. This enforcement is necessary to know which images to actually check and pull. If an image ID was used, Podman does not know which image to check/pull anymore.",
					"- `local`: Tells Podman to compare the image a container is using to the image with its raw name in local storage. If an image is updated locally, Podman simply restarts the systemd unit executing the Kubernetes Quadlet.",
					"- `name/(local|registry)`: Tells Podman to perform the `local` or `registry` autoupdate on the specified container name.",
				},
				Parameters: []string{"registry", "local"},
			},
			{
				Label: "ConfigMap",
				Hover: []string{
					"Pass the Kubernetes ConfigMap YAML path to `podman kube play` via the `--configmap` argument. Unlike the `configmap` argument, the value may contain only one path but it may be absolute or relative to the location of the unit file.",
					"",
					"This key may be used multiple times",
				},
			},
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman `--module` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "ExitCodePropagation",
				Hover: []string{
					"Control how the main PID of the systemd service should exit. The following values are supported:",
					"- `all`: exit non-zero if all containers have failed (i.e., exited non-zero)",
					" `any`: exit non-zero if any container has failed",
					"- `none`: exit zero and ignore failed containers",
					"",
					"The current default value is `none`.",
				},
			},
			{
				Label: "GlobalArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `kube` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "KubeDownForce",
				Hover: []string{
					"Remove all resources, including volumes, when calling `podman kube down`. Equivalent to the Podman `--force` option.",
				},
			},
			{
				Label: "LogDriver",
				Hover: []string{
					"Set the log-driver Podman uses when running the container. Equivalent to the Podman `--log-driver` option.",
				},
			},
			{
				Label: "Network",
				Hover: []string{
					"Specify a custom network for the container. This has the same format as the `--network` option to `podman kube play`. For example, use `host` to use the host network in the container, or `none` to not set up networking in the container.",
					"",
					"Special case:",
					"- If the `name` of the network ends with `.network`, a Podman network called `systemd-$name` is used, and the generated systemd service contains a dependency on the `$name-network.service`. Such a network can be automatically created by using a `$name.network` Quadlet file. Note: the corresponding `.network` file must exist.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman kube play` command in the generated file (right before the path to the yaml file in the command line). It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "PublishPort",
				Hover: []string{
					"Exposes a port, or a range of ports (e.g. `50-59`), from the container to the host. Equivalent to the `podman kube play`'s `--publish` option. The format is similar to the Podman options, which is of the form `ip:hostPort:containerPort`, `ip::containerPort`, `hostPort:containerPort` or `containerPort`, where the number of host and container ports must be the same (in the case of a range).",
					"",
					"If the IP is set to 0.0.0.0 or not set at all, the port is bound on all IPv4 addresses on the host; use [::] for IPv6.",
					"",
					"The list of published ports specified in the unit file is merged with the list of ports specified in the Kubernetes YAML file. If the same container port and protocol is specified in both, the entry from the unit file takes precedence",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "PublishPort=${1:interface}:${2:exposed}:${3:source}\n$0",
			},
			{
				Label: "SetWorkingDirectory",
				Hover: []string{
					"Set the `WorkingDirectory` field of the `Service` group of the Systemd service unit file. Used to allow `podman kube play` to correctly resolve relative paths. Supported values are `yaml` and `unit` to set the working directory to that of the YAML or Quadlet Unit file respectively.",
					"",
					"Alternatively, users can explicitly set the `WorkingDirectory` field of the `Service` group in the `.kube` file. Please note that if the `WorkingDirectory` field of the `Service` group is set, Quadlet will not set it even if `SetWorkingDirectory` is set",
				},
			},
			{
				Label: "UserNS",
				Hover: []string{
					"Set the user namespace mode for the container. This is equivalent to the Podman `--userns` option and generally has the form `MODE[:OPTIONS,...]`.",
				},
				Parameters: []string{
					"auto",
					"host",
					"keep-id",
					"nomap",
				},
			},
			{
				Label: "Yaml",
				Hover: []string{
					"The path, absolute or relative to the location of the unit file, to the Kubernetes YAML file to use.",
				},
			},
		},
		"Network": {
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman `--module` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "DisableDNS",
				Hover: []string{
					"If enabled, disables the DNS plugin for this network.",
					"",
					"This is equivalent to the Podman `--disable-dns` option",
				},
			},
			{
				Label: "DNS",
				Hover: []string{
					"Set network-scoped DNS resolver/nameserver for containers in this network.",
					"",
					"This key can be listed multiple times.",
				},
				Parameters: []string{"1.1.1.1", "8.8.8.8"},
			},
			{
				Label: "Driver",
				Hover: []string{
					"Driver to manage the network. Currently `bridge`, `macvlan` and `ipvlan` are supported.",
					"",
					"This is equivalent to the Podman `--driver` option",
				},
			},
			{
				Label: "Gateway",
				Hover: []string{
					"Define a gateway for the subnet. If you want to provide a gateway address, you must also provide a subnet option.",
					"",
					"This is equivalent to the Podman `--gateway` option",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "GlobalArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `network` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "InterfaceName",
				Hover: []string{
					"This option maps the *network_interface* option in the network config, see **podman network inspect**. Depending on the driver, this can have different effects; for `bridge`, it uses the bridge interface name. For `macvlan` and `ipvlan`, it is the parent device on the host. It is the same as `--opt parent=...`.",
					"",
					"This is equivalent to the Podman `--interface-name` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 6, 0),
			},
			{
				Label: "Internal",
				Hover: []string{
					"Restrict external access of this network.",
					"",
					"This is equivalent to the Podman `--internal` option",
				},
			},
			{
				Label: "IPAMDriver",
				Hover: []string{
					"Set the ipam driver (IP Address Management Driver) for the network. Currently `host-local`, `dhcp` and `none` are supported.",
					"",
					"This is equivalent to the Podman `--ipam-driver` option",
				},
			},
			{
				Label: "IPRange",
				Hover: []string{
					"Allocate container IP from a range. The range must be a either a complete subnet in CIDR notation or be in the `<startIP>-<endIP>` syntax which allows for a more flexible range compared to the CIDR subnet. The ip-range option must be used with a subnet option.",
					"",
					"This is equivalent to the Podman `--ip-range` option",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "IPv6",
				Hover: []string{
					"Enable IPv6 (Dual Stack) networking.",
					"",
					"This is equivalent to the Podman `--ipv6` option",
				},
			},
			{
				Label: "Label",
				Hover: []string{
					"Set one or more OCI labels on the network. The format is a list of `key=value` items, similar to `Environment`.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "Label=\"${1:key}=${2:value}\"\n$0",
			},
			{
				Label: "NetworkDeleteOnStop",
				Hover: []string{
					"When set to `true` the network is deleted when the service is stopped",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "NetworkName",
				Hover: []string{
					"The (optional) name of the Podman network. If this is not specified, the default value is the same name as the unit, but with a `systemd-` prefix, i.e. a `$name.network` file creates a `systemd-$name` Podman network to avoid conflicts with user-managed network.",
				},
			},
			{
				Label: "Options",
				Hover: []string{
					"Set driver specific options.",
					"",
					"This is equivalent to the Podman `--opt` option",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman network create` command in the generated file (right before the name of the network in the command line). It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Subnet",
				Hover: []string{
					"The subnet in CIDR notation.",
					"",
					"This is equivalent to the Podman `--subnet` option",
					"",
					"This key can be listed multiple times.",
				},
			},
		},
		"Volume": {
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman `--module` option.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Copy",
				Hover: []string{
					"If enabled, the content of the image located at the mountpoint of the volume is copied into the volume on the first run.",
				},
			},
			{
				Label: "Device",
				Hover: []string{
					"The path of a device which is mounted for the volume.",
				},
			},
			{
				Label: "Driver",
				Hover: []string{
					"Specify the volume driver name. When set to `image`, the `Image` key must also be set.",
					"",
					"This is equivalent to the Podman `--driver` option.",
				},
			},
			{
				Label: "GlobalArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `volume` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Group",
				Hover: []string{
					"The host (numeric) GID, or group name to use as the group for the volume",
				},
			},
			{
				Label: "Image",
				Hover: []string{
					"Specifies the image the volume is based on when `Driver` is set to the `image`. It is recommended to use a fully qualified image name rather than a short name, both for performance and robustness reasons.",
					"",
					"The format of the name is the same as when passed to `podman pull`. So, it supports using `:tag` or digests to guarantee the specific image version.",
					"",
					"Special case:",
					"- If the `name` of the image ends with `.image`, Quadlet will use the image pulled by the corresponding `.image` file, and the generated systemd service contains a dependency on the `$name-image.service` (or the service name set in the .image file). Note: the corresponding `.image` file must exist.",
				},
			},
			{
				Label: "Label",
				Hover: []string{
					"Set one or more OCI labels on the volume. The format is a list of `key=value` items, similar to `Environment`.",
					"",
					"This key can be listed multiple times.",
				},
				Macro: "Label=\"${1:key}=${2:value}\"\n$0",
			},
			{
				Label: "Options",
				Hover: []string{
					"The mount options to use for a filesystem as used by the **mount(8)** command `-o` option.",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman volume create` command in the generated file (right before the name of the volume in the command line). It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Type",
				Hover: []string{
					"The filesystem type of `Device` as used by the **mount(8)** commands `-t` option.",
				},
			},
			{
				Label: "User",
				Hover: []string{
					"The host (numeric) UID, or user name to use as the owner for the volume",
				},
			},
			{
				Label: "VolumeName",
				Hover: []string{
					"The (optional) name of the Podman volume. If this is not specified, the default value is the same name as the unit, but with a `systemd-` prefix, i.e. a `$name.volume` file creates a `systemd-$name` Podman volume to avoid conflicts with user-managed volumes.",
				},
			},
		},
		"Image": {
			{
				Label: "AllTags",
				Hover: []string{
					"All tagged images in the repository are pulled.",
					"This is equivalent to the Podman `--all-tags` option.",
				},
			},
			{
				Label: "Arch",
				Hover: []string{
					"Override the architecture, defaults to hosts, of the image to be pulled.",
					"This is equivalent to the Podman `--arch` option.",
				},
			},
			{
				Label: "AuthFile",
				Hover: []string{
					"Path of the authentication file.",
					"This is equivalent to the Podman `--authfile` option.",
				},
			},
			{
				Label: "CertDir",
				Hover: []string{
					"Use certificates at path (*.crt, *.cert, *.key) to connect to the registry.",
					"This is equivalent to the Podman `--cert-dir` option.",
				},
			},
			{
				Label: "ContainersConfModule",
				Hover: []string{
					"Load the specified containers.conf(5) module. Equivalent to the Podman `--module` option.",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Creds",
				Hover: []string{
					"The `[username[:password]]` to use to authenticate with the registry, if required.",
					"This is equivalent to the Podman `--creds` option.",
				},
			},
			{
				Label: "DecryptionKey",
				Hover: []string{
					"This key contains a list of arguments passed directly between `podman` and `image` in the generated file. It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Image",
				Hover: []string{
					"The image to pull. It is recommended to use a fully qualified image name rather than a short name, both for performance and robustness reasons.",
					"",
					"The format of the name is the same as when passed to `podman pull`. So, it supports using `:tag` or digests to guarantee the specific image version.",
				},
			},
			{
				Label: "ImageTag",
				Hover: []string{
					"Actual FQIN of the referenced `Image`. Only meaningful when source is a file or directory archive.",
					"",
					"For example, an image saved into a `docker-archive` with the following Podman command:",
					"",
					"`podman image save --format docker-archive --output /tmp/archive-file.tar quay.io/podman/stable:latest`",
					"",
					"requires setting",
					"- `Image=docker-archive:/tmp/archive-file.tar`",
					"- `ImageTag=quay.io/podman/stable:latest`",
				},
			},
			{
				Label: "OS",
				Hover: []string{
					"Override the OS, defaults to hosts, of the image to be pulled.",
					"This is equivalent to the Podman `--os` option.",
				},
			},
			{
				Label: "PodmanArgs",
				Hover: []string{
					"This key contains a list of arguments passed directly to the end of the `podman image pull` command in the generated file (right before the image name in the command line). It can be used to access Podman features otherwise unsupported by the generator. Since the generator is unaware of what unexpected interactions can be caused by these arguments, it is not recommended to use this option.",
					"",
					"The format of this is a space separated list of arguments, which can optionally be individually escaped to allow inclusion of whitespace and other control characters.",
					"",
					"This key can be listed multiple times.",
				},
			},
			{
				Label: "Policy",
				Hover: []string{
					"The pull policy to use when pulling the image.",
					"This is equivalent to the Podman `--policy` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 6, 0),
			},
			{
				Label: "Retry",
				Hover: []string{
					"Number of times to retry the image pull when a HTTP error occurs. Equivalent to the Podman `--retry` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "RetryDelay",
				Hover: []string{
					"Delay between retries. Equivalent to the Podman `--retry-delay` option.",
				},
				MinVersion: utils.BuildPodmanVersion(5, 5, 0),
			},
			{
				Label: "TLSVerify",
				Hover: []string{
					"Require HTTPS and verification of certificates when contacting registries.",
					"This is equivalent to the Podman `--tls-verify` option.",
				},
			},
			{
				Label: "Variant",
				Hover: []string{
					"Override the default architecture variant of the container image.",
					"This is equivalent to the Podman `--variant` option.",
				},
			},
		},
	}
)
