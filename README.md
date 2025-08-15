# Podman Quadlet language server

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

### Download the latest version

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
