package cmd

import (
	"fmt"
	"os"

	"github.com/firmotecnologia/devbox/internal/config"
	"github.com/firmotecnologia/devbox/internal/docker"
	"github.com/firmotecnologia/devbox/internal/git"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devbox",
	Short: "Run Claude Code inside an isolated Docker container",
	RunE: func(cmd *cobra.Command, args []string) error {
		image, _ := cmd.Flags().GetString("image")
		noPull, _ := cmd.Flags().GetBool("no-pull")
		shell, _ := cmd.Flags().GetBool("shell")
		dotbinsConf, _ := cmd.Flags().GetString("dotbins-config")
		worktreeName, _ := cmd.Flags().GetString("worktree")

		workspaceDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}

		var gitRoot string
		if worktreeName != "" {
			gitRoot, err = git.FindRoot(workspaceDir)
			if err != nil {
				return err
			}
			workspaceDir, err = git.EnsureWorktree(gitRoot, worktreeName)
			if err != nil {
				return err
			}
			fmt.Printf("worktree: %s (branch: %s)\n", workspaceDir, worktreeName)
		}

		cfg, err := config.New(image, noPull, shell, workspaceDir, gitRoot, dotbinsConf)
		if err != nil {
			return err
		}

		return docker.Run(cfg)
	},
}

func init() {
	rootCmd.Flags().StringP("image", "i", "firmotecnologia/devbox:latest", "Docker image to use")
	rootCmd.Flags().Bool("no-pull", false, "Skip docker pull before running")
	rootCmd.Flags().Bool("shell", false, "Start a bash shell instead of Claude Code")
	rootCmd.Flags().StringP("dotbins-config", "d", "", "Path to dotbins dotbins.yaml (default: ~/.dotbins/dotbins.yaml)")
	rootCmd.Flags().StringP("worktree", "w", "", "Create or reuse a git worktree with this name and run Claude Code inside it")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
