package cmd

import (
	"fmt"
	"os"
	"os/exec"
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
	return os.WriteFile("deploy.yaml", []byte(template), 0644)
}

func Check() error {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return fmt.Errorf("SSH_AUTH_SOCK not set. Run: eval $(ssh-agent) && ssh-add")
	}
	fmt.Printf("✓ SSH Agent: %s\n", sock)

	cmd := exec.Command("ssh-add", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("no keys in agent. Run: ssh-add")
	}
	fmt.Printf("✓ Keys loaded:\n%s", output)
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
