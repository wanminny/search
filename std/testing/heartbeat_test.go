package testing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
)


type testGetHandler struct {
	JsonReturn string
	HttpError  string
	HttpCode   int
}

func (t testGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.HttpError != "" {
		http.Error(w, t.HttpError, t.HttpCode)
	} else {
		w.Write([]byte(t.JsonReturn))
	}
}

func TestGet(t *testing.T) {

	var TestValues = []struct {
		Json      string
		HttpError string
		HttpCode  int
		RetMsg    HeartbeatMessage
		RetErr    string
	}{
		{
			`{"status":"running","build":"testHash","uptime":"5m31.5s"}`,
			"",
			http.StatusOK,
			HeartbeatMessage{"running", "testHash", "5m31.5s"},
			"",
		},
		{
			"",
			"Error",
			http.StatusBadRequest,
			HeartbeatMessage{},
			"Wrong status code: 400",
		},
		{
			"",
			"",
			http.StatusOK,
			HeartbeatMessage{},
			"Error occured unmarshalling the response",
		},
	}

	for _, tv := range TestValues {

		ts := httptest.NewServer(testGetHandler{tv.Json, tv.HttpError, tv.HttpCode})
		defer ts.Close()

		hm, err := Get(ts.URL)
		if err != nil && err.Error() != tv.RetErr {
			t.Fatal("Wrong error result! Expected:", tv.RetErr, "Got:", err)
		}

		if !reflect.DeepEqual(hm, tv.RetMsg) {
			t.Fatal("Wrong result object! Expected:", tv.RetMsg, "Got:", hm)
		}
	}
}
