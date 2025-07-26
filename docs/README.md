# Documentation

## Details

- [Language features](./features.md)
- [Quadlet Syntax Rules](./qsr.md)

## Usage with Neovim

During the tests, I've used LazyVim distribution of Neovim with the following
LSP configuration in `~/.config/nvim/lua/plugin/quadlet_lsp.lua` file.

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
                        kube = "quadlet",
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
