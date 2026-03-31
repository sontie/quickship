package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"quickship/internal/config"
	"quickship/internal/executor"
)

func Init() error {
	template := `version: "1.0"

hosts:
  ali-dev:
    addr: "192.168.1.100:22"
    user: "deploy"
  ali-prod:
    addr: "192.168.1.101:22"
    user: "deploy"

projects:
  - name: "api-server"
    repo: "git@github.com:user/api-server.git"
    path: "/opt/api-server"
    envs:
      dev: ["ali-dev"]
      prod: ["ali-prod"]
    scripts:
      deploy: |
        docker-compose down
        docker-compose up -d --build
`
	return os.WriteFile("qship.yaml", []byte(template), 0644)
}

func Check() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}
	sshDir := filepath.Join(home, ".ssh")

	allPassed := true

	// Step 1: 检查 SSH 密钥是否存在
	fmt.Println("Step 1: Checking SSH key...")
	keyPath := ""
	candidates := []string{
		filepath.Join(sshDir, "id_ed25519"),
		filepath.Join(sshDir, "id_rsa"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			keyPath = c
			break
		}
	}
	if keyPath == "" {
		fmt.Println("  ✗ No SSH key found (~/.ssh/id_ed25519 or ~/.ssh/id_rsa)")
		fmt.Println("  → Run: ssh-keygen -t ed25519 -C \"your@email.com\"")
		allPassed = false
	} else {
		fmt.Printf("  ✓ SSH key:   %s\n", keyPath)
	}

	// Step 2: 检查 SSH Agent 是否运行
	fmt.Println("Step 2: Checking SSH Agent...")
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		fmt.Println("  ✗ SSH_AUTH_SOCK is not set (SSH Agent not running)")
		fmt.Println("  → Run: eval $(ssh-agent)")
		allPassed = false
	} else {
		fmt.Printf("  ✓ SSH Agent: %s\n", sock)
	}

	// Step 3: 检查 Agent 是否已加载密钥
	fmt.Println("Step 3: Checking loaded keys...")
	var keyCount int
	if sock != "" {
		cmd := exec.Command("ssh-add", "-l")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("  ✗ No keys loaded in SSH Agent")
			addTarget := "~/.ssh/id_ed25519"
			if keyPath != "" {
				addTarget = keyPath
			}
			fmt.Printf("  → Run: ssh-add %s\n", addTarget)
			allPassed = false
		} else {
			lines := strings.TrimSpace(string(output))
			keyCount = len(strings.Split(lines, "\n"))
			fmt.Printf("  ✓ Keys loaded: %d key(s)\n", keyCount)
		}
	} else {
		fmt.Println("  - Skipped (SSH Agent not running)")
	}

	// Step 4: 建议 ~/.ssh/config 持久化配置
	fmt.Println("Step 4: Checking SSH config (optional)...")
	configPath := filepath.Join(sshDir, "config")
	configContent, err := os.ReadFile(configPath)
	hasAddKeysToAgent := err == nil && strings.Contains(string(configContent), "AddKeysToAgent")
	if hasAddKeysToAgent {
		fmt.Println("  ✓ SSH config has AddKeysToAgent configured")
	} else {
		fmt.Println("  ⓘ Tip: Add the following to ~/.ssh/config to avoid running ssh-add after each reboot:")
		fmt.Println("")
		fmt.Println("    Host *")
		fmt.Println("      AddKeysToAgent yes")
		if runtime.GOOS == "darwin" {
			fmt.Println("      UseKeychain yes")
		}
		fmt.Println("")
	}

	// 汇总结果
	fmt.Println("---")
	if allPassed {
		fmt.Println("✓ SSH environment is ready!")
	} else {
		return fmt.Errorf("SSH environment is not ready. Please fix the issues above and re-run: qship check")
	}
	return nil
}

func Auth(host string) error {
	cmd := exec.Command("ssh-copy-id", host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func List(cfg *config.Config) {
	fmt.Println("Hosts:")
	for name, host := range cfg.Hosts {
		fmt.Printf("  %s: %s@%s\n", name, host.User, host.Addr)
	}

	fmt.Println("\nProjects:")
	for _, proj := range cfg.Projects {
		fmt.Printf("  %s (%s)\n", proj.Name, proj.Repo)
		for env, hosts := range proj.Envs {
			fmt.Printf("    %s: %s\n", env, strings.Join(hosts, ", "))
		}
	}
}

func Deploy(env string, cfg *config.Config) error {
	fmt.Printf("Deploying to %s environment...\n", env)
	return executor.Deploy(env, cfg)
}

func Exec(cmd string, hosts []string, cfg *config.Config) error {
	return executor.ExecCommand(hosts, cmd, cfg)
}
