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

For a more details overview, with visual examples see the
[document](/docs/README.md).

<img src="docs/assets/overall_demo.gif" style="width: 100%;"/>

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

## Usage with Neovim

I'm using this language server with the following `neovim/nvim-lspconfig` setup
with LazyVim distribution. I'm relatively new in this world, so this setup may
have mistakes, but I was able to keep it working. Let's just consider it as an
example.

```lua
return {
    "neovim/nvim-lspconfig",
    opts = {
        servers = {
            quadlet_lsp = {},
        },
        setup = {
            quadlet_lsp = function(_, _)
                require("lspconfig.configs").quadlet_lsp = {
                    default_config = {
                        cmd = { "quadlet-lsp" },
                        filetypes = { "quadlet" },
                        root_dir = function(_)
                            return vim.fn.getcwd()
                        end,
                    },
                }

                vim.filetype.add({
                    extension = {
                        container = "quadlet",
                        volume = "quadlet",
                        network = "quadlet",
                        pod = "quadlet",
                        image = "quadlet",
                    },
                })

                vim.api.nvim_create_autocmd("FileType", {
                    pattern = "quadlet",
                    callback = function()
                        vim.bo.syntax = "systemd"
                    end,
                })

                return false
            end,
        },
    },
}
```

## Usage with VS Code

There is a simple VS Code extension to use it:

- [MarketPlace](https://marketplace.visualstudio.com/items?itemName=onlyati.quadlet-lsp)
- [Repository](https://github.com/onlyati/quadlet-lsp-vscode-extension)
