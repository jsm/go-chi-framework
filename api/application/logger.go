package application

import (
	"os"

	"github.com/op/go-logging"
)

func getLogger() *logging.Logger {
	log := logging.MustGetLogger("api")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	var format logging.Formatter
	if Env.IsLive {
		format = logging.MustStringFormatter(
			`%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
		)
	} else {
		format = logging.MustStringFormatter(
			`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
	}
	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)
	log.SetBackend(leveledBackend)
	return log
}
