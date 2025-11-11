# Quadlet Syntax Rules

<!-- toc -->

- [`QSR001` - Missing section header](#qsr001---missing-section-header)
- [`QSR002` - Unfinished line](#qsr002---unfinished-line)
- [`QSR003` - Invalid property](#qsr003---invalid-property)
- [`QSR004` - Image name is not fully qualified](#qsr004---image-name-is-not-fully-qualified)
- [`QSR005` - Invalid value of AutoUpdate](#qsr005---invalid-value-of-autoupdate)
- [`QSR006` - Image file does not exists](#qsr006---image-file-does-not-exists)
- [`QSR007` - Invalid format of Environment variable](#qsr007---invalid-format-of-environment-variable)
- [`QSR008` - Invalid format of Annotation](#qsr008---invalid-format-of-annotation)
- [`QSR009` - Invalid format of Label](#qsr009---invalid-format-of-label)
- [`QSR010` - Incorrect format of PublishPort](#qsr010---incorrect-format-of-publishport)
- [`QSR011` - Port is not exposed in image](#qsr011---port-is-not-exposed-in-image)
- [`QSR012` - Invalid format of secret specification](#qsr012---invalid-format-of-secret-specification)
- [`QSR013` - Volume file does not exists](#qsr013---volume-file-does-not-exists)
- [`QSR014` - Network file does not exists](#qsr014---network-file-does-not-exists)
- [`QSR015` - Invalid format of Volume specification](#qsr015---invalid-format-of-volume-specification)
- [`QSR016` - Invalid value of UserNS specification](#qsr016---invalid-value-of-userns-specification)
- [`QSR017` - Pod file does not exists](#qsr017---pod-file-does-not-exists)
- [`QSR018` - Container cannot publish port with pod](#qsr018---container-cannot-publish-port-with-pod)
- [`QSR019` - Container cannot have network with pod](#qsr019---container-cannot-have-network-with-pod)
- [`QSR020` - Naming of unit is invalid](#qsr020---naming-of-unit-is-invalid)
- [`QSR021` - Unit points to not a systemd unit](#qsr021---unit-points-to-not-a-systemd-unit)
- [`QSR022` - '/' is before systemd directory specifier](#qsr022----is-before-systemd-directory-specifier)
- [`QSR023` - Invalid systemd specifier is used](#qsr023---invalid-systemd-specifier-is-used)
- [`QSR024` - Not recommended property usage in Service section](#qsr024---not-recommended-property-usage-in-service-section)
- [`QSR025` - Image is missing in container](#qsr025---image-is-missing-in-container)
- [`QSR026` - Artifact is missing in artifact](#qsr026---artifact-is-missing-in-artifact)

<!-- tocstop -->

## `QSR001` - Missing section header

**Message**

> Missing any of these sections: _%section_list%_

**Explanation**

This is a Quadlet file, but lack of section headers:

- Image
- Container
- Volume
- Network
- Kube
- Pod

## `QSR002` - Unfinished line

**Message**

> Line is unfinished

**Explanation**

This is error is reported when a keyword is used but nothing after equation
sign.

```ini
[Container]
Image=     # <-- Error here because it is unfinished
AutoUpdate=registry
```

## `QSR003` - Invalid property

**Message**

> Invalid property is found: _%section%.%property_name%_

**Explanation**

The typed property does not exists.

```ini
[Container]
Image=docker.io/library/nextcloud:fpm
AutoUpdat=registry # <-- Invalid value at the left side of '='
```

## `QSR004` - Image name is not fully qualified

**Message**

> Image name is not fully qualified

**Explanation**

The specified image name is not fully qualified:

```ini
[Container]
Image=debian:bookworm-slim
```

Use fully qualified image name instead. It means that the image must begin with
`localhost` or with a valid address (e.g.: something.tld). It can also contains
port number.

Some examples:

```ini
Image=docker.io/library/debian:bookworm-slim
Image=ghcr.io/henrygd/beszel/beszel-agent
Image=localhost/test
Image=localhost:5000/test
Image=example.com:5000/test
```

## `QSR005` - Invalid value of AutoUpdate

**Message**

> Invalid value of AutoUpdate: _%value%_

**Explanation**

The `AutoUpdate` can only have `local` and `registry` values.

## `QSR006` - Image file does not exists

**Message**

> Image file does not exists: _%name%_

**Explanation**

The specified `*.image` or `*.build` file does not exists that is used in the
`Image=` line.

## `QSR007` - Invalid format of Environment variable

**Message**

> Invalid format: _%reason%_

**Explanation**

Environment variables are represented as key-value pairs. If you need to assign
a value containing spaces or the equals sign to a variable, put quotes around
the whole assignment. Variable expansion is not performed inside the strings and
the "$" character has no special meaning.

This option may be specified more than once, in which case all listed variables
will be set. If the same variable is listed twice, the later setting will
override the earlier setting. If the empty string is assigned to this option,
the list of environment variables is reset, all prior assignments have no
effect.

Correct examples:

```ini
Environment=FOO=BAR "MyVar=MyValue" 'foo=bar'
Environment=FOO=
Environment='fooVariable=barValue'
```

After `5.6.0` version, it accept environment variable without assignment. In
this case, the value is get from the host.

```ini
Environment=FOO
```

## `QSR008` - Invalid format of Annotation

**Message**

> Invalid format: _%reason%_

**Explanation**

Annotation variables are represented as key-value pairs. If you need to assign a
value containing spaces or the equals sign to a variable, put quotes around the
whole assignment. Variable expansion is not performed inside the strings and the
"$" character has no special meaning.

This option may be specified more than once, in which case all listed variables
will be set. If the same variable is listed twice, the later setting will
override the earlier setting. If the empty string is assigned to this option,
the list of environment variables is reset, all prior assignments have no
effect.

Correct examples:

```ini
Annotation=FOO=BAR "MyVar=MyValue" 'foo=bar'
Annotation=FOO=
Annotation='fooVariable=barValue'
```

## `QSR009` - Invalid format of Label

**Message**

> Invalid format: _%reason%_

**Explanation**

Label variables are represented as key-value pairs. If you need to assign a
value containing spaces or the equals sign to a variable, put quotes around the
whole assignment. Variable expansion is not performed inside the strings and the
"$" character has no special meaning.

This option may be specified more than once, in which case all listed variables
will be set. If the same variable is listed twice, the later setting will
override the earlier setting. If the empty string is assigned to this option,
the list of environment variables is reset, all prior assignments have no
effect.

Correct examples:

```ini
Label=FOO=BAR "MyVar=MyValue" 'foo=bar'
Label=FOO=
Label='fooVariable=barValue'
```

## `QSR010` - Incorrect format of PublishPort

**Message**

> Incorrect format of PublishPort

**Explanation**

Valid formats for `PublishPort`:

```ini
PublishPort=10.0.0.1:10069:69
PublishPort=10420:420
```

## `QSR011` - Port is not exposed in image

**Message**

> Port is not exposed in the image, exposed ports: _%port_list%_

or

> Not able to verify exposed ports, because image not pulled: _%image_list%_

**Explanation**

Port is used in container or pod that is not exposed by the image. In case of
pod, first it discover which other container files are linked for the pod and
analyze those images.

## `QSR012` - Invalid format of secret specification

**Message**

> Invalid format of secret specification: _%reason%_

**Explanation**

Depends on `reason` text:

- `%opt% has no value`: Invalid option
- `'type' can be either 'mount' or 'env'`: Target is specified but with invalid
  value
- `'%opt%' only allowed if type=mount`: Using `uid`, `gid` or `mode` meanwhile
  not `type=env` is set

## `QSR013` - Volume file does not exists

**Message**

> Volume file does not exists: _%volume_file%_

**Explanation**

The defined file, e.g.: `Volume=data.volume:/data`, does not exists in the
current working directory.

## `QSR014` - Network file does not exists

**Message**

> Network file does not exists: _%network_file%_

**Explanation**

The defined file, e.g.: `Network=my.network`, does not exists in the current
working directory.

## `QSR015` - Invalid format of Volume specification

**Message**

> Invalid format of Volume specification: _%reason%_

**Explanation**

Depends on the `reason`:

- `container directory is not absolute`: Container directory must be absolute
- `'%flag%' is unkown`: Not existing flag is used

## `QSR016` - Invalid value of UserNS specification

**Message**

> Invalid value of UserNS: allowed values: _%reason%_

**Explanation**

Depends on the values of `reason`:

- `%opt% has no paramerets`: Only `keep-id` can have further parameters
- `[uid gid] allowed but found %opt%`: Uses `keep-id` with other parameters than
  `uid` or `gid`
- `allowed values: [auto host keep-id nomap] but found %opt%`: Using invalid
  value of `UserNS`

## `QSR017` - Pod file does not exists

**Message**

> Pod file does not exists: _%pod_file%_

**Explanation**

The defined file, e.g.: `Pod=my.pod`, does not exists in the current working
directory.

## `QSR018` - Container cannot publish port with pod

**Message**

> Container cannot have PublishPort because belongs to a pod: _%pod_file%_

**Explanation**

A Pod in Podman shares a network namespace across all containers inside it. The
pod is the unit that binds to the host network (e.g., 127.0.0.1:8080), not the
individual containers.

Each container in the pod uses 127.0.0.1 to reach other containers in the same
pod.

## `QSR019` - Container cannot have network with pod

**Message**

> Container cannot have Network because belongs to a pod: _%pod_file%_

**Explanation**

When you create a pod, it gets a single network namespace that all containers in
the pod share. So: Containers in the same pod communicate over localhost
(127.0.0.1). You assign the network (e.g. --network) when creating the pod, not
per container.

## `QSR020` - Naming of unit is invalid

**Message**

> Invalid name of unit: _%name%_

**Explanation**

Container, Volume, Pod and Network naming must match with
`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$` regular expression. The specified name at
`ContainerName`, `VolumeName`, `PodName` or `Networkname` does not match with
the expression.

## `QSR021` - Unit points to not a systemd unit

Quadlet will automatically translate dependencies, specified in the keys
`Wants`, `Requires`, `Requisite`, `BindsTo`, `PartOf`, `Upholds`, `Conflicts`,
`Before` and `After` of the `[Unit]` section, between different Quadlet units.

But this is true only after `5.5.0` version. This rule gives error if dependency
not translated and version if before `5.5.0`.

For example the `fedora.container` unit below specifies a dependency on the
`basic.container` unit.

```ini
[Unit]
After=basic.container
Requires=basic.container

[Container]
Image=registry.fedoraproject.org/fedora:41
```

Before `5.5.0` version, file above should look:

```ini
[Unit]
After=basic-container.service
Requires=basic-container.service

[Container]
Image=registry.fedoraproject.org/fedora:41
```

## `QSR022` - '/' is before systemd directory specifier

**Message**

> Specifier, _%specifier%_, already starts with '/' sign

**Explanation**

The _specifier_ is a directory and it already begin with `/` sign. For example
`/%T` would evaluate as `//tmp` which is invalid path.

## `QSR023` - Invalid systemd specifier is used

**Message**

> Specifier '_%specifier%_' is invalid

**Explanation**

Specifier does not exists:
[specifiers](https://www.freedesktop.org/software/systemd/man/latest/systemd.unit.html#Specifiers).
If you want to write '%' then escape it as '%%'.

## `QSR024` - Not recommended property usage in Service section

**Message**

> Usage in rootless podman is not recommended: Service._%property%_

**Explanation**

Note that Quadlet units do not support running as a non-root user by defining
the
[User, Group](https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html#User=),
or
[DynamicUser](https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html#DynamicUser=)
systemd options. If you want to run a rootless Quadlet, you will need to create
the user and add the unit file to one of the above rootless unit search paths.

If this is a rootful Podman, you can ignore the message.

## `QSR025` - Image is missing in container

**Message**

> Container Quadlet file does not have Image property

**Explanation**

The `Image` property is mandatory to specify in container Quadlet files.

## `QSR026` - Artifact is missing in artifact

**Message**

> Artifact Quadlet file does not have Artifact property

**Explanation**

The `Artifact` property is mandatory to specify in artifact Quadlet files.
