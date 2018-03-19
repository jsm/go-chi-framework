package application

import (
	"encoding/json"
	"fmt"
)

func formatError(mainErr error, context map[string]string) string {
	if context != nil {
		contextJSON, err := json.Marshal(context)
		if err == nil {
			return fmt.Sprintf("Error: %+v\nError Context: %s", mainErr, string(contextJSON))
		}
		return fmt.Sprintf("Error: %+v\nError Context: %s", mainErr, "PARSINGFAILED")
	}
	return fmt.Sprintf("Error: %+v", mainErr)
}

// CaptureError and route to appropriate mechanisms
func CaptureError(mainErr error, context map[string]string) {
	Instance.CaptureError(mainErr, context)
}
