package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/dikmit/gocron/slack"
	"github.com/mattn/go-shellwords"
)

type Job struct {
	Name     string
	Command  string
	Schedule string
	Cwd      string
	Slack    string
	Stdout   io.Writer
	Stderr   io.Writer
}

func (j *Job) Run() {
	if err := j.run(); err != nil {
		fmt.Fprintln(j.Stderr, err)
	}
}

func (j *Job) run() error {
	execCmd, err := createCommand(j.Command)
	if err != nil {
		return err
	}
	if j.Cwd != "" {
		path, err := filepath.Abs(j.Cwd)
		if err != nil {
			return fmt.Errorf("%s change dir error: %w", j.Name, err)
		}
		execCmd.Dir = path
	}

	buf := new(bytes.Buffer)
	outWtr := io.MultiWriter(j.Stdout, buf)
	errWtr := io.MultiWriter(j.Stderr, buf)
	execCmd.Stdout = outWtr
	execCmd.Stderr = errWtr

	fmt.Fprintf(outWtr, "Job %s running.\n", j.Name)
	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("%s execute error: %w\n", j.Name, err)
	}
	if err := execCmd.Wait(); err != nil {
		fmt.Fprintf(errWtr, "%s cmd wait error: %s\n", j.Name, err)
	}

	if j.Slack != "" {
		sc := slack.NewSlackClient(j.Slack)
		if err := sc.Post(buf.String()); err != nil {
			return fmt.Errorf("slack post error: %w", err)
		}
	}
	return nil
}

func createCommand(command string) (*exec.Cmd, error) {
	cmd, err := shellwords.Parse(command)
	log.Printf("Command : %q\n", cmd)
	if err != nil {
		return nil, err
	}
	switch len(command) {
	case 0:
		return nil, fmt.Errorf("blank command")
	case 1:
		return exec.Command(cmd[0]), nil
	default:
		return exec.Command(cmd[0], cmd[1:]...), nil
	}
}
