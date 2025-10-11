package info

type VersionInfo struct {
	AppVersion         string `json:"app_version"`          // 程序版本号
	SupportDesktopMode bool   `json:"support_desktop_mode"` // 是否支持桌面模式
	CommitHash         string `json:"commit_hash"`          // GIt Commit Hash
	BuildTime          string `json:"build_time"`           // 编译时间
	GoVersion          string `json:"go_version"`           // 编译 Golang 版本
	OS                 string `json:"os"`                   // 操作系统
	Arch               string `json:"arch"`                 // 架构
}
