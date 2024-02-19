/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var host = "0.0.0.0:50051"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clhi",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Display a prompt to select a command
		commands := []string{"login", "join-room", "other-command"}
		for _, c := range commands {
			fmt.Printf(c)
		}
	},
	// PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	conn := viper.Get("Connection").(*grpc.ClientConn)
	// 	defer conn.Close()
	// },
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clhi.yaml)")
	rootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0:50051", "chat server host (default is 0.0.0.0:50051)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
