package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

// prepareXauth creates a temporary Xauthority file with wildcard family entries,
// allowing the container (with a different hostname) to authenticate with the X server.
func prepareXauth(display string) string {
	const tmpFile = "/tmp/.devbox.xauth"

	nlist, err := exec.Command("xauth", "nlist", display).Output()
	if err != nil || len(strings.TrimSpace(string(nlist))) == 0 {
		if xa := os.Getenv("XAUTHORITY"); xa != "" {
			return xa
		}
		return os.Getenv("HOME") + "/.Xauthority"
	}

	lines := strings.Split(string(nlist), "\n")
	var wildcard strings.Builder
	for _, line := range lines {
		if len(line) >= 4 {
			wildcard.WriteString("ffff" + line[4:] + "\n")
		}
	}

	os.Remove(tmpFile)
	merge := exec.Command("xauth", "-f", tmpFile, "nmerge", "-")
	merge.Stdin = strings.NewReader(wildcard.String())
	if err := merge.Run(); err != nil {
		if xa := os.Getenv("XAUTHORITY"); xa != "" {
			return xa
		}
		return os.Getenv("HOME") + "/.Xauthority"
	}

	return tmpFile
}

func runContainer(cfg *config.Config) error {
	args := []string{
		"run", "--rm", "-it",
		"--workdir", "/workspace",
	}

	if cfg.GitRoot != "" {
		args = append(args, "-v", cfg.GitRoot+":"+cfg.GitRoot)
	}

	args = append(args,
		"-v", cfg.WorkspaceDir+":/workspace",
		"-v", cfg.ClaudeDir+":"+containerHome+"/.claude",
		"-v", cfg.ClaudeJSON+":"+containerHome+"/.claude.json",
	)

	if cfg.HasDotbinsConf() {
		args = append(args,
			"-v", cfg.DotbinsConf+":"+containerHome+"/.config/dotbins/dotbins.yaml:ro",
			"-v", cfg.DotbinsCache+":"+containerHome+"/.dotbins",
		)
	}

	if display := os.Getenv("DISPLAY"); display != "" {
		xauthFile := prepareXauth(display)
		args = append(args,
			"-v", "/tmp/.X11-unix:/tmp/.X11-unix",
			"-v", xauthFile+":"+containerHome+"/.Xauthority:ro",
			"-e", "DISPLAY="+display,
			"-e", "XAUTHORITY="+containerHome+"/.Xauthority",
		)
	} else if waylandDisplay := os.Getenv("WAYLAND_DISPLAY"); waylandDisplay != "" {
		xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")
		args = append(args,
			"-v", xdgRuntime+"/"+waylandDisplay+":/tmp/"+waylandDisplay,
			"-e", "WAYLAND_DISPLAY=/tmp/"+waylandDisplay,
		)
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		args = append(args, "-e", "GITHUB_TOKEN="+token)
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
