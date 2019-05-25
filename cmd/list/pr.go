package list

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest/search"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "List pull requests.",
	Args:  prArgs,
	RunE:  prRunE,
}

func init() {
	Cmd.AddCommand(prCmd)

	prCmd.Flags().BoolVar(&prState.Mine, "mine", false, "List my pull requests. (Set by default; use --author to override.)")
	prCmd.Flags().StringVarP(&prState.Author, "author", "a", "", "Specify the author (must not be used with explicit --mine)")
}

var prState struct {
	Author string
	Mine   bool
}

func prArgs(cmd *cobra.Command, args []string) error {
	switch {
	case prState.Author != "" && !prState.Mine:
	case prState.Author != "" && prState.Mine:
		return errors.New("Explicit mixed use of --mine and --author")
	case prState.Author == "" && !prState.Mine:
		prState.Mine = true
		fallthrough
	case prState.Author == "" && prState.Mine:
		selectedProfile, err := profile.SelectedProfile()
		if err != nil {
			return err
		}
		prState.Author = selectedProfile.Userid
	}
	return nil
}

func prRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	query := search.IssuesQuery{Author: prState.Author, Type: "pr"}
	response, err := search.Issues(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(response)
	return nil
}
