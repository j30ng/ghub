package list

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest/search"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "List commits.",
	Args:  commitArgs,
	RunE:  commitRunE,
}

func init() {
	Cmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVar(&commitState.Mine, "mine", false, "List my commits. (Set by default; use --author to override.)")
	commitCmd.Flags().StringVarP(&commitState.Author, "author", "a", "", "Specify the author (must not be used with explicit --mine)")
}

var commitState struct {
	Author string
	Mine   bool
}

func commitArgs(cmd *cobra.Command, args []string) error {
	switch {
	case commitState.Author != "" && !commitState.Mine:
	case commitState.Author != "" && commitState.Mine:
		return errors.New("Explicit mixed use of --mine and --author")
	case commitState.Author == "" && !commitState.Mine:
		commitState.Mine = true
		fallthrough
	case commitState.Author == "" && commitState.Mine:
		selectedProfile, err := profile.SelectedProfile()
		if err != nil {
			return err
		}
		commitState.Author = selectedProfile.Userid
	}
	return nil
}

func commitRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	query := search.CommitsQuery{Author: commitState.Author}
	response, err := search.Commits(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(response.Items)
	return nil
}
