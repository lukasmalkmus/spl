package version

import (
	"fmt"
	"runtime"
)

var (
	// BuildTime is a timestamp of the moment when the binary was built.
	BuildTime = "-"

	// Commit is the last commit hash of the source repository at the moment
	// the binary was built.
	Commit = "-"

	// Release is the semantic version of the current build.
	Release = "-"

	// GoVersion is the go version the build utilizes.
	GoVersion = runtime.Version()

	// User is the username of the user who performed the build.
	User = "-"
)

// String returns the version and build information.
func String() string {
	return fmt.Sprintf("build=%s, commit=%s, release=%s, go=%s, user=%s",
		BuildTime, Commit, Release, GoVersion, User)
}
