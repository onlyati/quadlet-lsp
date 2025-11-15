# Podman Quadlet language server

<!-- toc -->

- [Usage with Neovim](#usage-with-neovim)
- [Usage with VS Code](#usage-with-vs-code)
- [Usage with Zed](#usage-with-zed)
- [Alternate usage](#alternate-usage)
- [Get the executable](#get-the-executable)
    * [Use with mise](#use-with-mise)
    * [Install from Fedora copr](#install-from-fedora-copr)
    * [Install from Debian registry](#install-from-debian-registry)
    * [Install from .deb and .rpm package](#install-from-deb-and-rpm-package)
    * [Download the compiled version](#download-the-compiled-version)
    * [Compile with Go](#compile-with-go)

<!-- tocstop -->

> [!CAUTION]
>
> The main branch may unstable. Use version tagged code to get stable code.

This is an implementation of language server for
[Podman Quadlet](https://docs.podman.io/en/latest/markdown/podman-systemd.unit.5.html#description)
files.

Following features are currently available:

- Code completion
  - Provide static completion based on Podman Quadlet documentation
  - Query images, volumes, networks, pods, and so on, and provide completion
    based on real configuration
- Hover menu
- Implemented "go definition" and "go references" functions
- Provide syntax checking
- Execute built-in commands

For a more details overview, with visual examples see the
[document](/docs/README.md).

<img src="docs/assets/overall_demo.gif" style="width: 100%;"/>

## Usage with Neovim

There is a plugin made for this language server:

- [Repository](https://github.com/onlyati/quadlet-lsp.nvim)

## Usage with VS Code

There is a simple VS Code extension to use it:

- [MarketPlace](https://marketplace.visualstudio.com/items?itemName=onlyati.quadlet-lsp)
- [Repository](https://github.com/onlyati/quadlet-lsp-vscode-extension)

## Usage with Zed

There is a third-party Zed extension that makes use of `quadlet-lsp`:

- [Marketplace](https://zed.dev/extensions/quadlet)
- [Repository](https://github.com/mufeedali/zed-quadlet/)

## Alternate usage

This binary can be used as a CLI syntax checker for Quadlet files. This can be
useful, for example in CI/CD pipeline to verify Quadlets before packaging and
later deploying.

Same rule applied for CLI that is also applied for editor syntax checking. The
`.quadletrc.json` file is also used on same way.

Example for usage, that monitor the current working directory:

```bash
$ quadlet-lsp check .
nc-db.container     , quadlet-lsp.qsr003, 09.000-09.010, Invalid property is found: Container.Memory
nc-app.container    , quadlet-lsp.qsr003, 08.000-08.010, Invalid property is found: Container.Memory
$ echo $?
4
```

Binary return with non-zero return code if it find any non information finding.

## Get the executable

### Use with mise

You can easily download, update and use with
[mise](https://github.com/jdx/mise).

```bash
mise use -g ubi:onlyati/quadlet-lsp
```

### Install from Fedora copr

This is supported on Fedora 42/43, RHEL 10, Alma Linux 10, Rocky Linux 10.

```bash
sudo dnf copr enable onlyati/quadlet-lsp
sudo dnf install quadlet-lsp
```

### Install from Debian registry

This method is supported for Debian 13. Add the following registry, then update:

```bash
$ sudo curl \
    https://git.thinkaboutit.tech/api/packages/pandora/debian/repository.key \
    -o /etc/apt/keyrings/gitea-pandora.asc
$ sudo tee /etc/apt/sources.list.d/onlyati.sources > /dev/null <<'EOF'
Types: deb
URIs: https://git.thinkaboutit.tech/api/packages/pandora/debian
Suites: trixie
Components: main
Signed-By: /etc/apt/keyrings/gitea-pandora.asc
EOF
$ sudo apt update
```

Then simply install:

```bash
sudo apt install quadlet-lsp
```

### Install from .deb and .rpm package

Check GitHub [release page](https://github.com/onlyati/quadlet-lsp/releases) and
download the version you need, then install it manually.

### Install from Nix flake

Add this repo to your flake's inputs.

```nix
inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
    quadlet-lsp.url = "github:onlyati/quadlet-lsp";
};
```

Add the default package to your packages list.

```nix
environment.systemPackages = [
    inputs.quadlet-lsp.packages.${system}.default
];

# Home Manager
home.packages = [
    inputs.quadlet-lsp.packages.${system}.default
];
```

### Download the compiled version

Check GitHub [release page](https://github.com/onlyati/quadlet-lsp/releases) and
download the version you need. The archive contains only the binary of language
server.

On Linux, you can get it quicker from terminal. See example commands.

```bash
ARCH="amd64"
OS="linux"
LATEST_VERSION=$(curl -s -H "Accept: application/vnd.github+json" \
    https://api.github.com/repos/onlyati/quadlet-lsp/releases/latest \
    | jq -r .tag_name)
rm quadlet-lsp-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz
wget "https://github.com/onlyati/quadlet-lsp/releases/download/${LATEST_VERSION}/quadlet-lsp-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"
sudo tar -xvf "quadlet-lsp-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz" \
    -C /usr/local/bin/
```

### Compile with Go

You can also install the binary using Go.

```bash
go install github.com/onlyati/quadlet-lsp@latest
```
