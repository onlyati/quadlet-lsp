# Development documentation

<!-- toc -->

- [Getting started](#getting-started)
- [Project structure](#project-structure)
    * [Syntax rules](#syntax-rules)
    * [Hover](#hover)
    * [Commands](#commands)
- [Contribution](#contribution)

<!-- tocstop -->

This is a
[language server](https://microsoft.github.io/language-server-protocol/overviews/lsp/overview/)
for Podman Quadlets. Purpose of this language server is to provide a simple and
easy way to modify Quadlet files.

Basically this language server is also designed to work with Quadlets without
running them. For example, Quadlets are in a repository, it is pulled and
modified, but not directly run on the developer machine but on a server. In this
situation, language server can help with syntax rules to provide better way to
validate the files. Because of this, language server does not assume that user
is working in `~/.config/containers/systemd` directory. It checks files in the
current working directory of IDE, which information is sent by IDE during
`initialize` call.

Besides language server nature, this also can be used in CI/CD to check syntax
of Quadlet files, by using `check` parameter.

Following key features are supported:

- **Completion:** Provide static completion (e.g. keywords) and also provide
  dynamic completions (e.g.: listing currently existing volumes).
- **Hover:** Provide explanation for keywords and values for better
  understanding of Quadlet.
- **Go definition/Go reference:** Quadlets are split to multiple files, not like
  compose files. This split makes it easier to read, since they are shorter. But
  it makes also more difficult because it required to move among files. The
  `go definition` and `go reference` features can help to navigate to the proper
  file.
- **Commands:** Commands that can be useful, like pull down all images.
- **Extra options:** Like disable syntax rules or specify Podman version (in
  case of the deployment server has different Podman than the developer
  machine). These are done by `.quadletrc.json` file in the root working
  directory.

This project uses [tliron/glsp](https://github.com/tliron/glsp) 3.16 LSP
implementation.

## Getting started

For simplicity this project uses [mise](https://github.com/jdx/mise). Visit the
website about installation. The `mise.toml` file contains the uses go version
which is picked up by `mise` with the proper terminal integration.

You can use the following `mise tasks` to easily perform actions:

```bash
$ mise task ls
Name       Description
build      Build languager server to your system
build_all  Make an offline release
test       Run unit tests
```

The `build_all` command uses `goreleaser` but this is included as tool onto
`mise.toml` file.

## Project structure

Documents are stored in the memory of language server and they are modified
during `TextDocumentDidOpen`, `TextDocumentDidChang` and `TextDocumentDidClose`
calls.

All components (e.g.: syntax, hover) are read the documents here, from the
memory of the language server.

### Syntax rules

Quadlet Syntax Rules (QSR) are in `internal/syntax` directory. The file name
must match with the syntax rule identifier. For example, for `QSR021` you must
have:

- `internal/syntax/qsr021.go`: Store the logic.
- `internal/syntax/qsr021_test.go`: Store unit tests for syntax rule.
- `internal/syntax/syntax.go`: Register rule to the `NewSyntaxChecker` function.

The `internal/syntax/common.go` has a function called `canFileBeApplied`. This
should be called to verify which file extension are eligible for the syntax
check.

Then file is scanned using `utils.ScanQadlet` function. This function has the
following parameters:

- _document text_: Full text of the document.
- _podman version_: This version is passed down to the _action_ parameter. If
  the subject is independent from the version, then placeholder is enough.
- _selector map_: specify a `map[utils.ScanProperty]struct{}` hash set which
  line the _action_ should run. If action needs to be called at all lines, then
  map must have a `utils.ScanProperty{ Section: "*", Property: "*" }` element.
- _action_: This action doing the main part of the syntax check. This action is
  called every line that is filtered by line.

Example for syntax rule:

```go
func qsr024(s SyntaxChecker) []protocol.Diagnostic {
    diags := []protocol.Diagnostic{}

    allowedFiles := []string{
        "image",
        "container",
        "volume",
        "network",
        "kube",
        "pod",
        "build",
    }
    if c := canFileBeApplied(s.uri, allowedFiles); c != "" {
        diags = utils.ScanQadlet(
            s.documentText,
            // Placeholder here, but actual can be read from `s.config.Podman`
            utils.PodmanVersion{},
            map[utils.ScanProperty]struct{}{
                {Section: "[Service]", Property: "User"}:        {},
                {Section: "[Service]", Property: "Group"}:       {},
                {Section: "[Service]", Property: "DynamicUser"}: {},
            },
            qsr024Action,
        )
    }

    return diags
}
```

The action function must have parameters like this.

```go
func qsr024Action(q utils.QuadletLine, _ utils.PodmanVersion) []protocol.Diagnostic {
    // Do syntax validation logic here

    // One or more diagnostic element can be send back, sample:
    return []protocol.Diagnostic{
        {
            Range: protocol.Range{
                Start: protocol.Position{Line: q.LineNumber, Character: 0},
                End:   protocol.Position{Line: q.LineNumber, Character: q.Length},
            },
            Severity: &warnDiag,
            Message:  fmt.Sprintf("Usage in rootless podman is not recommended: %s.%s", "Service", q.Property),
            Source:   utils.ReturnAsStringPtr("quadlet-lsp.qsr024"),
        },
    }
}
```

Don't forget to write unit tests to validate the functions and update
`docs/qsr.md` file accordingly.

### Hover

For static completion the `internal/data/properties.go` file is used. It has map
structures that hold information and completion is suggested based on this. For
example, we have this item in the map, in the `Container` section:

````go
PropertyMapItem {
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
````

Language server will provide property completion if there is no `=` sign in the
line yet. If the line start with `DNS=` and the `Parameters` list is not null,
then provide their values as completion. By this, static completion is done
purely on data basis.

For dynamic completions that requires more logic (e.g.: listing secrets, cursor
position matter), then it must be place into a separate file, like
`internal/completion/property_image.go` and corresponding unit test
(`internal/completion/property_image_test.go`).

The completion function must has parameters and return value, like:

```go
func propertyListImages(s Completion) []protocol.CompletionItem {
    var completionItems []protocol.CompletionItem

    // Do the logic here

    return completionItems
}
```

That the completion function is triggered, it must be registers in
`listPropertyCompletions` function in `internal/completion/property.go` file.

Don't forget to write unit tests and update `docs/features.md` file accordingly.

### Commands

The commander module is not finalized, it will be probably redesigned in the
near future.

## Contribution

All contribution are welcomed, not just PRs but also issues to highlight issues
or new potential features in the project.
