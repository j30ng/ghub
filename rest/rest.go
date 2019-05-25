package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/j30ng/ghub/profile"
)

func MakeRawAPICallWithProfile(profile profile.Profile, path string) (*string, error) {
	return MakeRawAPICall(profile.Endpoint, profile.Token, path)
}

func MakeRawAPICall(endpoint string, token string, path string) (*string, error) {
	apiBaseUrl := strings.TrimRight(endpoint, "/")
	fullUrl := apiBaseUrl + "/" + strings.TrimLeft(path, "/")

	req, err := generateRequest(fullUrl, token)
	if err != nil {
		return nil, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if b, err := ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	} else {
		str := string(b)
		return &str, nil
	}
}

func MakeAPICall(profile profile.Profile, path string) (interface{}, error) {
	jsonString, err := MakeRawAPICallWithProfile(profile, path)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	reader := strings.NewReader(*jsonString)
	if err = json.NewDecoder(reader).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func generateRequest(url string, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return req, nil
}
