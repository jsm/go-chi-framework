package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jsm/gode/utils/errors"
)

type googleTokenInfo struct {
	UserID string `json:"sub"`
}

func getGoogleTokenInfo(token string) (googleTokenInfo, error) {
	var tokenInfo googleTokenInfo

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token))
	if err != nil {
		return tokenInfo, errors.WithStack(err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&tokenInfo); err != nil {
		return tokenInfo, errors.WithStack(err)
	}

	return tokenInfo, nil
}
