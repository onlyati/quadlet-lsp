# Quadlet syntax rules

<!-- toc -->

- [`QSR001` - Missing section header](#qsr001---missing-section-header)
- [`QSR002` - Unfinished line](#qsr002---unfinished-line)
- [`QSR003` - Invalid property](#qsr003---invalid-property)
- [`QSR004` - Image name is not fully qualified](#qsr004---image-name-is-not-fully-qualified)
- [`QSR005` - Invalid value of AutoUpdate](#qsr005---invalid-value-of-autoupdate)
- [`QSR006` - Image file does not exists](#qsr006---image-file-does-not-exists)
- [`QSR007` - Invalid format of Environment variable](#qsr007---invalid-format-of-environment-variable)

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

Use fully qualified image name instead:

```ini
[Container]
Image=docker.io/library/debian:bookworm-slim

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

> Invalid format of Environment variable specification

**Explanation**

Environment variable must be specified as key-value pairs without having space
before or after the `=` sign.

Examples:

```ini
Environment=FOO=BAR   <-- Correct
Environment=FOO       <-- Incorrect
Environment=FOO = BAR <-- Incorrect
```
