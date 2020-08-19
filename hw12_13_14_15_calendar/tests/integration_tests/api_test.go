package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/messages-go/v10"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

const address = "http://calendar:1337"

var (
	errEmptyResponse = errors.New("received empty response")

	ev1 = internal.Event{
		Title:       "event1",
		Start:       "2049-07-08T16:00:00",
		End:         "2049-07-08T16:30:00",
		Description: "event 1",
		OwnerID:     "9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a",
		NotifyIn:    "1h",
	}
	ev2 = internal.Event{
		Title:       "event3",
		Start:       "2049-07-25T11:05:00",
		End:         "2049-07-25T19:41:00",
		Description: "event 2",
		OwnerID:     "9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a",
		NotifyIn:    "1h",
	}
)

type getResponse struct {
	UUID string `json:"uuid"`
}

type apiTest struct {
	resp               *httptest.ResponseRecorder
	responseStatusCode int
	getResponseBody    []byte
	postResponseBody   []byte

	eventUUID uuid.UUID
	event     internal.Event
}

func (test *apiTest) iSendRequestTo(httpMethod, path string) error {
	switch httpMethod {
	case http.MethodGet:
		r, err := http.Get(address + path)
		if err != nil {
			return err
		}
		test.responseStatusCode = r.StatusCode

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		str, err := strconv.Unquote(string(bodyBytes))
		if err != nil {
			return err
		}
		test.getResponseBody, err = base64.StdEncoding.DecodeString(str)

	case http.MethodDelete:
		req, err := http.NewRequest(http.MethodDelete, address+path+"/"+test.eventUUID.String(), bytes.NewBuffer([]byte{}))
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		test.responseStatusCode = resp.StatusCode
	default:
		return fmt.Errorf("unknown method: %s", httpMethod)
	}

	return nil
}

func (test *apiTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *apiTest) iSendRequestToWithBody(httpMethod, path string, body *messages.PickleStepArgument_PickleDocString) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(body.Content)
		r, err = http.Post(address+path, "application/json", bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseStatusCode = r.StatusCode
	test.postResponseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *apiTest) iReceiveEvent() error {
	var ev internal.EventList
	err := ev.UnmarshalJSON(test.getResponseBody)
	if err != nil {
		return err
	}
	if len(ev) == 0 {
		return errEmptyResponse
	}
	if ev[0].UUID != test.eventUUID.String() {
		return fmt.Errorf("event uuid %s received in get request doesn't match %s", ev[0].UUID, test.eventUUID.String())
	}
	test.event = ev[0]
	return nil
}

func (test *apiTest) iReceiveUUID() error {
	var resp getResponse
	err := json.Unmarshal(test.postResponseBody, &resp)
	if err != nil {
		return err
	}
	test.eventUUID, err = uuid.Parse(resp.UUID)
	return err
}

func (test *apiTest) iSendRequestToUpdatingTitleTo(path, title string) error {
	test.event.Title = title
	evBytes, err := test.event.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, address+path+"/"+test.eventUUID.String(), bytes.NewReader(evBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	test.responseStatusCode = resp.StatusCode
	return nil
}

func (test *apiTest) iReceiveUpdatedEvent() error {
	var ev internal.EventList
	err := ev.UnmarshalJSON(test.getResponseBody)
	if err != nil {
		return err
	}
	if len(ev) == 0 {
		return errEmptyResponse
	}
	if ev[0].Title != test.event.Title {
		return fmt.Errorf("updated event %s doesn't match %s", ev[0].Title, test.event.Title)
	}
	return nil
}

func (test *apiTest) iReceiveListOfEventsFor(time string) error {
	var ev internal.EventList
	err := ev.UnmarshalJSON(test.getResponseBody)
	if err != nil {
		return err
	}
	if len(ev) == 0 {
		return errEmptyResponse
	}

	switch time {
	case "week":
		if len(ev) != 2 {
			return fmt.Errorf("wrong quantity of events: expected 2, got %d", len(ev))
		}
		if ev[0].UUID != test.event.UUID {
			return fmt.Errorf("wrong event: expected %+v, got %+v", test.event, ev[0])
		}
		if ev[1].UUID != ev1.UUID {
			return fmt.Errorf("wrong event: expected %+v, got %+v", ev1, ev[1])
		}
	case "month":
		if len(ev) != 3 {
			return fmt.Errorf("wrong quantity of events: expected 2, got %d", len(ev))
		}
		if ev[2].UUID != ev2.UUID {
			return fmt.Errorf("wrong event: expected %+v, got %+v", ev2, ev[2])
		}
	}

	return nil
}

func (test *apiTest) iReceiveNoEventsAt(path string) error {
	r, err := http.Get(address + path)
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(bodyBytes) != 0 {
		return fmt.Errorf("received not nil result: %+v", bodyBytes)
	}

	return nil
}

func (test *apiTest) thereAreEventsInDB() error {
	evBytes, err := ev1.MarshalJSON()
	if err != nil {
		return err
	}
	req, err := http.Post(address+"/create", "application/json", bytes.NewReader([]byte(evBytes)))
	if err != nil {
		return err
	}
	resp, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	var uuid getResponse
	err = json.Unmarshal(resp, &uuid)
	if err != nil {
		return err
	}
	ev1.UUID = uuid.UUID

	evBytes, err = ev2.MarshalJSON()
	if err != nil {
		return err
	}
	req, err = http.Post(address+"/create", "application/json", bytes.NewReader([]byte(evBytes)))
	if err != nil {
		return err
	}
	resp, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &uuid)
	if err != nil {
		return err
	}
	ev2.UUID = uuid.UUID

	return nil
}
