# Features

<!-- toc -->

- [Hover menu](#hover-menu)
- [Completion](#completion)
    * [Starter template](#starter-template)
    * [Static completion](#static-completion)
    * [Templates](#templates)
    * [Dynamic completion](#dynamic-completion)
    * [Go definition/references](#go-definitionreferences)

<!-- tocstop -->

Features below are implemented to support following file extensions:

- `*.image`
- `*.contianer`
- `*.volume`
- `*.network`
- `*.kube`
- `*.pod`

> [!IMPORTANT]
>
> Only Quadlet part has features in the files below. The generic systemd related
> parts are not covered.

## Hover menu

Provide some information about specific property. See a demo about a container
file in the following video.

<img src="assets/hover_demo.gif" style="width: 100%;"/>

## Completion

### Starter template

The `newContainer`, `newVolume`, and so on, provide a started template for
specific files.

<img src="assets/newContainer_demo.gif" style="width: 100%;"/>

### Static completion

Language server provide some static completion based on Podman Quadlet
Documentation, like `Exec`, `Environment`, and so on.

<img src="assets/static_comp_demo.gif" style="width: 100%;"/>

### Templates

Some property has a "new template". If you type them you can get predefined
snippets. Currently supported new templates:

- `new.Annotation`
- `new.AddHost`
- `new.Environment`
- `new.Label`
- `new.PublishPort`
- `new.Secret`
- `new.Volume`

<img src="assets/new_env_demo.gif" style="width: 100%;"/>

### Dynamic completion

Language server provide some dynamic completion:

- List pulled images and `*.image` files at `Image=`
- List defined secrets at `Secret=`. Also further parameters (type, target)
- List created volumes and `*.volume` files at `Volume=`. Also further
  parameters (rw, ro, z, Z)
- List `*.pod` files at `Pod=`
- List created networks and `*.network` files at `Network=`
- Gather and list `uid` and `gid` from image if `UserNS=keep-id:` is specified
- Gather exposed ports from image and provide them when `PublishPort` is
  specified

<img src="assets/din_comp_demo.gif" style="width: 100%;"/>

### Go definition/references

If you are on a line that points to another file, e.g.: `Pod=nc.pod` and using
the `go definition` function, the file is open.

If you are on a line like `[Pod]`, `[Volume]`, `[Network]`, `[Image]`, then
current work directory is searched for any references to that specific file.

<img src="assets/go_def_ref_demo.gif" style="width: 100%; max-width: 800px;"/>
