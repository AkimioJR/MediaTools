package config

import "time"

var (
	appVersion string = "v0.0.1"
	commitHash string = "Unkown"
	buildDate  string = "Unkown"
)

type VersionInfo struct {
	AppVersion string `json:"app_version"` // 程序版本号
	CommitHash string `json:"commit_hash"` // GIt Commit Hash
	BuildDate  string `json:"build_date"`  // 编译时间
	GoVersion  string `json:"go_version"`  // 编译 Golang 版本
	OS         string `json:"os"`          // 操作系统
	Arch       string `json:"arch"`        // 架构
}

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unkown"
	} else {
		return t.Local().Format(time.RFC822Z)
	}
}
