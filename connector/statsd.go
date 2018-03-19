package connector

import (
	"log"
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
	statsd "gopkg.in/alexcesaro/statsd.v2"

	"github.com/jsm/gode/utils/errors"
)

func ConnectStatsD() *statsd.Client {
	var client *statsd.Client

	action := func(attempt uint) error {
		var err error
		client, err = statsd.New(statsd.Address("stats:8125"))
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	err := retry.Retry(
		action,
		strategy.Limit(5),
		strategy.BackoffWithJitter(
			backoff.Linear(time.Duration(1)*time.Second),
			jitter.Deviation(random, 0.5)),
	)

	if err != nil {
		log.Println("Failed to connect statsd")
		log.Fatalln(err)
	}

	return client
}

func ConnectMockStatsD() *statsd.Client {
	client, _ := statsd.New()
	return client
}
