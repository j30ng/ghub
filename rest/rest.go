package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/j30ng/ghub/profile"
)

// MakeRawAPICall makes an API call to the URL passed to it, then
// returns the response body as a string.
func MakeRawAPICall(url string, token string, headers map[string]string) (*string, error) {
	println(url)
	req, err := generateRequest(url, token)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	str := string(b)
	return &str, nil
}

// MakeAPICall is like MakeAPICallWithHeaders, except it sets no headers to the request.
func MakeAPICall(profile profile.Profile, path string) (interface{}, error) {
	return MakeAPICallWithHeaders(profile, path, map[string]string{})
}

// MakeAPICallWithHeaders makes an API call via MakeRawAPICall, tries to convert the response
// into a map[string]interface{}, then returns the converted response if the conversion was successful.
func MakeAPICallWithHeaders(profile profile.Profile, path string, headers map[string]string) (interface{}, error) {
	url := strings.TrimRight(profile.APIBaseURL, "/") + "/" + strings.TrimLeft(path, "/")
	jsonString, err := MakeRawAPICall(url, profile.Token, headers)
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
