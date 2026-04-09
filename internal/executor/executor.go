package executor

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	"quickship/internal/config"
	"quickship/internal/ssh"
)

var colors = []color.Attribute{
	color.FgCyan,
	color.FgGreen,
	color.FgYellow,
	color.FgMagenta,
	color.FgBlue,
}

func Deploy(env string, gitOnly bool, rmAfter bool, cfg *config.Config) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(cfg.Projects))

	for i, project := range cfg.Projects {
		hosts, ok := project.Envs[env]
		if !ok {
			continue
		}

		for _, hostName := range hosts {
			host, ok := cfg.Hosts[hostName]
			if !ok {
				errChan <- fmt.Errorf("host %s not found", hostName)
				continue
			}

			wg.Add(1)
			go func(h config.Host, p config.Project, colorIdx int) {
				defer wg.Done()

				client, err := ssh.NewClient(h.Addr, h.User, colors[colorIdx%len(colors)])
				if err != nil {
					errChan <- err
					return
				}
				defer client.Close()

				if err := client.DeployProject(p, gitOnly, rmAfter); err != nil {
					errChan <- err
				}
			}(host, project, i)
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func Clean(env string, cfg *config.Config) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(cfg.Projects))

	for i, project := range cfg.Projects {
		hosts, ok := project.Envs[env]
		if !ok {
			continue
		}

		for _, hostName := range hosts {
			host, ok := cfg.Hosts[hostName]
			if !ok {
				errChan <- fmt.Errorf("host %s not found", hostName)
				continue
			}

			wg.Add(1)
			go func(h config.Host, p config.Project, colorIdx int) {
				defer wg.Done()

				client, err := ssh.NewClient(h.Addr, h.User, colors[colorIdx%len(colors)])
				if err != nil {
					errChan <- err
					return
				}
				defer client.Close()

				if err := client.CleanProject(p); err != nil {
					errChan <- err
				}
			}(host, project, i)
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecCommand(hostNames []string, cmd string, cfg *config.Config) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(hostNames))

	for i, hostName := range hostNames {
		host, ok := cfg.Hosts[hostName]
		if !ok {
			errChan <- fmt.Errorf("host %s not found", hostName)
			continue
		}

		wg.Add(1)
		go func(h config.Host, colorIdx int) {
			defer wg.Done()

			client, err := ssh.NewClient(h.Addr, h.User, colors[colorIdx%len(colors)])
			if err != nil {
				errChan <- err
				return
			}
			defer client.Close()

			if err := client.ExecuteCommand(cmd); err != nil {
				errChan <- err
			}
		}(host, i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
