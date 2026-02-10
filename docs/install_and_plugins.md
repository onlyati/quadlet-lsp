# Installation & Plugins

Language server can be used with various code editors out of the box:

- [Neovim](https://github.com/onlyati/quadlet-lsp.nvim)
- [VS Code](https://marketplace.visualstudio.com/items?itemName=onlyati.quadlet-lsp)
- [Zed](https://zed.dev/extensions/quadlet) (3rd party)

## Download

If you would like to use it with Neovim or as a CLI, there are several source
where it can be installed. The VS Code and Zed plugins are already contains the
language server.

### Mise

Easy to download and use it with [mise](https://github.com/jdx/mise).

```bash
mise use -g github:onlyati/quadlet-lsp
```

### Install from Fedora copr

This is supported on Fedora 43, RHEL 10, Alma Linux 10, Rocky Linux 10.

```bash
sudo dnf copr enable onlyati/quadlet-lsp
sudo dnf install quadlet-lsp
```

### Install from Debian registry

This method is supported for Debian 13. Add the following registry, that hosted
by me, then update:

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

### Download with Go

You can also install the binary using Go.

```bash
go install github.com/onlyati/quadlet-lsp@latest
```
