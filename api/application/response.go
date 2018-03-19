package application

import (
	"encoding/json"
)

// Response struct for API responses
type Response struct {
	ok     bool
	errors []error
	data   interface{}
}

type responseJSON struct {
	Ok     bool        `json:"ok"`
	Errors []string    `json:"errors"`
	Data   interface{} `json:"data"`
}

// DefaultOKResponse for returning 200s
var DefaultOKResponse = CreateResponseJSON(true, nil, "Meow")

// CreateResponse for API response
func CreateResponse(ok bool, errors []error, data interface{}) Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	return Response{
		ok:     ok,
		errors: errors,
		data:   data,
	}
}

// ResponseToJSON converts a response into a JSON byteslice
func ResponseToJSON(resp Response) []byte {
	errStrings := []string{}

	for _, err := range resp.errors {
		errStrings = append(errStrings, err.Error())
	}

	respJ := responseJSON{
		Ok:     resp.ok,
		Errors: errStrings,
		Data:   resp.data,
	}

	respBytes, err := json.Marshal(respJ)
	if err != nil {
		panic(err)
	}
	return respBytes
}

// CreateResponseJSON creates a JSON byteslice directly from given parameters
func CreateResponseJSON(ok bool, errors []error, data interface{}) []byte {
	resp := CreateResponse(ok, errors, data)
	return ResponseToJSON(resp)
}
