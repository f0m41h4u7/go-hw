package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/cucumber/godog"
)

//const delay = 20 * time.Second

func reconnect() error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = 3 * time.Minute
	be.InitialInterval = 3 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 20 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}
		select {
		case <-time.After(d):
			_, err := http.Get("http://calendar:1337/")
			if err != nil {
				log.Printf("could not reconnect")
				continue
			}
			return nil
		}
	}
}

func TestMain(m *testing.M) {
	//	log.Printf("wait %s for service availability...", delay)
	//	time.Sleep(delay)
	err := reconnect()
	if err != nil {
		log.Fatal(err)
	}

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	test := new(apiTest)

	s.Step(`^I send "(GET|POST)" request to "([^"]*)" with body:$`, test.iSendRequestToWithBody)
	s.Step(`^I send "(GET|POST)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^I receive event$`, test.iReceiveEvent)
	s.Step(`^I receive UUID$`, test.iReceiveUUID)

	notifyTest := new(notifyTest)
	s.BeforeScenario(notifyTest.startConsuming)
	s.Step(`^An event is about to start$`, notifyTest.anEventIsAboutToStart)
	s.Step(`^I receive notification$`, notifyTest.iReceiveNotification)
	s.AfterScenario(notifyTest.stopConsuming)
}
