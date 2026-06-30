// Package version holds the application version, overridable at build time via
//
//	-ldflags "-X ncp-nuke/pkg/version.Version=1.2.3"
package version

// Version is the current application version (without a leading "v").
var Version = "1.0.0"

// Repo is the GitHub owner/name used for update checks.
const Repo = "enbraining/ncp-nuke"
