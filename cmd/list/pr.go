package list

import (
	"errors"
	"fmt"
	"strings"

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

	prCmd.Flags().BoolVar(&prFlags.Mine, "mine", false, "List my pull requests. (Set by default; use --author to override.)")
	prCmd.Flags().StringVarP(&prFlags.Author, "author", "a", "", "Specify the author (must not be used with explicit --mine)")
	prCmd.Flags().StringVar(&prFlags.Order, "order", "desc", "Whether to print the output in descending or ascending order.")
	prCmd.Flags().IntVarP(&prFlags.Cols, "cols", "c", 120, "Fold lines longer than this value. 0 indicates no folding.")
	prCmd.Flags().IntVar(&prFlags.OutputSizeLimit, "limit", 0, "The max number of output elements to print. 0 indicates no limit. (>= 0)")
	prCmd.Flags().StringVarP(&prFlags.SortBy, "sort-by", "s", "updated",
		"Sort the result of the API request by this measure.\n"+"https://developer.github.com/v3/search/#parameters-3")

}

var prFlags struct {
	Author          string
	Mine            bool
	Order           string
	Cols            int
	SortBy          string
	OutputSizeLimit int
}

func prArgs(cmd *cobra.Command, args []string) error {
	switch {
	case prFlags.Author != "" && !prFlags.Mine:
	case prFlags.Author != "" && prFlags.Mine:
		return errors.New("Explicit mixed use of --mine and --author")
	case prFlags.Author == "" && !prFlags.Mine:
		prFlags.Mine = true
		fallthrough
	case prFlags.Author == "" && prFlags.Mine:
		selectedProfile, err := profile.SelectedProfile()
		if err != nil {
			return err
		}
		prFlags.Author = selectedProfile.Userid
	}
	switch prFlags.SortBy {
	case "comments", "reactions", "reactions-+1", "reactions--1", "reactions-smile", "reactions-thinking_face", "reactions-heart", "reactions-tada", "interactions", "created", "updated":
	default:
		return fmt.Errorf("Invalid value for --sort-by: %s", prFlags.SortBy)
	}
	prFlags.Order = strings.ToLower(prFlags.Order)
	if prFlags.Order != "desc" && prFlags.Order != "asc" {
		return fmt.Errorf("Invalid value for --order: %s", prFlags.Order)
	}
	if prFlags.Cols < 0 {
		return fmt.Errorf("Invalid value for --cols: %d", prFlags.Cols)
	}
	return nil
}

func prRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	query := search.IssuesQuery{Q: search.IssuesQueryQ{Author: prFlags.Author, Type: "pr"}, Order: prFlags.Order, Sort: prFlags.SortBy}
	response, err := search.Issues(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(issueFormatResponse(response, prFlags.Order, prFlags.Cols, prFlags.OutputSizeLimit))
	return nil
}
