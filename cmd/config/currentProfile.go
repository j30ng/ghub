package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currentProfileCmd = &cobra.Command{
	Use:   "create-profile",
	Short: "Create a ghub profile.",
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
	configCmd.AddCommand(currentProfileCmd)
}
