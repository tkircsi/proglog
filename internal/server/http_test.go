package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleProduce(t *testing.T) {
	pr := ProduceRequest{
		Record: Record{
			Value: []byte("Test Record-001"),
		},
	}

	bb, err := json.Marshal(pr)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(bb))
	if err != nil {
		t.Fatal(err)
	}

	newHTTPServer().handleProduce(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := strings.TrimRight(`{"offset":0}`, "\n")
	got := strings.TrimRight(string(body), "\n")
	if got != want {
		t.Errorf("want %s got %s", want, got)
	}
}

func TestHandleConsume(t *testing.T) {
	srv := newHTTPServer()
	srv.Log.Append(Record{
		Value:  []byte("Test Record-001"),
		Offset: 0,
	})

	cr := ConsumeRequest{
		Offset: 0,
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", bytes.NewReader(bb))
	if err != nil {
		t.Fatal(err)
	}

	srv.handleConsume(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := strings.TrimRight(`{"record":{"value":"VGVzdCBSZWNvcmQtMDAx","offset":0}}`, "\n")
	got := strings.TrimRight(string(body), "\n")
	if got != want {
		t.Errorf("want %s got %s", want, got)
	}
}

func TestHandleProduceSrv(t *testing.T) {
	srv := httptest.NewServer(NewRouter())
	defer srv.Close()

	pr := ProduceRequest{
		Record: Record{
			Value: []byte("Test Record-001"),
		},
	}

	bb, err := json.Marshal(pr)
	if err != nil {
		t.Fatal(err)
	}
	rs, err := http.Post(srv.URL+"/", "application/json", bytes.NewReader(bb))
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, rs.StatusCode)
	}

	want := strings.TrimRight(`{"offset":0}`, "\n")
	got := strings.TrimRight(string(body), "\n")
	if got != want {
		t.Errorf("want %s got %s", want, got)
	}
}

func TestHandleConsumeSrv(t *testing.T) {
	s := newHTTPServer()
	s.Log.Append(Record{
		Value:  []byte("Test Record-001"),
		Offset: 0,
	})

	srv := httptest.NewServer(NewRouter())
	defer srv.Close()

	cr := ConsumeRequest{
		Offset: 0,
	}

	bb, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/", bytes.NewReader(bb))
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	rs, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusNotFound {
		t.Errorf("want %d got %d", http.StatusOK, rs.StatusCode)
	}

}
