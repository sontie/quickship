package main

import (
	"fmt"
	"os"
	"strings"

	"quickship/cmd"
	"quickship/internal/config"
)

const Version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Printf("QuickShip v%s\n", Version)
		return

	case "init":
		if err := cmd.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Created deploy.yaml")

	case "check":
		if err := cmd.Check(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "auth":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: qship auth <host>")
			os.Exit(1)
		}
		if err := cmd.Auth(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "list":
		cfg, err := config.LoadConfig("deploy.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		cmd.List(cfg)

	case "deploy":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: qship deploy <env>")
			os.Exit(1)
		}
		cfg, err := config.LoadConfig("deploy.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := cmd.Deploy(os.Args[2], cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Deployment completed")

	case "exec":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: qship exec \"<command>\" [host1,host2,...]")
			os.Exit(1)
		}
		cfg, err := config.LoadConfig("deploy.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		hosts := []string{}
		if len(os.Args) > 3 {
			hosts = strings.Split(os.Args[3], ",")
		} else {
			for name := range cfg.Hosts {
				hosts = append(hosts, name)
			}
		}

		if err := cmd.Exec(os.Args[2], hosts, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`QuickShip - SSH Agent Forwarding Deployment Tool

Usage:
  qship version                 Show version information
  qship init                    Generate deploy.yaml template
  qship check                   Check SSH agent status
  qship auth <host>             Copy SSH key to host
  qship list                    List hosts and projects
  qship deploy <env>            Deploy to environment
  qship exec "<cmd>" [hosts]    Execute command on hosts`)
}
