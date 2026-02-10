# Podman Quadlet language server

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

For a more details overview, see
[documentation site](https://quadlet-lsp.thinkaboutit.tech).

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
