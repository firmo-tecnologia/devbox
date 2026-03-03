package docker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/firmotecnologia/devbox/internal/config"
)

const containerHome = "/home/claude"

func Run(cfg *config.Config) error {
	if err := cfg.EnsureDirs(); err != nil {
		return err
	}

	if !cfg.NoPull {
		if err := pullImage(cfg.Image); err != nil {
			return err
		}
	}

	return runContainer(cfg)
}

func pullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pulling image %s: %w", image, err)
	}
	return nil
}

func runContainer(cfg *config.Config) error {
	args := []string{
		"run", "--rm", "-it",
		"--workdir", "/workspace",
		"-v", cfg.WorkspaceDir + ":/workspace",
		"-v", cfg.ClaudeDir + ":" + containerHome + "/.claude",
		"-v", cfg.ClaudeJSON + ":" + containerHome + "/.claude.json",
	}

	if cfg.HasDotbinsConf() {
		args = append(args,
			"-v", cfg.DotbinsConf+":"+containerHome+"/.config/dotbins/dotbins.yaml:ro",
			"-v", cfg.DotbinsCache+":"+containerHome+"/.dotbins",
		)
	}

	args = append(args, cfg.Image)

	if cfg.Shell {
		args = append(args, "/bin/bash")
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("running container: %w", err)
	}

	return nil
}
