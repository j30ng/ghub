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
	RunE:  prRunE,
}

func init() {
	Cmd.AddCommand(prCmd)

	prCmd.Flags().Bool("mine", false, "List my pull requests. (Set by default; use --author to override.)")
	prCmd.Flags().StringP("authors", "a", "", "Array of authors. Comma-separated.")
	prCmd.Flags().Bool("all-authors", false, "Search for pull requests by all authors. Overrides --mine and --authors.")

	prCmd.Flags().StringP("repos", "r", "", "Array of repos. Comma-separated.")
	prCmd.Flags().StringP("orgs", "o", "", "Array of organizations. Comma-separated.")
	prCmd.Flags().Bool("open", false, "List open pull requests.")
	prCmd.Flags().Bool("closed", false, "List closed pull requests.")
	prCmd.Flags().String("order", "desc", "Print the output in descending or ascending order. (desc, asc)")
	prCmd.Flags().StringP("sort-by", "s", "updated",
		"Sort the result of the API request by this measure.\n"+"https://developer.github.com/v3/search/#parameters-3")

	prCmd.Flags().IntP("cols", "c", 120, "Fold lines longer than this value. 0 indicates no folding.")
	prCmd.Flags().Int("limit", 0, "The max number of output elements to print. 0 indicates no limit. (>= 0)")
}

type prCmdParam struct {
	Authors         []string
	Repos           []string
	Orgs            []string
	States          []string
	Order           string
	Sort            string
	Cols            int
	OutputCountLimit int
}

func parsePrFlag(cmd *cobra.Command) (param *prCmdParam, err error) {
	param, err = &prCmdParam{}, nil
	err, mine, csvAuthors, allAuthors, csvRepos, csvOrgs, open, closed, order, sortBy, cols, limit := getPrFlagVals(cmd)
	if err != nil {
		return nil, err
	}

	param.Authors, err = determineAuthors(allAuthors, mine, csvAuthors)
	if err != nil {
		return nil, err
	}

	if csvRepos != "" {
		for _, repo := range strings.Split(csvRepos, ",") {
			param.Repos = append(param.Repos, strings.TrimSpace(repo))
		}
	}

	if csvOrgs != "" {
		for _, org := range strings.Split(csvOrgs, ",") {
			param.Orgs = append(param.Orgs, strings.TrimSpace(org))
		}
	}

	if order != "desc" && order != "asc" {
		return nil, errors.New("--order must be either one of (asc, desc)")
	}
	param.Order = order

	switch sortBy {
	case "comments", "reactions", "reactions-+1", "reactions--1", "reactions-smile", "reactions-thinking_face", "reactions-heart", "reactions-tada", "interactions", "created", "updated":
		param.Sort = sortBy
	default:
		return nil, fmt.Errorf("Invalid value for --sort-by: %s", sortBy)
	}

	if open {
		param.States = append(param.States, "open")
	}
	if closed {
		param.States = append(param.States, "closed")
	}

	if limit < 0 {
		return nil, fmt.Errorf("--limit must be >= 0; got %d", limit)
	}
	param.OutputCountLimit = limit

	if cols < 0 {
		return nil, fmt.Errorf("--cols must be >= 0; got %d", cols)
	}
	param.Cols = cols

	return
}

func getPrFlagVals(cmd *cobra.Command) (err error, mine bool, csvAuthors string, allAuthors bool, csvRepos string, csvOrgs string, open bool, closed bool, order string, sortBy string, cols int, limit int) {
	flags := cmd.Flags()
	mine      , err = flags.GetBool("mine");        if err != nil { return }
	csvAuthors, err = flags.GetString("authors");   if err != nil { return }
	allAuthors, err = flags.GetBool("all-authors"); if err != nil { return }
	csvRepos  , err = flags.GetString("repos");     if err != nil { return }
	csvOrgs   , err = flags.GetString("orgs");      if err != nil { return }
	open      , err = flags.GetBool("open");        if err != nil { return }
	closed    , err = flags.GetBool("closed");      if err != nil { return }
	order     , err = flags.GetString("order");     if err != nil { return }
	sortBy    , err = flags.GetString("sort-by");   if err != nil { return }
	cols      , err = flags.GetInt("cols");         if err != nil { return }
	limit     , err = flags.GetInt("limit");        if err != nil { return }
	return
}

func determineAuthors(includeAll bool, includeMine bool, csvAuthors string) (authors []string, err error) {
	if includeAll {
		return
	}
	return determineCommitters(includeMine, csvAuthors)
}

func prRunE(cmd *cobra.Command, args []string) error {
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("An error occurred while loading profile.\n\n" + err.Error())
	}
	param, err := parsePrFlag(cmd)
	if err != nil {
		return err
	}
	query := search.IssuesQuery{Q: search.IssuesQueryQ{Author: param.Authors, Repo: param.Repos, Org: param.Orgs, Type: []string{"pr"}, State: param.States}, Order: param.Order, Sort: param.Sort}
	response, err := search.Issues(*profile, query)
	if err != nil {
		return errors.New("An error occurred while making an API call.\n\n" + err.Error())
	}
	fmt.Println(issueFormatResponse(response, param.Order, param.Cols, param.OutputCountLimit))
	return nil
}
