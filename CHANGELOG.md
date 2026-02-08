# Changelog of Quadlet-LSP

## V0.7.2

### Features

- Add file preview for `Network`, `Pod` and `Volume` on hover
  (<https://github.com/onlyati/quadlet-lsp/pull/176>)
- Remove templates and build into the property completion
  (<https://github.com/onlyati/quadlet-lsp/pull/182>)

### Bugfixes

- `QSR006` looking for in nested directories
  (<https://github.com/onlyati/quadlet-lsp/pull/181>)

### Repo

- Update to go 1.25.5 (<https://github.com/onlyati/quadlet-lsp/pull/177>)
- Fix lint issues, add linter for verification
  (<https://github.com/onlyati/quadlet-lsp/pull/178>)
- Review unit tests (<https://github.com/onlyati/quadlet-lsp/pull/180>)

## v0.7.1

### Bugfixes

- Completion crashes when previous line is empty
  (<https://github.com/onlyati/quadlet-lsp/pull/172>)

## v0.7.0

### Features

- Support nested directories for each function of languager server
  (<https://github.com/onlyati/quadlet-lsp/issues/148>)
- Add protocol suffix support to PublishPort validation
  (<https://github.com/onlyati/quadlet-lsp/pull/150>)
- Add warning (instead of just a crash) if not directory is open
  (<https://github.com/onlyati/quadlet-lsp/pull/165>)

### Bugfixes

- Add completion item type (e.g.: `Value`, `Keyword`, etc.)
  (<https://github.com/onlyati/quadlet-lsp/pull/154>)
- Fix property name for disable rules in documentation
  (<https://github.com/onlyati/quadlet-lsp/pull/159>)
- Fix `rshared` flag in QSR15
  (<https://github.com/onlyati/quadlet-lsp/pull/161>)
- Property completion did not work if `=` was already in the line
  (<https://github.com/onlyati/quadlet-lsp/pull/166>)

## v0.6.0

### Features

- Support Podman v5.7.0 features
  <https://github.com/onlyati/quadlet-lsp/issues/136>:
  - Quadlet .container files now support a new key, HttpProxy, to disable the
    automatic forwarding of HTTP proxy options from the host into the container.
  - Quadlet .pod files now support a new key, StopTimeout, to configure the stop
    timeout for the pod
  - Quadlet .build files now support two new keys, BuildArg and IgnoreFile, to
    specify build arguments and an ignore file
  - Implement hover, code completion and starter template for artifact Quadlets.
  - Implement new rule (QSR026).
- Add completion for `AddCapabilities`
  <https://github.com/onlyati/quadlet-lsp/issues/128>
- Add formatting <https://github.com/onlyati/quadlet-lsp/pull/124>

### Bugfixes

- QSR010 only accepted ports without trailing `/tcp` or `/udp`
  <https://github.com/onlyati/quadlet-lsp/issues/133>
- Code completion was also generated in comment lines too
  <https://github.com/onlyati/quadlet-lsp/issues/122>
- The `quadlet-lsp check` CLI command wasn't aware of drop-ins directory
  <https://github.com/onlyati/quadlet-lsp/issues/99>
- The `PullAll` command wasn't aware of drop-ins directory
  <https://github.com/onlyati/quadlet-lsp/issues/118>

## v0.5.0

### Feature

- Add dropins file support generally
  <https://github.com/onlyati/quadlet-lsp/pull/95>
  <https://github.com/onlyati/quadlet-lsp/pull/100>
  <https://github.com/onlyati/quadlet-lsp/pull/101>
  <https://github.com/onlyati/quadlet-lsp/pull/119>
- New syntax rule: QSR025 <https://github.com/onlyati/quadlet-lsp/pull/117>

### Bugfixes

- Fix regex in QSR008 and QSR009
  <https://github.com/onlyati/quadlet-lsp/pull/114>

### Repo

- Add tool to generate rpm and deb files at release
  <https://github.com/onlyati/quadlet-lsp/pull/102>
- Add mise tool <https://github.com/onlyati/quadlet-lsp/pull/108>
  <https://github.com/onlyati/quadlet-lsp/pull/109>
  <https://github.com/onlyati/quadlet-lsp/pull/110>
- Add spec file for Fedora copr
  <https://github.com/onlyati/quadlet-lsp/pull/112>
- Add more ways to download the language server
  <https://github.com/onlyati/quadlet-lsp/pull/113>

## v0.4.0

### Features

- 3rd party extension for Zed editor (#61)
- Hover explanation for systemd specifiers (#57)
- Hover explanation for `UserNS` (#78)
- Hover explanation for `Volume` (#84)
- Hover explanation for `Secret` (#85)
- `QSR022`: validate path with systemd specifier (#57)
- `QSR023`: validate systemd specifiers (#57)
- `QSR024`: warn for forbidden properties in `[Service]` (#77)
- Completion for systemd specifiers (#72)
- The `go definition` and `go references` works with template files (#75)
- Rule disabling on file basis (#81)

### Bugfixes

- `QSR003` pointed the error to the previous line (#64)
- `QSR021` accept all accept all systemd unit types (#65)
- `QSR008`, `QSR009`: fix naming convention checking (#66)
- The `;` also count as valid comment character besides `#` (#76)

### Other

- Improve on documentation (#86)
- Make QoL changes on github repo (#90) (#91)

## New Contributors

- @mufeedali made their first contribution in
  <https://github.com/onlyati/quadlet-lsp/pull/61>

## v0.3.1

### Bugfixes

- The '@' character caused false positive checks in QSR021
  <https://github.com/onlyati/quadlet-lsp/pull/53>
- DefaultInstance was missing in Install section
  <https://github.com/onlyati/quadlet-lsp/pull/53>
- The value of Exec property can be split to multiple line and multi line was
  handled individually <https://github.com/onlyati/quadlet-lsp/pull/52>
- Fix fully qualified syntax checking
  <https://github.com/onlyati/quadlet-lsp/pull/51>

## v0.3.0

### Features

- New syntax validation: checking container, volume, pod and network name
  <https://github.com/onlyati/quadlet-lsp/pull/33>
- Set properties of Podman 5.6.0
  <https://github.com/onlyati/quadlet-lsp/pull/34>
- Build files has static completion and new template
  <https://github.com/onlyati/quadlet-lsp/pull/35>
- Add completion for `Unit` and `Service` sections and validate automatic
  dependency translation <https://github.com/onlyati/quadlet-lsp/pull/36>
- Modify syntax rule, from Podman 5.6.0, environment variable can be specified
  without value <https://github.com/onlyati/quadlet-lsp/pull/39>
- Language server listing the exposed ports based on the image. But if image is
  not pulled, it cannot read. From now it gives an information message if
  exposed port is not found and could not check all images
  <https://github.com/onlyati/quadlet-lsp/pull/40>
- Add new language server commands: list jobs and pull all image
  <https://github.com/onlyati/quadlet-lsp/pull/41>

### Bugfixes

- Label, Annotation and Environment variables only accepted one style
  specification. Syntax check has been updated to accept all possible variation
  <https://github.com/onlyati/quadlet-lsp/pull/32>
- Invalid property was checking the commented lines too
  <https://github.com/onlyati/quadlet-lsp/pull/32>

## v0.2.1

### Bugfixes

- [fix] If .quadletrc.json does not exists, exit with error by @onlyati in
  <https://github.com/onlyati/quadlet-lsp/pull/29>

## v0.2.0

### Completions

- Add completion for PublishPort <https://github.com/onlyati/quadlet-lsp/pull/2>
- Add user namespace keep id completion
  <https://github.com/onlyati/quadlet-lsp/pull/4>
- Add `new.Image` completion <https://github.com/onlyati/quadlet-lsp/pull/12>
- Add completion for PublishPort in pod files
  <https://github.com/onlyati/quadlet-lsp/pull/22>

### Syntax checking

- Add disable function based on file in working directory
  <https://github.com/onlyati/quadlet-lsp/pull/23>
- Basic syntax checker implementation
  <https://github.com/onlyati/quadlet-lsp/pull/6>
  - QSR001, QSR002, QSR003 implemented altogether
  - QSR001 - Missing section header
  - QSR002 - Unfinished line
  - QSR003 - Invalid property
- QSR004 - Not fully qualified image
  <https://github.com/onlyati/quadlet-lsp/pull/8>
- QSR005 - Invalid value of AutoUpdate
  <https://github.com/onlyati/quadlet-lsp/pull/10>
- QSR006 - Image file does not exists
  <https://github.com/onlyati/quadlet-lsp/pull/11>
- QSR007 - Invalid format of Environment variable specification
  <https://github.com/onlyati/quadlet-lsp/pull/13>
- QSR008 - Invalid format of Annotation specification
  <https://github.com/onlyati/quadlet-lsp/pull/14>
- QSR009 - Invalid format of Label specification
  <https://github.com/onlyati/quadlet-lsp/pull/15>
- QSR010 - Invalid port number is used at published port
  <https://github.com/onlyati/quadlet-lsp/pull/16>
- QSR011 - Port is not exposed in image
  <https://github.com/onlyati/quadlet-lsp/pull/17>
- QSR012 - Invalid format of Secret specification
  <https://github.com/onlyati/quadlet-lsp/pull/26>
- QSR013 - Volume file does not exists
  <https://github.com/onlyati/quadlet-lsp/pull/18>
- QSR014 - Network file does not exists
  <https://github.com/onlyati/quadlet-lsp/pull/19>
- QSR015 - Invalid format of Volume specification
  <https://github.com/onlyati/quadlet-lsp/pull/27>
- QSR016 - Invalid value of UserNS specification
  <https://github.com/onlyati/quadlet-lsp/pull/28>
- QSR017 - Pod file does not exists
  <https://github.com/onlyati/quadlet-lsp/pull/20>
- QSR018 - Container cannot have PublishPort if belongs to a pod
  <https://github.com/onlyati/quadlet-lsp/pull/21>
- QSR019 - Container cannot have Network if belongs to a pod
  <https://github.com/onlyati/quadlet-lsp/pull/21>

### Other

- Add feature to run syntax checker from CLI and pass over directory or file
  <https://github.com/onlyati/quadlet-lsp/pull/23>
- Make it version aware <https://github.com/onlyati/quadlet-lsp/pull/23>
  - From version 5.4.0
  - Podman version can also be set in `quadletrc.json` file
