package main

import (
	"MediaTools/internal/info"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Version struct {
	major     uint
	minor     uint
	patch     uint
	isRelease bool
}

func (v Version) String() string {
	if v.isRelease {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	} else {
		return fmt.Sprintf("PreRelease-%d.%d.%d-%s", v.major, v.minor, v.patch+1, getGitCommitHash(true))
	}
}

func ParseVersion(isRelease bool) *Version {
	var v Version

	v.isRelease = isRelease

	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return &v
	}

	parts := strings.Split(strings.Replace(string(out), "v", "", 1), ".")
	if len(parts) != 3 {
		return &v
	}
	fmt.Sscanf(parts[0], "%d", &v.major)
	fmt.Sscanf(parts[1], "%d", &v.minor)
	fmt.Sscanf(parts[2], "%d", &v.patch)

	return &v
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

func useWebkit2_41() bool {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var id, version string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}
	if err := scanner.Err(); err != nil {
		return false
	}

	if id != "ubuntu" {
		return false
	}

	// 解析版本号，提取主版本号
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return false
	}
	var major int
	_, err = fmt.Sscanf(parts[0], "%d", &major)
	if err != nil {
		return false
	}

	return major >= 24
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
		"-X", "MediaTools/internal/info.appVersion=" + ParseVersion(*isRelease).String(),
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
		args := append([]string{"build", "-clean", "-skipbindings", "-ldflags", strings.Join(infoFlags, " ")}, platformArgs...)
		if !*buildFrontend {
			args = append(args, "-s")
		}

		if *targetOS == "linux" && useWebkit2_41() {
			args = append(args, "-tags", "webkit2_41")
		}

		args = append(args, ".")
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

		cmd = exec.Command("go", args...)
		cmd.Env = append(os.Environ(), "GOOS"+"="+*targetOS, "GOARCH"+"="+*targetArch)
	}

	if *targetOS == "windows" && *targetArch == "arm64" {
		cmd.Env = append(os.Environ(), `CC=zig cc`)
	}

	fmt.Printf("构建命令: %s\n", cmd.String())
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic("构建失败: " + err.Error() + "\n" + "\n\n" + string(output))
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
		fmt.Println(ParseVersion(*isRelease).String())
		return
	}
	build()
}
