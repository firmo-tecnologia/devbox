package docker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/firmotecnologia/devbox/internal/config"
)

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
		"-v", cfg.ClaudeDir + ":/root/.claude",
	}

	if cfg.HasDotbinsConf() {
		args = append(args,
			"-v", cfg.DotbinsConf+":/root/.config/dotbins/config.yaml:ro",
			"-v", cfg.DotbinsCache+":/root/.dotbins",
		)
	}

	args = append(args, cfg.Image)

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
