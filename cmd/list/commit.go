package list

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
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

	commitCmd.Flags().BoolVar(&commitFlags.Mine, "mine", false, "List commits I committed. (Set by default; use --committer to override.)")
	commitCmd.Flags().StringVarP(&commitFlags.Committer, "committer", "c", "", "Specify the committer (must not be used with explicit --mine)")
	commitCmd.Flags().StringVar(&commitFlags.Order, "order", "desc", "Whether to print the output in descending or ascending order. (desc, asc)")
	commitCmd.Flags().IntVar(&commitFlags.OutputSizeLimit, "limit", 0, "The max number of output elements to print. 0 indicates no limit. (>= 0)")
	commitCmd.Flags().StringVarP(&commitFlags.SortBy, "sort-by", "s", "author-date",
		"Sort the result of the API request by this measure.\n"+"https://developer.github.com/v3/search/#parameters-1")
}

var commitFlags struct {
	Committer       string
	Mine            bool
	SortBy          string
	Order           string
	OutputSizeLimit int
}

func commitArgs(cmd *cobra.Command, args []string) error {
	switch {
	case commitFlags.Committer != "" && !commitFlags.Mine:
	case commitFlags.Committer != "" && commitFlags.Mine:
		return errors.New("Explicit mixed use of --mine and --author")
	case commitFlags.Committer == "" && !commitFlags.Mine:
		commitFlags.Mine = true
		fallthrough
	case commitFlags.Committer == "" && commitFlags.Mine:
		selectedProfile, err := profile.SelectedProfile()
		if err != nil {
			return err
		}
		commitFlags.Committer = selectedProfile.Userid
	}
	switch commitFlags.SortBy {
	case "author-date", "committer-date":
	default:
		return fmt.Errorf("Invalid value for --sort-by: %s", commitFlags.SortBy)
	}
	commitFlags.Order = strings.ToLower(commitFlags.Order)
	if commitFlags.Order != "desc" && commitFlags.Order != "asc" {
		return fmt.Errorf("Invalid value for --order: %s", commitFlags.Order)
	}
	if commitFlags.OutputSizeLimit < 0 {
		return fmt.Errorf("Invalid value for --limit: %d", commitFlags.OutputSizeLimit)
	}
	return nil
}

func commitRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	query := search.CommitsQuery{Q: search.CommitsQueryQ{Committer: commitFlags.Committer},
		Sort: commitFlags.SortBy, Order: commitFlags.Order}
	response, err := search.Commits(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(commitFormatResponse(response, commitFlags.OutputSizeLimit))
	return nil
}

func commitFormatResponse(res *search.CommitsResponse, size int) string {
	result := []string{}
	for i, item := range res.Items {
		if size > 0 && i >= size {
			break
		}

		itemString := strings.Join([]string{
			"------------------------------------------------",
			color.GreenString("By: %s", item.Commit.Committer.Name) + " / " +
				color.CyanString("At: %s", reformatDate(item.Commit.Committer.Date)),
			color.YellowString("On: %s", item.Repository.Full_name),
			"",
			wrapString(item.Commit.Message, 120, 2),
		}, "\n")

		result = append([]string{itemString}, result...)
	}
	return strings.Join(result, "\n")
}

func reformatDate(dateString string) string {
	t, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		return err.Error()
	}
	return t.Format("2006-01-02 15:04")
}
