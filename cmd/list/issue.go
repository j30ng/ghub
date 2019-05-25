package list

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest/search"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "List issues.",
	Args: issueArgs,
	RunE: issueRunE,
}

func init() {
	Cmd.AddCommand(issueCmd)

	issueCmd.Flags().BoolVar(&issueState.Mine, "mine", false, "List issues I filed. (Set by default; use --author to override.)")
	issueCmd.Flags().StringVarP(&issueState.Author, "author", "a", "", "Specify the author (must not be used with explicit --mine)")
}

var issueState struct {
	Author string
	Mine   bool
}

func issueArgs(cmd *cobra.Command, args []string) error {
	switch {
	case issueState.Author != "" && !issueState.Mine:
	case issueState.Author != "" && issueState.Mine:
		return errors.New("Don't mix --mine and --author.")
	case issueState.Author == "" && !issueState.Mine:
		issueState.Mine = true
		fallthrough
	case issueState.Author == "" && issueState.Mine:
		if selectedProfile, err := profile.GetSelectedProfile(); err != nil {
			return err
		} else {
			issueState.Author = selectedProfile.Userid
		}
	}
	return nil
}

func issueRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.GetSelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	response, err := search.Issues(*profile, search.IssuesQuery{Author: issueState.Author, Type: "issue"})
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(response)
	return nil
}

