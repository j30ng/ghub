package search

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"
	"github.com/mitchellh/mapstructure"
)

// Issues makes an API call to the path /search/issues.
func Issues(profile profile.Profile, query IssuesQuery) (*IssuesResponse, error) {
	queryString, err := generateQueryString(query)
	if err != nil {
		return nil, err
	}
	response, err := rest.MakeAPICall(profile, "/search/issues"+queryString)
	if err != nil {
		return nil, err
	}
	var ret IssuesResponse
	if err = mapstructure.Decode(response, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func generateQueryString(query IssuesQuery) (string, error) {
	marshalled, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	var queryMap map[string]interface{}
	if err = json.Unmarshal(marshalled, &queryMap); err != nil {
		return "", err
	}
	var queryStringItems []string
	for k, v := range queryMap {
		queryStringItems = append(queryStringItems, fmt.Sprintf("%s:%s", k, v))
	}
	return "?q=" + strings.Join(queryStringItems, "+"), nil
}

// IssuesQuery represents the parameters of the query string to the path /search/issues.
type IssuesQuery struct {
	Author string
	State  string
	Type   string
}

// IssuesResponse represents the response from the path /search/issues.
type IssuesResponse struct {
	Total_count        int
	Incomplete_results bool
	Items              []struct {
		Url                string
		Repository_url     string
		Labels_url         string
		Comments_url       string
		Events_url         string
		Html_url           string
		Id                 int
		Node_id            string
		Number             int
		Title              string
		User               IssuesUser
		Labels             []IssuesLabel
		State              string
		Locked             bool
		Assignee           IssuesUser
		Assignees          []IssuesUser
		Milestone          IssuesMilestone
		Comments           int
		Created_at         string
		Updated_at         string
		Closed_at          string
		Author_association string
		Body               string
		Score              int
	}
}

// IssuesUser represents a sub-structure used inside IssuesResponse.
type IssuesUser struct {
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

// IssuesLabel represents a sub-structure used inside IssuesResponse.
type IssuesLabel struct {
	Id      int
	Node_id string
	Url     string
	Name    string
	Color   string
	Default bool
}

// IssuesMilestone represents a sub-structure used inside IssuesResponse.
type IssuesMilestone interface{}
