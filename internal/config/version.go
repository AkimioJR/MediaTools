package config

import "time"

var (
	appVersion string = "v0.0.1"
	commitHash string = "Unkown"
	buildDate  string = "Unkown"
)

type VersionInfo struct {
	AppVersion string // 程序版本号
	CommitHash string // GIt Commit Hash
	BuildDate  string // 编译时间
	GoVersion  string // 编译 Golang 版本
	OS         string // 操作系统
	Arch       string // 架构
}

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unkown"
	} else {
		return t.Local().Format(time.RFC822Z)
	}
}
