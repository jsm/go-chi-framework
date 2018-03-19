package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jsm/gode/utils/errors"
)

type linkedInUserInfo struct {
	ID string `json:"id"`
}

func getLinkedInUserInfo(token string) (linkedInUserInfo, error) {
	var userInfo linkedInUserInfo

	req, err := http.NewRequest("GET", "https://api.linkedin.com/v1/people/~?format=json", nil)
	if err != nil {
		return userInfo, errors.WithStack(err)
	}

	req.Header.Add("oauth_token", token)

	// TODO: Don't use Default Client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return userInfo, errors.WithStack(err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&userInfo); err != nil {
		return userInfo, errors.WithStack(err)
	}

	return userInfo, nil
}
