package main

import (
	"MediaTools/internal/info"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func parseVersion(versionStr string) (int, int, int) {
	parts := strings.Split(strings.Replace(versionStr, "v", "", 1), ".")
	major := 0
	minor := 0
	patch := 0
	if len(parts) >= 3 {
		fmt.Sscanf(parts[0], "%d", &major)
		fmt.Sscanf(parts[1], "%d", &minor)
		fmt.Sscanf(parts[2], "%d", &patch)
	}
	return major, minor, patch
}

func getVersion(isRelease bool) string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		return "PreRelease-0.0.0-" + getGitCommitHash(true)
	}
	major, minor, patch := parseVersion(strings.TrimSpace(string(out)))
	if isRelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	} else {
		return fmt.Sprintf("PreRelease-%d.%d.%d-%s", major, minor, patch+1, getGitCommitHash(true))
	}
}

func getGitCommitHash(isShort bool) string {
	args := []string{"rev-parse"}
	if isShort {
		args = append(args, "--short")
	}
	args = append(args, "HEAD")

	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		panic("获取 git commit 失败: " + err.Error() + "\n" + string(out))
	}
	return strings.ReplaceAll(string(out), "\n", "")
}

func getServerName() string {
	name := info.ProjectName + "-" + *targetOS + "-" + *targetArch
	if *targetOS == "windows" {
		name += ".exe"
	}
	return name
}

func getDesktopName() string {
	name := filepath.Join("build", "bin", info.ProjectName)
	switch *targetOS {
	case "windows":
		name += ".exe"
	case "darwin":
		name += ".app"
	}
	return name
}

func needBuildFrontend() bool {
	if _, err := os.Stat("web/dist/index.html"); err == nil {
		return false
	}
	return true
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

	output, err = exec.Command("swag", "init").CombinedOutput()
	if err != nil {
		fmt.Println("更新 Swagger 文档失败: \n" + string(output))
		panic(err.Error())
	}
	fmt.Println("更新 Swagger 文档成功🎉")

	infoFlags := []string{
		"-X", "MediaTools/internal/info.appVersion=" + getVersion(*isRelease),
		"-X", "MediaTools/internal/info.buildTime=" + time.Now().Format(time.RFC3339),
		"-X", "MediaTools/internal/info.commitHash=" + getGitCommitHash(false),
	}
	ldFlags := []string{
		"-s",
		"-w",
	}

	var cmd *exec.Cmd
	if *desktopMode {
		platformArgs := []string{"-platform", *targetOS + "/" + *targetArch}
		args := append([]string{"build", "-skipbindings", "-ldflags", strings.Join(infoFlags, " ")}, platformArgs...)
		if !*buildFrontend {
			args = append(args, "-s")
		}
		args = append(args, ".")
		fmt.Println("执行命令: wails", strings.Join(args, " "))
		print("\n\n")
		cmd = exec.Command("wails", args...)

	} else {
		if *buildFrontend {
			fmt.Println("开始构建前端...")
			err = buildWeb()
			if err != nil {
				panic("构建前端失败: \n" + err.Error())
			}
			fmt.Println("构建前端成功🎉")
		}

		fmt.Println("设置 GOOS 和 GOARCH 成功🎉")

		args := []string{"build", "-o", getServerName()}
		args = append(args, "-ldflags", strings.Join(append(ldFlags, infoFlags...), " "), ".")
		fmt.Println("执行命令: go", strings.Join(args, " "))
		print("\n\n")
		cmd = exec.Command("go", args...)
		cmd.Env = append(os.Environ(), "GOOS"+"="+*targetOS, "GOARCH"+"="+*targetArch)
	}

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("构建命令输出:")
		fmt.Println(string(output))
		panic("构建失败: " + err.Error())
	} else {
		fmt.Println("构建成功！🎉🎉🎉")
	}

	if *outputName != "" {
		if *desktopMode {
			err = os.Rename(getDesktopName(), *outputName)
		} else {
			err = os.Rename(getServerName(), *outputName)
		}
		if err != nil {
			fmt.Println("重命名输出文件失败: " + err.Error())
		} else {
			fmt.Println("重命名输出文件成功！🎉🎉🎉 输出文件: " + *outputName)
		}
	}
}

var (
	targetOS      = flag.String("os", runtime.GOOS, "目标操作系统\nTarget operating system")
	targetArch    = flag.String("arch", runtime.GOARCH, "目标架构\nTarget architecture")
	desktopMode   = flag.Bool("desktop", false, "桌面模式\nDesktop mode")
	buildFrontend = flag.Bool("build-frontend", needBuildFrontend(), fmt.Sprintf("强制构建前端(默认: %v)\nForce build frontend(Defaults: %v)", needBuildFrontend(), needBuildFrontend()))
	outputName    = flag.String("output", "", "输出文件名(默认: 根据 os/arch 和 desktop 自动生成)\nOutput file name(Defaults: auto generate by os/arch and desktop)")
	isRelease     = flag.Bool("release", false, "发布模式\nRelease mode")

	showVersionFlag = flag.Bool("version", false, "显示版本信息\nShow version information")
)

func init() {
	flag.Parse()
}

func main() {
	if *showVersionFlag {
		fmt.Println(getVersion(*isRelease))
		return
	}
	build()
}
