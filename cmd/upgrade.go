package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const repo = "sontie/quickship"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func Upgrade(currentVersion string) error {
	// 获取最新版本
	latest, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check latest version: %w", err)
	}

	latestClean := strings.TrimPrefix(latest, "v")
	currentClean := strings.TrimPrefix(currentVersion, "v")

	if latestClean == currentClean {
		fmt.Printf("✓ Already up to date (v%s)\n", currentClean)
		return nil
	}

	fmt.Printf("Upgrading: v%s → v%s\n", currentClean, latestClean)

	// 构建下载 URL
	osName := runtime.GOOS
	arch := runtime.GOARCH

	filename := fmt.Sprintf("qship-%s-%s", osName, arch)
	if osName == "windows" {
		filename += ".exe"
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, latest, filename)

	fmt.Printf("Downloading %s...\n", filename)

	// 下载新版本
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// 写入临时文件
	tmpFile, err := os.CreateTemp("", "qship-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write file: %w", err)
	}
	tmpFile.Close()

	// 设置可执行权限
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// 获取当前二进制路径
	execPath, err := os.Executable()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// 替换自身
	if err := os.Rename(tmpPath, execPath); err != nil {
		// rename 跨设备可能失败，尝试 copy
		if err := copyFile(tmpPath, execPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("failed to replace binary: %w (try: sudo qship upgrade)", err)
		}
		os.Remove(tmpPath)
	}

	fmt.Printf("✓ Upgraded to v%s\n", latestClean)
	return nil
}

func getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
