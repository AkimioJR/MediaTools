package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func getVersion(isRelease bool) string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		if isRelease {
			panic("获取版本号失败: " + err.Error() + "\n" + string(out))
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

func getOutputName(isPackage bool) string {
	baseName := "MediaTools-" + targetOS + "-" + targetArch
	if runtime.GOOS == "windows" {
		baseName += ".exe"
	}
	if desktopMode && isPackage && runtime.GOOS == "darwin" {
		baseName += ".app"
	}
	return baseName
}

func needBuildFrontend() bool {
	if _, err := os.Stat("web/dist/index.html"); err == nil {
		return false
	}
	return true
}

var (
	appVersion    string
	buildTime     string
	commitHash    string
	desktopMode   bool
	buildFrontend bool
	targetOS      string
	targetArch    string
	outputName    string

	isRelease bool

	showVersion bool
)

func init() {
	flag.BoolVar(&isRelease, "release", false, "是否为发布版本 (default false)")
	flag.StringVar(&appVersion, "version", getVersion(false), "应用版本")
	flag.StringVar(&buildTime, "build-time", getTimeStr(), "构建时间")
	flag.StringVar(&commitHash, "commit-hash", getGitCommitHash(false), "Git 提交哈希值")
	flag.BoolVar(&desktopMode, "desktop", false, "编译桌面模式 (default false)")
	flag.BoolVar(&buildFrontend, "web", needBuildFrontend(), fmt.Sprintf("是否构建前端 (default %v)", needBuildFrontend()))
	flag.StringVar(&targetOS, "os", runtime.GOOS, "目标操作系统")
	flag.StringVar(&targetArch, "arch", runtime.GOARCH, "目标架构")
	flag.StringVar(&outputName, "output", getOutputName(false), "输出文件名")

	flag.BoolVar(&showVersion, "version-info", false, "显示版本信息并退出")

	flag.Parse()

	if isRelease && strings.HasPrefix(appVersion, "dev-") {
		appVersion = getVersion(isRelease)
	}
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
	println("是否构建前端:", strconv.FormatBool(buildFrontend))
	println("是否为发布版本:", strconv.FormatBool(isRelease))

	println(strings.Repeat("=", 70))
	print("\n\n")
}

func buildWeb() error {
	// 安装前端依赖
	cmd := exec.Command("pnpm", "install")
	cmd.Dir = "web" // 设置工作目录
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("安装前端依赖失败: \n%s", string(output))
	}

	// 构建前端
	cmd = exec.Command("pnpm", "build")
	cmd.Dir = "web" // 设置工作目录
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("构建前端失败: \n%s", string(output))
	}
	return nil
}

func build() {
	output, err := exec.Command("go", "mod", "download").CombinedOutput()
	if err != nil {
		fmt.Println("下载依赖失败: \n" + string(output))
		panic(err.Error())
	}
	fmt.Println("下载依赖成功🎉")

	infoFlags := []string{
		"-X", "MediaTools/internal/version.appVersion=" + appVersion,
		"-X", "MediaTools/internal/version.buildTime=" + buildTime,
		"-X", "MediaTools/internal/version.commitHash=" + commitHash,
	}
	ldFlags := []string{
		"-s",
		"-w",
	}

	var cmd *exec.Cmd
	if desktopMode {
		platformArgs := []string{"-platform", targetOS + "/" + targetArch}
		outputArgs := []string{"-o", outputName}
		args := append([]string{"build", "-skipbindings", "-ldflags", strings.Join(infoFlags, " ")}, append(platformArgs, outputArgs...)...)
		if !buildFrontend {
			args = append(args, "-s")
		}
		args = append(args, ".")
		fmt.Println("执行命令: wails", strings.Join(args, " "))
		print("\n\n")
		cmd = exec.Command("wails", args...)

	} else {
		if buildFrontend {
			fmt.Println("开始构建前端...")
			err = buildWeb()
			if err != nil {
				panic("构建前端失败: \n" + err.Error())
			}
			fmt.Println("构建前端成功🎉")
		}

		fmt.Println("设置 GOOS 和 GOARCH 成功🎉")

		args := []string{"build", "-o", outputName}
		args = append(args, "-ldflags", strings.Join(append(ldFlags, infoFlags...), " "), ".")
		fmt.Println("执行命令: go", strings.Join(args, " "))
		print("\n\n")
		cmd = exec.Command("go", args...)
		cmd.Env = append(os.Environ(), "GOOS"+"="+targetOS, "GOARCH"+"="+targetArch)
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("构建命令输出:")
		fmt.Println(string(output))
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
