package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest/search"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "List issues.",
	Args:  issueArgs,
	RunE:  issueRunE,
}

func init() {
	Cmd.AddCommand(issueCmd)

	issueCmd.Flags().BoolVar(&issueFlags.Mine, "mine", false, "List issues I filed. (Set by default; use --author to override.)")
	issueCmd.Flags().StringVarP(&issueFlags.Author, "author", "a", "", "Specify the author (must not be used with explicit --mine)")
	issueCmd.Flags().StringVar(&issueFlags.Order, "order", "desc", "Whether to print the output in descending or ascending order. (desc, asc)")
	issueCmd.Flags().IntVarP(&issueFlags.Cols, "cols", "c", 120, "Fold lines longer than this value. 0 indicates no folding. (>= 0)")
	issueCmd.Flags().IntVar(&issueFlags.OutputSizeLimit, "limit", 0, "The max number of output elements to print. 0 indicates no limit. (>= 0)")
	issueCmd.Flags().StringVarP(&issueFlags.SortBy, "sort-by", "s", "updated",
		"Sort the result of the API request by this measure.\n"+"https://developer.github.com/v3/search/#parameters-3")
}

var issueFlags struct {
	Author          string
	Mine            bool
	Order           string
	SortBy          string
	Cols            int
	OutputSizeLimit int
}

func issueArgs(cmd *cobra.Command, args []string) error {
	switch {
	case issueFlags.Author != "" && !issueFlags.Mine:
	case issueFlags.Author != "" && issueFlags.Mine:
		return errors.New("Mixed use of --mine and --author")
	case issueFlags.Author == "" && !issueFlags.Mine:
		issueFlags.Mine = true
		fallthrough
	case issueFlags.Author == "" && issueFlags.Mine:
		selectedProfile, err := profile.SelectedProfile()
		if err != nil {
			return err
		}
		issueFlags.Author = selectedProfile.Userid
	}
	switch issueFlags.SortBy {
	case "comments", "reactions", "reactions-+1", "reactions--1", "reactions-smile", "reactions-thinking_face", "reactions-heart", "reactions-tada", "interactions", "created", "updated":
	default:
		return fmt.Errorf("Invalid value for --sort-by: %s", issueFlags.SortBy)
	}
	issueFlags.Order = strings.ToLower(issueFlags.Order)
	if issueFlags.Order != "desc" && issueFlags.Order != "asc" {
		return fmt.Errorf("Invalid value for --order: %s", issueFlags.Order)
	}
	if issueFlags.Cols < 0 {
		return fmt.Errorf("Invalid value for --cols: %d", issueFlags.Cols)
	}
	return nil
}

func issueRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	query := search.IssuesQuery{Q: search.IssuesQueryQ{Author: issueFlags.Author, Type: "issue"}, Order: issueFlags.Order, Sort: issueFlags.SortBy}
	response, err := search.Issues(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(issueFormatResponse(response, issueFlags.Order, issueFlags.Cols, issueFlags.OutputSizeLimit))
	return nil
}

func issueFormatResponse(res *search.IssuesResponse, order string, cols int, size int) string {
	result := []string{}
	for i, item := range res.Items {
		if size > 0 && i >= size {
			break
		}

		itemString := strings.Join([]string{
			"------------------------------------------------",
			color.RedString("[%s] ", item.State) + fmt.Sprintf("%s", item.Title),
			color.YellowString("URL: %s", item.Html_url),
			color.CyanString("Updated At: %s", reformatDate(item.Updated_at)) + " / " + color.CyanString("Created At: %s", reformatDate(item.Created_at)),
			"",
			wrapString(item.Body, cols, 2),
		}, "\n")

		result = append([]string{itemString}, result...)
	}
	return strings.Join(result, "\n")
}

func wrapString(in string, cols int, indent int) string {
	if indent < 0 {
		indent = 0
	}

	lines := strings.Split(in, "\n")
	b := strings.Builder{}
	for _, line := range lines {
		b.WriteString(strings.Repeat(" ", indent) + strings.Join(partitionString(line, cols), "\n") + "\n")
	}
	return b.String()
}

func partitionString(in string, size int) []string {
	runeIn := []rune(in)
	result := []string{}
	var line_size int
	for ; len(runeIn) > 0; runeIn = runeIn[line_size:] {
		line_size = size
		if line_size > len(runeIn) {
			line_size = len(runeIn)
		}
		result = append(result, string(runeIn[:line_size]))
	}
	return result
}
