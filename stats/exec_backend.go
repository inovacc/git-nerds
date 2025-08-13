package stats

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecBackend implementa Backend usando el binario `git` instalado.
type ExecBackend struct {
	RepoPath string
}

func NewExecBackend(repoPath string) *ExecBackend {
	return &ExecBackend{RepoPath: repoPath}
}

func (b *ExecBackend) runGit(args ...string) (string, error) {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command("git", args...)
	cmd.Dir = b.RepoPath
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %v - %s", strings.Join(args, " "), err, stderr.String())
	}

	return out.String(), nil
}

func (b *ExecBackend) Log(args ...string) (string, error) {
	return b.runGit(append([]string{"log"}, args...)...)
}

func (b *ExecBackend) Branches(args ...string) (string, error) {
	return b.runGit(append([]string{"branch"}, args...)...)
}
