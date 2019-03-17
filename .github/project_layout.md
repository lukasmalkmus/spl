# Project Layout

This document describes the project layout of this repository. It is derived
from the [Standard Go Project Layout].

Additional background information is available at [Go Project Layout].

## Go Directories

### `/cmd`

Main applications for this project.

The directory name for each application should match the name of the executable
(e.g., `/cmd/myapp`).

If code can be imported and used in other projects, then it should live in the
`/pkg` directory. If the code is not or should not be reusable, that code should
live in the `/internal` directory.

It's common to have a small main function that imports and invokes the code from
the `/internal` and `/pkg` directories.

### `/internal`

Private application and library code.

Put actual application code in the `/internal/app` directory (e.g.,
`/internal/app/myapp`) and the code shared by those apps in the `/internal/pkg`
directory (e.g., `/internal/pkg/myprivlib`).

### `/pkg`

Library code that's safe to use by external applications (e.g.,
`/pkg/mypubliclib`). Other projects will import these libraries expecting them
to work!

### `/vendor`

Application dependencies managed by [go mod].

## Common Application Directories

### `/configs`

Configuration file templates or default configs.

<!-- Links -->
[Standard Go Project Layout]: https://github.com/golang-standards/project-layout
[Go Project Layout]: https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2
[go mod]: https://golang.org/cmd/go/#hdr-Module_maintenance
