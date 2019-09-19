# lukasmalkmus/spl

> A toolchain for a simple programming language, inspired by the Go toolchain. - by **[Lukas Malkmus]**

[![Build Status][build_badge]][build]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![GoDoc][docs_badge]][docs]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]
[![License Status][license_status_badge]][license_status]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

*spl* is a compiler and toolchain for a **s**imple **p**rogramming **l**anguage
I had to build a compiler in an university course on building compilers. The
design is heavily inspired by the packages of the [Go][go] language and features
a handwritten parser as well as the two excellent books from [Thorsten Ball] on
building an [interpreter][interpreter_book] and a [compiler][compiler_book]. 

The languages specification can be reviewed [here][spl_spec] (in german).
A machine translated one in english (yes, I'm a lazy fuck) can be reviewed
[here][docs].

## Usage

### Installation

The easiest way to run *spl* is by grabbing the latest standalone binary from
the [release page][release].

This project uses native [go mod] support for vendoring and requires a working
`go` toolchain installation when installing via `go get` or from source.

#### Install using `go get`

```bash
GO111MODULE=on go install github.com/lukasmalkmus/spl/cmd/spl
```

#### Install from source

```bash
git clone https://github.com/lukasmalkmus/spl.git
cd spl
make # Build production binary
make install # Build and install binary into $GOPATH
```

#### Validate installation

The installation can be validated by running `spl version` in the terminal.

### Configuration

*spl* is a [Twelve Factor Application] and can be configured by either
configuration file, the environment or command line flags. It provides a basic
*help flag* `--help` which prints out application and configuration help. See
[using the application](#using-the-application).

Configuration files are [TOML] formatted:

```toml
[format]
indent = 8
```

Sections which are in TOML indicated by `[...]` are mapped to their respective
environment variables by seperating sections and values with an underscore `_`.
However, they are prefixed by the application name:

```bash
export SPL_FORMAT_INDENT=8
```

The same is true for command line flags but they are separated by a dot `.` and
not prefixed:

```bash
spl --format.indent=8
```

Configuration values **without** a default value must be set explicitly.

The application itself can echo out its configuration by calling the `config`
command:

```bash
spl config > spl.toml
```

Configuration priority from lowest to highest is like presented above:
Configuration file, environment, command line option (flag).

<details>
<summary>Click to expand default configuration file:</summary>

```toml
# SPL COMPILER TOOLCHAIN CONFIGURATION

# Source code formatter configuration.
[format]
# Indentation width used.
indent = 4

```

</details>

### Using the application

```bash
spl [flags] [commands]
```

Help on flags and commands:

```bash
spl --help
```

## Contributing

Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

More information about the project layout is documented
[here](.github/project_layout.md).

## License

© Lukas Malkmus, 2019

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

[![License Status Large][license_status_large_badge]][license_status_large]

<!-- Links -->
[Lukas Malkmus]: https://github.com/lukasmalkmus
[go]: https://golang.org
[Thorsten Ball]: https://github.com/mrnugget
[interpreter_book]: https://interpreterbook.com
[compiler_book]: https://compilerbook.com
[spl_spec]: https://homepages.thm.de/~hg52/lv/compiler/praktikum/SPL-1.2.html
[go mod]: https://golang.org/cmd/go/#hdr-Module_maintenance
[Twelve Factor Application]: https://12factor.net
[TOML]: https://github.com/toml-lang/toml

<!-- Badges -->
[build]: https://travis-ci.com/lukasmalkmus/spl
[build_badge]: https://img.shields.io/travis/com/lukasmalkmus/spl.svg?style=flat-square
[coverage]: https://codecov.io/gh/lukasmalkmus/spl
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/spl.svg?style=flat-square
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/spl
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/spl?style=flat-square
[docs]: https://godoc.org/github.com/lukasmalkmus/spl
[docs_badge]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square
[release]: https://github.com/lukasmalkmus/spl/releases
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/spl.svg?style=flat-square
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/spl.svg?color=blue&style=flat-square
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fspl?ref=badge_shield
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Fspl.svg
[license_status_large]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fspl?ref=badge_large
[license_status_large_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Fspl.svg?type=large
