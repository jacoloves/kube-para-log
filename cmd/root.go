/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jacoloves/kube-para-log/internal/kubectl"
	"github.com/jacoloves/kube-para-log/internal/tmux"
	"github.com/spf13/cobra"
)

var namespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-para-log [keyword]",
	Short: "Kubernetes pod los in parallel tmux panes",
	Args:  cobra.ExactArgs(1), // Specify only one argument.
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		fmt.Printf("👓 Searching for pods containing keyword: '%s' in namespace: '%s'\n", keyword, namespace)

		// internal/kubectl package function
		pods, err := kubectl.FindMatchingPods(keyword, namespace)
		if err != nil {
			fmt.Println("❎ Error:", err)
			os.Exit(1)
		}

		if len(pods) == 0 {
			fmt.Println("🚧 No mathing pods found.")
			return
		}

		fmt.Println("✅ Stating tmux session with logs...")
		err = tmux.StartTmuxWithLogs("kube-para-log", pods, namespace)
		if err != nil {
			fmt.Println("❎ tmux error:", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-para-log.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes maespace to search pods in")
}
