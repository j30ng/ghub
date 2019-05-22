package config

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currentProfileCmd = &cobra.Command{
	Use:   "current-profile",
	Short: "Prints out current profile.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("create-profile action takes no arguments.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if selectedProfile := viper.Get("selectedProfile"); selectedProfile == nil {
			fmt.Println("No profile is selected yet.")
		} else {
			fmt.Println(selectedProfile)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(currentProfileCmd)
}
