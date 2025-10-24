/*
Copyright © 2025 Mburu edupablo72@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version = "v0.0.1-beta"

// rootCmd is the main command for the CLI tool
var rootCmd = &cobra.Command{
	Use:   "domsnif",
	Short: "Domsnif helps web developers discover new projects and potential gigs.",
	Long: `Domsnif is a CLI tool that helps web developers find websites that are coming soon or under development, 
	extract associated emails, verify them, and forward the contacts to the developer. 
	This tool was originally created for personal use and is now shared with the public. 
	Feel free to reach out for any questions or support.`,
	Version: version,
}

// Execute runs the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.domsnif.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig loads config file and environment variables
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".domsnif")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
