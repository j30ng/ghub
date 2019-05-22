package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigFile string

var RootCmd = &cobra.Command{
	Use:   "ghub",
	Short: "A simple wrapper for GitHub REST API calls.",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defaultConfigFile := getHomeDirectory() + "/.config/ghub/ghub.yaml"
	description := "config file (default is $HOME/.config/ghub/ghub.yaml)"

	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", defaultConfigFile, description)
}

func initConfig() {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(ConfigFile), os.ModeDir|os.ModePerm)
		os.Create(ConfigFile)
	}

	viper.SetConfigFile(ConfigFile)

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
