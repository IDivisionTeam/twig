package git

import (
	"errors"
	"os/exec"
)

type Cmd int

const (
	Branch = iota
	Checkout
	Fetch
	Push
	Status
	Version
)

const directive = "git"

// Git is a wrapper over exec.Command with command range restrictions.
type Git struct {
	name    string
	command Cmd
	arg     []string
	Err     error
}

func Command(cmd Cmd, arg ...string) *Git {
	git := &Git{
		name: directive,
	}

	command, err := fromCommand(cmd)
	if err != nil {
		git.Err = err
	}

	git.arg = append([]string{command}, arg...)

	return git
}

func (g *Git) Run() error {
	if g.Err != nil {
		return g.Err
	}

	return exec.Command(g.name, g.arg...).Run()
}

func (g *Git) CombinedOutput() ([]byte, error) {
	if g.Err != nil {
		return []byte{}, g.Err
	}

	return exec.Command(g.name, g.arg...).CombinedOutput()
}

func fromCommand(cmd Cmd) (string, error) {
	switch cmd {
	case Branch:
		return "branch", nil
	case Checkout:
		return "checkout", nil
	case Fetch:
		return "fetch", nil
	case Push:
		return "push", nil
	case Status:
		return "status", nil
	case Version:
		return "version", nil
	default:
		return "", errors.New("git: undefined command")
	}
}
