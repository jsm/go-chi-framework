package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gopkg.in/gographics/imagick.v2/imagick"

	"github.com/jsm/gode/worker/application"
	"github.com/jsm/gode/worker/tasks"
)

var version string

func main() {
	instance := application.Initialize(version)
	defer instance.Close()

	m := instance.Machinery

	imagick.Initialize()
	defer imagick.Terminate()

	m.RegisterTasks(tasks.TaskMap)

	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		application.Instance.Log.Fatal(err)
	}

	// HTTP Server for health checks
	go http.ListenAndServe(":8088", router())

	// Task for monitoring amount of pending tasks
	go monitorPendingTasks(instance)

	// Setup Worker
	worker := m.NewWorker("worker", numWorkers)
	worker.ErrorHandler(func(err error) {
		instance.CaptureError(err, nil)
	})

	// Startup Metrics
	instance.StatsD.Increment("startup.worker")

	// Launch the worker
	if err := worker.Launch(); err != nil {
		application.Instance.Log.Info(err.Error())
	}
}

func router() http.Handler {
	// Initialize Main Router
	r := chi.NewRouter()

	// Setup Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)

	// Define routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Meow"))
	})

	// Return router
	return r
}

func monitorPendingTasks(instance *application.Application) {
	for {
		pendingTasks, err := instance.Machinery.GetBroker().GetPendingTasks("machinery_tasks")
		if err != nil {
			application.CaptureError(err, nil)
			continue
		}

		instance.StatsD.Gauge("worker.pending_task_count", len(pendingTasks))

		time.Sleep(time.Duration(1) * time.Second)
	}
}
