package info

import "time"

type PlatformInfo struct {
	OS   string `json:"os"`   // 操作系统
	Arch string `json:"arch"` // 架构
}
type VersionInfo struct {
	PlatformInfo
	AppVersion         string `json:"app_version"`          // 程序版本号
	SupportDesktopMode bool   `json:"support_desktop_mode"` // 是否支持桌面模式
	CommitHash         string `json:"commit_hash"`          // GIt Commit Hash
	BuildTime          string `json:"build_time"`           // 编译时间
	GoVersion          string `json:"go_version"`           // 编译 Golang 版本
}

type RuntimeAppStatusInfo struct {
	PlatformInfo
	DesktopMode bool      `json:"desktop_mode"` // 是否桌面模式
	IsDev       bool      `json:"is_dev"`       // 是否开发者模式
	Port        uint16    `json:"port"`         // 访问端口
	BootTime    time.Time `json:"boot_time"`    // 启动时间
}
