package data

type SystemdSpecifier struct {
	ShortDescription string
	LongDescription  []string
	IsDirectory      bool
}

var SystemdSpecifierSet = map[string]SystemdSpecifier{
	"%a": {
		ShortDescription: "Architecture",
		LongDescription: []string{
			"A short string identifying the architecture of the local system. A string such as x86, x86-64 or arm64.",
			"See the architectures defined for `ConditionArchitecture=` above for a full list.",
		},
		IsDirectory: false,
	},
	"%A": {
		ShortDescription: "Operating system image version",
		LongDescription: []string{
			"The operating system image version identifier of the running system, as read from the IMAGE_VERSION= field of /etc/os-release. ",
			"If not set, resolves to an empty string. See os-release(5) for more information.",
		},
		IsDirectory: false,
	},
	"%b": {
		ShortDescription: "Boot ID",
		LongDescription: []string{
			"The boot ID of the running system, formatted as string. See random(4) for more information.",
		},
		IsDirectory: false,
	},
	"%B": {
		ShortDescription: "Operating system build ID",
		LongDescription: []string{
			"The operating system build identifier of the running system, as read from the BUILD_ID= field of /etc/os-release. ",
			"If not set, resolves to an empty string. See os-release(5) for more information.",
		},
	},
	"%C": {
		ShortDescription: "Cache directory",
		LongDescription: []string{
			"This is either /var/cache (for the system manager) or the path `$XDG_CACHE_HOME` resolves to (for user managers).",
		},
		IsDirectory: true,
	},
	"%d": {
		ShortDescription: "Credentials directory",
		LongDescription: []string{
			"This is the value of the `$CREDENTIALS_DIRECTORY` environment variable if available.",
			"See section `Credentials` in systemd.exec(5) for more information.",
		},
		IsDirectory: true,
	},
	"%D": {
		ShortDescription: "Shared data directory",
		LongDescription: []string{
			"This is either /usr/share/ (for the system manager) or the path `$XDG_DATA_HOME` resolves to (for user managers).",
		},
		IsDirectory: true,
	},
	"%E": {
		ShortDescription: "Configuration directory root",
		LongDescription: []string{
			"This is either /etc/ (for the system manager) or the path `$XDG_CONFIG_HOME` resolves to (for user managers).",
		},
		IsDirectory: true,
	},
	"%f": {
		ShortDescription: "Unescaped file name",
		LongDescription: []string{
			"This is either the unescaped instance name (if applicable) with / prepended (if applicable), or the unescaped prefix name prepended with /. ",
			"This implements unescaping according to the rules for escaping absolute file system paths discussed above.",
		},
		IsDirectory: false,
	},
	"%g": {
		ShortDescription: "User group",
		LongDescription: []string{
			"This is the name of the group running the service manager instance. In case of the system manager this resolves to `root`.",
		},
		IsDirectory: false,
	},
	"%G": {
		ShortDescription: "User GID",
		LongDescription: []string{
			"This is the numeric GID of the user running the service manager instance. In case of the system manager this resolves to `0`.",
		},
	},
	"%h": {
		ShortDescription: "User home directory",
		LongDescription: []string{
			"This is the home directory of the user running the service manager instance. In case of the system manager this resolves to `/root`. ",
			"Note that this setting is not influenced by the User= setting configurable in the [Service] section of the service unit.",
		},
		IsDirectory: true,
	},
	"%H": {
		ShortDescription: "Host name",
		LongDescription: []string{
			"The hostname of the running system at the point in time the unit configuration is loaded.",
		},
		IsDirectory: false,
	},
	"%i": {
		ShortDescription: "Instance name",
		LongDescription: []string{
			"For instantiated units this is the string between the first `@` character and the type suffix. Empty for non-instantiated units.",
		},
		IsDirectory: false,
	},
	"%I": {
		ShortDescription: "Unescaped instance name",
		LongDescription: []string{
			"Same as `%i`, but with escaping undone.",
		},
		IsDirectory: false,
	},
	"%j": {
		ShortDescription: "Final component of the prefix",
		LongDescription: []string{
			"This is the string between the last `-` and the end of the prefix name. If there is no `-`, this is the same as `%p`.",
		},
		IsDirectory: false,
	},
	"%J": {
		ShortDescription: "Unescaped final component of the prefix",
		LongDescription: []string{
			"Same as `%j`, but with escaping undone.",
		},
		IsDirectory: false,
	},
	"%l": {
		ShortDescription: "Short host name",
		LongDescription: []string{
			"The hostname of the running system at the point in time the unit configuration is loaded, truncated at the first dot to remove any domain component.",
		},
		IsDirectory: false,
	},
	"%L": {
		ShortDescription: "Log directory root",
		LongDescription: []string{
			"This is either /var/log (for the system manager) or the path $XDG_STATE_HOME resolves to with /log appended (for user managers).",
		},
		IsDirectory: true,
	},
	"%m": {
		ShortDescription: "Machine ID",
		LongDescription: []string{
			"The machine ID of the running system, formatted as string. See machine-id(5) for more information.",
		},
		IsDirectory: false,
	},
	"%M": {
		ShortDescription: "Operqating system image identifier",
		LongDescription: []string{
			"The operating system image identifier of the running system, as read from the IMAGE_ID= field of /etc/os-release. ",
			"If not set, resolves to an empty string. See os-release(5) for more information.",
		},
		IsDirectory: false,
	},
	"%n": {
		ShortDescription: "Full unit name",
		LongDescription:  []string{},
		IsDirectory:      false,
	},
	"%N": {
		ShortDescription: "Full unit name",
		LongDescription: []string{
			"Same as `%n`, but with the type suffix removed.",
		},
		IsDirectory: false,
	},
	"%o": {
		ShortDescription: "Operating system ID",
		LongDescription: []string{
			"The operating system identifier of the running system, as read from the ID= field of /etc/os-release. See os-release(5) for more information.",
		},
		IsDirectory: false,
	},
	"%p": {
		ShortDescription: "Prefix name",
		LongDescription: []string{
			"For instantiated units, this refers to the string before the first `@` character of the unit name. For non-instantiated units, same as `%N`.",
		},
		IsDirectory: false,
	},
	"%P": {
		ShortDescription: "Unescaped prefix name",
		LongDescription: []string{
			"Same as `%p`, but with escaping undone.",
		},
		IsDirectory: false,
	},
	"%q": {
		ShortDescription: "Pretty host name",
		LongDescription: []string{
			"The pretty hostname of the running system at the point in time the unit configuration is loaded, as read from the PRETTY_HOSTNAME= field of /etc/machine-info. ",
			"If not set, resolves to the short hostname. See machine-info(5) for more information.",
		},
		IsDirectory: false,
	},
	"%s": {
		ShortDescription: "User shell",
		LongDescription: []string{
			"This is the shell of the user running the service manager instance.",
		},
		IsDirectory: false,
	},
	"%S": {
		ShortDescription: "State directory root",
		LongDescription: []string{
			"This is either /var/lib (for the system manager) or the path $XDG_STATE_HOME resolves to (for user managers).",
		},
		IsDirectory: true,
	},
	"%t": {
		ShortDescription: "Runtime directory root",
		LongDescription: []string{
			"This is either /run/ (for the system manager) or the path `$XDG_RUNTIME_DIR` resolves to (for user managers).",
		},
		IsDirectory: true,
	},
	"%T": {
		ShortDescription: "Directory for temporary files",
		LongDescription: []string{
			"This is either /tmp or the path `$TMPDIR`, `$TEMP` or `$TMP` are set to. (Note that the directory may be specified without a trailing slash.)",
		},
		IsDirectory: true,
	},
	"%u": {
		ShortDescription: "User name",
		LongDescription: []string{
			"This is the name of the user running the service manager instance. In case of the system manager this resolves to `root`. ",
			"Note that this setting is not influenced by the User= setting configurable in the [Service] section of the service unit.",
		},
		IsDirectory: false,
	},
	"%U": {
		ShortDescription: "User UID",
		LongDescription: []string{
			"This is the numeric UID of the user running the service manager instance. In case of the system manager this resolves to `0`.",
			"Note that this setting is not influenced by the User= setting configurable in the [Service] section of the service unit.",
		},
		IsDirectory: false,
	},
	"%v": {
		ShortDescription: "Kernel release",
		LongDescription: []string{
			"Identical to uname -r output.",
		},
		IsDirectory: false,
	},
	"%V": {
		ShortDescription: "Directory for larger and persistent temporary files",
		LongDescription: []string{
			"This is either /var/tmp or the path `$TMPDIR`, `$TEMP` or `$TMP` are set to. (Note that the directory may be specified without a trailing slash.)",
		},
		IsDirectory: true,
	},
	"%w": {
		ShortDescription: "Operating system version ID",
		LongDescription: []string{
			"The operating system version identifier of the running system, as read from the VERSION_ID= field of /etc/os-release.",
			"If not set, resolves to an empty string. See os-release(5) for more information.",
		},
		IsDirectory: false,
	},
	"%W": {
		ShortDescription: "Operating system variant ID",
		LongDescription: []string{
			"The operating system variant identifier of the running system, as read from the VARIANT_ID= field of /etc/os-release. ",
			"If not set, resolves to an empty string. See os-release(5) for more information.",
		},
		IsDirectory: false,
	},
	"%y": {
		ShortDescription: "The path of the fragement",
		LongDescription: []string{
			"This is the path where the main part of the unit file is located. For linked unit files, the real path outside of the unit search directories is used. ",
			"For units that do not have a fragment file, this specifier will raise an error.",
		},
		IsDirectory: false,
	},
	"%Y": {
		ShortDescription: "The directory of the fragement",
		LongDescription: []string{
			"This is the directory part of `%y`.",
		},
		IsDirectory: false,
	},
	"%%": {
		ShortDescription: "Single percent sign",
		LongDescription: []string{
			"Use `%%` in place of `%` to specify a single percent sign.",
		},
		IsDirectory: false,
	},
}
