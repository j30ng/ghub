package search

import (
	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"
	"github.com/mitchellh/mapstructure"
)

// Commits makes an API call to the path /search/commits.
func Commits(profile profile.Profile, query CommitsQuery) (*CommitsResponse, error) {
	queryString, err := reqParamString(query)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{"Accept": "application/vnd.github.cloak-preview"}
	response, err := rest.MakeAPICallWithHeaders(profile, "/search/commits"+queryString, headers)
	if err != nil {
		return nil, err
	}
	var ret CommitsResponse
	if err = mapstructure.Decode(response, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type CommitsQuery struct {
	Q     CommitsQueryQ
	Sort  string
	Order string
}

// CommitsQuery represents the parameters of the query string to the path /search/commits.
type CommitsQueryQ struct {
	Author    []string
	Committer []string
	Org       []string
	Repo      []string
}

// CommitsResponse represents the response from the path /search/commits.
type CommitsResponse struct {
	Total_count        int
	Incomplete_results bool
	Items              []CommitsResponseItem
}

type CommitsResponseItem struct {
	Url          string
	Sha          string
	Node_id      string
	Html_url     string
	Comments_url string
	Commit       struct {
		Url    string
		Author struct {
			Date  string
			Name  string
			Email string
		}
		Committer struct {
			Date  string
			Name  string
			Email string
		}
		Message string
		Tree    struct {
			Url string
			Sha string
		}
		Comment_count int
	}
	Author struct {
		Login               string
		Id                  int
		Node_id             string
		Avatar_url          string
		Gravatar_id         string
		Url                 string
		Html_url            string
		Followers_url       string
		Following_url       string
		Gists_url           string
		Starred_url         string
		Subscriptions_url   string
		Organizations_url   string
		Repos_url           string
		Events_url          string
		Received_events_url string
		Type                string
		Site_admin          bool
	}
	Committer struct {
		Login               string
		Id                  int
		Node_id             string
		Avatar_url          string
		Gravatar_id         string
		Url                 string
		Html_url            string
		Followers_url       string
		Following_url       string
		Gists_url           string
		Starred_url         string
		Subscriptions_url   string
		Organizations_url   string
		Repos_url           string
		Events_url          string
		Received_events_url string
		Type                string
		Site_admin          bool
	}
	Parents []struct {
		Url      string
		Html_url string
		Sha      string
	}
	Repository struct {
		Id        int
		Node_id   string
		Name      string
		Full_name string
		Private   bool
		Owner     struct {
			Login               string
			Id                  int
			Node_id             string
			Avatar_url          string
			Gravatar_id         string
			Url                 string
			Html_url            string
			Followers_url       string
			Following_url       string
			Gists_url           string
			Starred_url         string
			Subscriptions_url   string
			Organizations_url   string
			Repos_url           string
			Events_url          string
			Received_events_url string
			Type                string
			Site_admin          bool
		}
		Html_url          string
		Description       interface{}
		Fork              bool
		Url               string
		Forks_url         string
		Keys_url          string
		Collaborators_url string
		Teams_url         string
		Hooks_url         string
		Issue_events_url  string
		Events_url        string
		Assignees_url     string
		Branches_url      string
		Tags_url          string
		Blobs_url         string
		Git_tags_url      string
		Git_refs_url      string
		Trees_url         string
		Statuses_url      string
		Languages_url     string
		Stargazers_url    string
		Contributors_url  string
		Subscribers_url   string
		Subscription_url  string
		Commits_url       string
		Git_commits_url   string
		Comments_url      string
		Issue_comment_url string
		Contents_url      string
		Compare_url       string
		Merges_url        string
		Archive_url       string
		Downloads_url     string
		Issues_url        string
		Pulls_url         string
		Milestones_url    string
		Notifications_url string
		Labels_url        string
		Releases_url      string
		Deployments_url   string
	}
	Score int
}

