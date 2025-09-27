package version

import (
	"runtime"
	"time"
)

var (
	appVersion string = "v0.0.1"
	commitHash string = "Unkown"
	buildTime  string = "Unkown"
)

var Version = VersionInfo{
	AppVersion:         appVersion,
	SupportDesktopMode: supportDesktopMode,
	CommitHash:         commitHash,
	BuildTime:          parseBuildTime(buildTime),
	GoVersion:          runtime.Version(),
	OS:                 runtime.GOOS,
	Arch:               runtime.GOARCH,
}

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unkown"
	} else {
		return t.Local().Format(time.DateTime)
	}
}
