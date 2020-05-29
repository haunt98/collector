package httpwrap

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Param struct {
	Name  string
	Value string
}

func AddParams(originalURL string, params ...Param) (string, error) {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for _, param := range params {
		if param.Name != "" && param.Value != "" {
			q.Set(param.Name, param.Value)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func DoRequestWithResult(client *http.Client, req *http.Request, result interface{}) error {
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}
