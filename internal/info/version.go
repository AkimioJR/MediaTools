package info

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
	SupportDesktopMode: true,
	CommitHash:         commitHash,
	BuildTime:          parseBuildTime(buildTime),
	GoVersion:          runtime.Version(),
	OS:                 runtime.GOOS,
	Arch:               runtime.GOARCH,
}

func (v VersionInfo) String() string {
	return fmt.Sprintf(
		"版本: %s\n编译时间: %s\nGit 提交: %s\nGo 版本: %s\n操作系统: %s\n架构: %s",
		v.AppVersion,
		v.BuildTime,
		v.CommitHash,
		v.GoVersion,
		v.OS,
		v.Arch,
	)
}

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unkown"
	} else {
		return t.Local().Format(time.DateTime)
	}
}
