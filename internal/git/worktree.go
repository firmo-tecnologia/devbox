package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindRoot(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func EnsureWorktree(gitRoot, name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	worktreeDir := filepath.Join(home, ".devbox", "worktrees", filepath.Base(gitRoot), name)

	if _, err := os.Stat(worktreeDir); err == nil {
		return worktreeDir, nil
	}

	exists, err := branchExists(gitRoot, name)
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if exists {
		cmd = exec.Command("git", "-C", gitRoot, "worktree", "add", worktreeDir, name)
	} else {
		cmd = exec.Command("git", "-C", gitRoot, "worktree", "add", "-b", name, worktreeDir)
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("creating worktree: %w\n%s", err, out)
	}
	return worktreeDir, nil
}

func branchExists(gitRoot, name string) (bool, error) {
	err := exec.Command("git", "-C", gitRoot, "rev-parse", "--verify", name).Run()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}
	return false, fmt.Errorf("checking branch %q: %w", name, err)
}
