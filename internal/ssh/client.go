package ssh

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/fatih/color"
	"quickship/internal/config"
)

type Client struct {
	host   string
	user   string
	client *ssh.Client
	color  *color.Color
}

func NewClient(host, user string, colorAttr color.Attribute) (*Client, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK not set")
	}

	agentConn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}

	agentClient := agent.NewClient(agentConn)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		host:   host,
		user:   user,
		client: client,
		color:  color.New(colorAttr),
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) ExecuteCommand(cmd string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if err := agent.RequestAgentForwarding(session); err != nil {
		return fmt.Errorf("failed to request agent forwarding: %w", err)
	}

	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	if err := session.Start(cmd); err != nil {
		return err
	}

	go c.streamOutput(stdout, false)
	go c.streamOutput(stderr, true)

	return session.Wait()
}

func (c *Client) streamOutput(r io.Reader, isErr bool) {
	scanner := bufio.NewScanner(r)
	prefix := fmt.Sprintf("[%s]", c.host)
	for scanner.Scan() {
		c.color.Printf("%s %s\n", prefix, scanner.Text())
	}
}

func (c *Client) DeployProject(project config.Project, gitOnly bool) error {
	deployScript := ""
	if !gitOnly {
		deployScript = project.Scripts["deploy"]
	}

	script := fmt.Sprintf(`
if [ ! -d "%s" ]; then
    mkdir -p %s 2>/dev/null
    if [ $? -ne 0 ]; then
        echo ""
        echo "ERROR: Permission denied - cannot create directory %s"
        echo ""
        echo "Please run the following commands on the server first:"
        echo "  sudo mkdir -p %s"
        echo "  sudo chown -R $USER:$USER %s"
        echo ""
        echo "Or configure passwordless sudo for this user:"
        echo "  echo '$USER ALL=(ALL) NOPASSWD: /bin/mkdir, /bin/chown' | sudo tee /etc/sudoers.d/$USER"
        echo ""
        exit 1
    fi
fi
if [ ! -w "%s" ]; then
    echo ""
    echo "ERROR: No write permission to directory %s"
    echo ""
    echo "Please run on the server:"
    echo "  sudo chown -R $USER:$USER %s"
    echo ""
    exit 1
fi
cd %s
if [ ! -d ".git" ]; then
    git clone %s .
else
    git pull
fi
%s
`, project.Path, project.Path,
		project.Path, project.Path, project.Path,
		project.Path, project.Path, project.Path,
		project.Path, project.Repo, deployScript)

	return c.ExecuteCommand(script)
}

