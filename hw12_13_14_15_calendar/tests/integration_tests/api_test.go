package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

type getResponse struct {
	UUID string `json:"uuid"`
}

type apiTest struct {
	resp               *httptest.ResponseRecorder
	responseStatusCode int
	getResponseBody    []byte
	postResponseBody   []byte

	eventUUID uuid.UUID
}

func (test *apiTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodGet:
		r, err = http.Get(addr)
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseStatusCode = r.StatusCode

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	str, err := strconv.Unquote(string(bodyBytes))
	if err != nil {
		return
	}
	test.getResponseBody, err = base64.StdEncoding.DecodeString(str)
	return
}

func (test *apiTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *apiTest) iSendRequestToWithBody(httpMethod, addr string, body *messages.PickleStepArgument_PickleDocString) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(body.Content)
		r, err = http.Post(addr, "application/json", bytes.NewReader([]byte(cleanJson)))
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
	if ev[0].UUID != test.eventUUID.String() {
		return fmt.Errorf("event uuid %s received in get request doesn't match %s", ev[0].UUID, test.eventUUID.String())
	}
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
