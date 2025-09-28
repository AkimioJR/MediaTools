package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func getVersion(isDev bool) string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		if !isDev {
			panic("获取版本号失败: " + err.Error())
		}
		return "dev-" + getGitCommitHash(true)
	}
	return strings.ReplaceAll(string(out), "\n", "")
}
func getGitCommitHash(isShort bool) string {
	args := []string{"rev-parse"}
	if isShort {
		args = append(args, "--short")
	}
	args = append(args, "HEAD")

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		panic("获取 git commit 失败: " + err.Error())
	}
	return strings.ReplaceAll(string(out), "\n", "")
}

func getTimeStr() string {
	return time.Now().Format(time.RFC3339)
}

func getOutputName() string {
	baseName := "MediaTools-" + targetOS + "-" + targetArch
	if runtime.GOOS == "windows" {
		baseName += ".exe"
	}
	return baseName
}

var (
	appVersion  string
	buildTime   string
	commitHash  string
	desktopMode bool
	targetOS    string
	targetArch  string
	outputName  string

	showVersion = false
)

func init() {
	flag.StringVar(&appVersion, "version", getVersion(true), "应用版本")
	flag.StringVar(&buildTime, "build-time", getTimeStr(), "构建时间")
	flag.StringVar(&commitHash, "commit-hash", getGitCommitHash(false), "Git 提交哈希值")
	flag.BoolVar(&desktopMode, "desktop", false, "编译桌面模式")
	flag.StringVar(&targetOS, "os", runtime.GOOS, "目标操作系统")
	flag.StringVar(&targetArch, "arch", runtime.GOARCH, "目标架构")
	flag.StringVar(&outputName, "output", getOutputName(), "输出文件名")

	flag.BoolVar(&showVersion, "version-info", false, "显示版本信息并退出")

	flag.Parse()
}

func showInfo() {
	println(strings.Repeat("=", 70))
	println("应用版本:", appVersion)
	println("构建时间:", buildTime)
	println("Git 提交哈希值:", commitHash)
	println("目标操作系统:", targetOS)
	println("目标架构:", targetArch)
	println("输出文件名:", outputName)
	if desktopMode {
		println("编译模式: 桌面模式")
	} else {
		println("编译模式: 服务器模式")
	}
	println(strings.Repeat("=", 70))
	print("\n\n")
}

func build() {
	err := exec.Command("go", "mod", "download").Run()
	if err != nil {
		panic("下载依赖失败: " + err.Error())
	}
	fmt.Println("下载依赖成功🎉")

	err = exec.Command("go", "env", "-w", "GOOS="+targetOS).Run()
	if err != nil {
		panic("设置 GOOS 失败: " + err.Error())
	}
	err = exec.Command("go", "env", "-w", "GOARCH="+targetArch).Run()
	if err != nil {
		panic("设置 GOARCH 失败: " + err.Error())
	}
	fmt.Println("设置 GOOS 和 GOARCH 成功🎉")

	args := []string{"build", "-o", outputName}
	if !desktopMode {
		args = append(args, "-tags=onlyServer")
	}
	ldFlags := []string{
		"-s",
		"-w",
		"-X", "MediaTools/internal/version.appVersion=" + appVersion,
		"-X", "MediaTools/internal/version.buildTime=" + buildTime,
		"-X", "MediaTools/internal/version.commitHash=" + commitHash,
	}
	if targetOS == "windows" && desktopMode {
		ldFlags = append(ldFlags, "-H", "windowsgui")
	}

	args = append(args, "-ldflags", strings.Join(ldFlags, " "), ".")
	// fmt.Println("执行命令: go", strings.Join(args, " "))
	print("\n\n")

	err = exec.Command("go", args...).Run()
	if err != nil {
		panic("构建失败: " + err.Error())
	} else {
		fmt.Println("构建成功！🎉🎉🎉")
	}
}

func main() {
	if showVersion {
		fmt.Println(appVersion)
		return
	}
	showInfo()
	build()
}
