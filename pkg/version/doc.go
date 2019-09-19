// Package version provides variables which, when used correctly, provide
// build information about the application referencing this package. At compile
// time the appropriate linker flags must be passed:
//
// 	go build -ldflags "-X github.com/lukasmalkmus/spl/pkg/version.release=1.0.0"
//
// Adapt the flags for all other exported variables. Eventually use vendored
// version.
package version
