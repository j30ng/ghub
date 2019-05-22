package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type pathFunction struct {
	Path (func(path string) tokenFunction)
}

type tokenFunction struct {
	Token (func(token string) (map[string]interface{}, error))
}

func Endpoint(endpoint string) pathFunction {
	client := &http.Client{}

	return generatePathFunction(client, strings.TrimRight(endpoint, "/"))
}

func generatePathFunction(client *http.Client, url string) pathFunction {
	return pathFunction{
		Path: func(path string) tokenFunction {
			url += "/" + strings.TrimLeft(path, "/")
			return generateTokenFunction(client, url)
		},
	}
}


func generateTokenFunction(client *http.Client, url string) tokenFunction {
	return tokenFunction{
		Token: func(token string) (map[string]interface{}, error) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Add("Authorization", "Bearer " + token)

			res, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			var m map[string]interface{}
			if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
				fmt.Println(m)
				fmt.Println(url)
				fmt.Println(res.StatusCode)
				return nil, err
			}
			return m, nil
		},
	}
}

