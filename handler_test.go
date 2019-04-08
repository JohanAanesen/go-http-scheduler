package main

import (
	"bytes"
	"encoding/json"
	"github.com/JohanAanesen/CSAMS/schedulerservice/model"
	"github.com/JohanAanesen/CSAMS/webservice/shared/util"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var dummyPayload = model.Payload{
	Authentication: os.Getenv("PEER_AUTH"),
	ScheduledTime:  util.GetTimeInCorrectTimeZone().Add(time.Hour * 24 * 31), // TODO time-norwegian
	Task:           "peer",
	AssignmentID:   1,
	Data: peerTasktoRawJSON(model.PeerTask{
		Authentication: os.Getenv("PEER_AUTH"),
		AssignmentID:   1,
		Reviewers:      20000000,
	}),
}

var dummyUpdate = struct {
	Authentication string    `json:"authentication"`
	AssignmentID   int       `json:"assignment_id"`
	ScheduledTime  time.Time `json:"scheduled_time"`
}{
	Authentication: os.Getenv("PEER_AUTH"),
	AssignmentID:   1,
	ScheduledTime:  util.GetTimeInCorrectTimeZone().Add(time.Hour * 24 * 31), // TODO time-norwegian
}

var dummyDelete = struct {
	Authentication string `json:"authentication"`
	AssignmentID   int    `json:"assignment_id"`
}{
	Authentication: os.Getenv("PEER_AUTH"),
	AssignmentID:   1,
}

func getReaderFromPayload(payload model.Payload) io.Reader {
	body, _ := json.Marshal(payload)
	return bytes.NewReader(body)
}

func getReaderFromUpdate(payload struct {
	Authentication string    `json:"authentication"`
	AssignmentID   int       `json:"assignment_id"`
	ScheduledTime  time.Time `json:"scheduled_time"`
}) io.Reader {
	body, _ := json.Marshal(payload)
	return bytes.NewReader(body)
}

func getReaderFromDelete(payload struct {
	Authentication string `json:"authentication"`
	AssignmentID   int    `json:"assignment_id"`
}) io.Reader {
	body, _ := json.Marshal(payload)
	return bytes.NewReader(body)
}

func peerTasktoRawJSON(peerTask model.PeerTask) json.RawMessage {
	byte, _ := json.Marshal(peerTask)

	return byte
}

func TestHandlers(t *testing.T) {
	Initialize()

	tests := []struct {
		name    string
		method  string
		url     string
		body    io.Reader
		handler func(w http.ResponseWriter, r *http.Request)

		expectedCode int
	}{
		{
			name:         "indexGET",
			method:       "GET",
			url:          "/",
			body:         nil,
			handler:      IndexGET,
			expectedCode: http.StatusOK,
		},
		{
			name:         "indexPOST_empty_body",
			method:       "POST",
			url:          "/",
			body:         nil,
			handler:      IndexPOST,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "indexPOST",
			method:       "POST",
			url:          "/",
			body:         getReaderFromPayload(dummyPayload),
			handler:      IndexPOST,
			expectedCode: http.StatusOK,
		},
		{
			name:         "indexPUT",
			method:       "PUT",
			url:          "/",
			body:         getReaderFromUpdate(dummyUpdate),
			handler:      IndexPUT,
			expectedCode: http.StatusOK,
		},
		{
			name:         "indexDELETE",
			method:       "DELETE",
			url:          "/",
			body:         getReaderFromDelete(dummyDelete),
			handler:      IndexDELETE,
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		r, _ := http.NewRequest(test.method, test.url, test.body)
		w := httptest.NewRecorder()

		test.handler(w, r)

		t.Run(test.name, func(t *testing.T) {
			if w.Body.String() == "" {
				t.Logf("error, response body was empty")
				t.Fail()
			}

			// Check status code
			if w.Code != test.expectedCode {
				t.Logf("expected: %v, got: %v\n", test.expectedCode, w.Code)
				t.Fail()
			}
		})
	}
}
