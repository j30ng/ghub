package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j30ng/ghub/cmd/config"
	"github.com/j30ng/ghub/cmd/list"
	"github.com/j30ng/ghub/cmd/raw"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "ghub",
	Short: "A simple wrapper for GitHub REST API calls.",
}

// Execute executes the command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defaultConfigFile := getHomeDirectory() + "/.config/ghub/ghub.yaml"
	description := "config file (default is $HOME/.config/ghub/ghub.yaml)"

	rootCmd.PersistentFlags().StringVar(&configFile, "config", defaultConfigFile, description)

	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(raw.Cmd)
	rootCmd.AddCommand(config.Cmd)
}

func initConfig() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(configFile), os.ModeDir|os.ModePerm)
		os.Create(configFile)
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getHomeDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
