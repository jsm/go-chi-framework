package tasks

import (
	"github.com/jsm/gode/worker/tasknames"
)

// TaskMap correlates task names to their respective functions
var TaskMap = map[string]interface{}{
	tasknames.Test: test,
}

// Test worker availability
func test() (string, error) {
	return "tested", nil
}
