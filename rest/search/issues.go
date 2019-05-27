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
	queryString, err := reqParamString(query)
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

func reqParamString(query interface{}) (string, error) {
	marshalled, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	var reqParamMap map[string]interface{}
	if err = json.Unmarshal(marshalled, &reqParamMap); err != nil {
		return "", err
	}
	var reqParamItems []string
	for reqParam, paramVal := range reqParamMap {
		if formattedVal, err := formatParamVal(paramVal); err != nil {
			return "", err
		} else if formattedVal != nil {
			reqParamItems = append(reqParamItems, fmt.Sprintf("%s=%s", strings.ToLower(reqParam), *formattedVal))
		}
	}
	return "?" + strings.Join(reqParamItems, "&"), nil
}

func formatParamVal(paramVal interface{}) (*string, error) {
	result := ""
	switch v := paramVal.(type) {
	case map[string]interface{}:
		resultItems := []string{}
		for key, val := range v {
			if val != "" {
				resultItems = append(resultItems, fmt.Sprintf("%s:%s", strings.ToLower(key), val))
			}
		}
		result = strings.Join(resultItems, "+")
	case string:
		result = v
	default:
		return nil, fmt.Errorf("Unexpected paramVal type: %T", v)

	}
	if result == "" {
		return nil, nil
	}
	return &result, nil
}

// IssuesQuery represents the parameters of the query string to the path /search/issues.
type IssuesQuery struct {
	Q     IssuesQueryQ
	Sort  string
	Order string
}

type IssuesQueryQ struct {
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
