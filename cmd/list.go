package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PrintMine bool
var Author string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List some useful stuffs.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New(fmt.Sprintf("Expecting only 1 argument. Got %d. %s", len(args), args))
		}
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}
		if Author != "" && PrintMine {
			return errors.New("Don't mix --mine and --author.")
		}
		if Author == "" {
			PrintMine = true
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token")
		fmt.Println(token)
		fmt.Println("list called")
		fmt.Println("args are: ", args)

	},
	ValidArgs: []string{
		"pr", "prs",
		"issue", "issues",
		"comment", "comments",
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&PrintMine, "mine", "m", false, "List only those which belong to me!")
	listCmd.Flags().StringVarP(&Author, "author", "a", "", "Specify the author (must not be used with --mine)")
}
