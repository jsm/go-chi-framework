package application

import (
	"net/http"

	"github.com/jsm/gode/services/errors"
)

func AuthenticateUser(r *http.Request, w http.ResponseWriter) (string, error) {
	userFirebaseKey := r.Header.Get(HeaderUserID)

	if userFirebaseKey == "" {
		jsonResponse := CreateResponseJSON(false, nil, nil)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)

		return "", serviceerrors.InvalidFirebaseKeyError(userFirebaseKey)
	}

	return userFirebaseKey, nil
}
