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
	RunE:  runCommitCmdWithError,
}

func init() {
	Cmd.AddCommand(commitCmd)

	commitCmd.Flags().Bool("mine", false, "List commits I committed.")
	commitCmd.Flags().StringP("committers", "c", "", "Array of committers. Comma-separated.")
	commitCmd.Flags().StringP("repos", "r", "", "Array of repos. Comma-separated. (e.g. --repos=/org/project,/user/project)")
	commitCmd.Flags().String("order", "desc", "Print the output in descending or ascending order. (desc, asc)")
	commitCmd.Flags().Int("limit", 0, "The max number of output elements to print. 0 indicates no limit. (>= 0)")
	commitCmd.Flags().StringP("sort-by", "s", "author-date",
		"Sort the result of the API request by this measure.\n"+"https://developer.github.com/v3/search/#parameters-1")
}

func runCommitCmdWithError(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	param, err := parseCommitFlag(cmd)
	if err != nil {
		return err
	}
	query := search.CommitsQuery{Q: search.CommitsQueryQ{Committer: param.Committers, Repo: param.Repos},
		Sort: param.Sort, Order: param.Order}
	response, err := search.Commits(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(commitFormatResponse(response, param.OutputCountLimit))
	return nil
}

type commitCmdParam struct {
	Committers       []string
	Repos            []string
	Sort             string
	Order            string
	OutputCountLimit int
}

func parseCommitFlag(cmd *cobra.Command) (param *commitCmdParam, err error) {
	param, err = nil, nil

	param = &commitCmdParam{}
	err, includeMine, csvCommitters, csvRepos, order, limit, sortBy := getCommitFlagVals(cmd)
	if err != nil {
		return nil, err
	}

	param.Committers, err = determineCommitters(includeMine, csvCommitters)
	if err != nil {
		return nil, err
	}

	for _, repo := range strings.Split(csvRepos, ",") {
		param.Repos = append(param.Repos, strings.TrimSpace(repo))
	}

	if order != "desc" && order != "asc" {
		return nil, errors.New("--order must be either one of (asc, desc)")
	}
	param.Order = order

	if limit < 0 {
		return nil, fmt.Errorf("--limit must be >= 0; got %d", limit)
	}
	param.OutputCountLimit = limit

	switch sortBy {
	case "author-date", "committer-date":
		param.Sort = sortBy
	default:
		return nil, fmt.Errorf("Invalid value for --sort-by: %s", sortBy)
	}

	return
}

func getCommitFlagVals(cmd *cobra.Command) (err error, mine bool, committers string, repos string, order string, limit int, sortBy string) {
	flags := cmd.Flags()
	mine      , err = flags.GetBool("mine");         if err != nil { return }
	committers, err = flags.GetString("committers"); if err != nil { return }
	repos     , err = flags.GetString("repos");      if err != nil { return }
	order     , err = flags.GetString("order");      if err != nil { return }
	limit     , err = flags.GetInt("limit");         if err != nil { return }
	sortBy    , err = flags.GetString("sort-by");    if err != nil { return }
	return
}

func determineCommitters(includeMine bool, csvCommitters string) (committers []string, err error) {
	for _, committer := range strings.Split(csvCommitters, ",") {
		trimmedCommitter := strings.TrimSpace(committer)
		if trimmedCommitter != "" {
			committers = append(committers, trimmedCommitter)
		}
	}

	if includeMine || len(committers) == 0 {
		profile, err := profile.SelectedProfile();
		if err != nil {
			return nil, err
		}
		if !contains(committers, profile.Userid) {
			committers = append(committers, profile.Userid)
		}
	}
	return
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}


func commitFormatResponse(res *search.CommitsResponse, size int) string {
	result := []string{}
	for i, item := range res.Items {
		if size > 0 && i >= size {
			break
		}
		result = append([]string{formatCommitResponseItem(&item)}, result...)
		//result = append(result, formatCommitResponseItem(&item))
	}
	line := "------------------------------------------------"
	return line + "\n" + strings.Join(result, "\n" + line + "\n") + "\n" + line
}

func formatCommitResponseItem(item *search.CommitsResponseItem) string {
	return strings.Join([]string{
		color.YellowString("Commit %s", item.Sha),
		color.RedString("By: %s", item.Commit.Committer.Name),
		color.CyanString("At: %s", reformatDate(item.Commit.Committer.Date)),
		color.GreenString("On: %s", item.Repository.Full_name),
		"",
		foldString(item.Commit.Message, 120, 2),
	}, "\n")

}

func reformatDate(dateString string) string {
	t, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		return err.Error()
	}
	return t.Local().Format("2006-01-02 15:04")
}
