package version

import (
	"fmt"
	"runtime"
	"time"
)

const tsFormat = "2006-01-02_15:04:05"

var (
	buildTime = "-"
	commit    = "-"
	release   = "-"
	goVersion = runtime.Version()
	user      = "-"
)

// BuildTime is a timestamp of the moment when the binary was built.
func BuildTime() string {
	ts, err := time.ParseInLocation(tsFormat, buildTime, time.UTC)
	if err != nil {
		return buildTime
	}
	return ts.UTC().Format(time.RFC3339)
}

// Commit is the last commit hash of the source repository at the moment
// the binary was built.
func Commit() string {
	return commit
}

// Release is the semantic version of the current build.
func Release() string {
	return release
}

// GoVersion is the go version the build utilizes.
func GoVersion() string {
	return goVersion
}

// User is the username of the user who performed the build.
func User() string {
	return user
}

// String returns the version and build information.
func String() string {
	return fmt.Sprintf("build=%s, commit=%s, release=%s, go=%s, user=%s",
		buildTime, commit, release, goVersion, user)
}
