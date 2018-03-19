package application

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HandleRequestBodyCheck(w http.ResponseWriter, r *http.Request) (success bool) {
	// Check for Body Presence
	if r.ContentLength < 1 {
		respJ := CreateResponseJSON(false, []error{ErrBodyExpected}, nil)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJ)
		return false
	}

	return true
}

func GetRequestBodyBytes(r *http.Request) []byte {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}

func HandleJSONBytesDecode(request interface{}, w http.ResponseWriter, bodyBytes []byte, errContext map[string]string) (success bool) {
	if errContext == nil {
		errContext = map[string]string{}
	}

	errContext["request_body"] = string(bodyBytes)

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	if decodeErr := decoder.Decode(request); decodeErr != nil {
		HandleJSONDecodeError(decodeErr, errContext, w)
		return false
	}

	return true
}

// HandleJSONDecode default handling for JSON Decodes
func HandleJSONDecode(request interface{}, w http.ResponseWriter, r *http.Request, errContext map[string]string) (success bool) {
	if !HandleRequestBodyCheck(w, r) {
		return false
	}

	bodyBytes := GetRequestBodyBytes(r)

	if !HandleJSONBytesDecode(request, w, bodyBytes, errContext) {
		return false
	}

	return true
}
