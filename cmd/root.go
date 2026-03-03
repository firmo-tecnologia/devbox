package cmd

import (
	"fmt"
	"os"

	"github.com/firmotecnologia/devbox/internal/config"
	"github.com/firmotecnologia/devbox/internal/docker"
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

		workspaceDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}

		cfg, err := config.New(image, noPull, shell, workspaceDir, dotbinsConf)
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
