package fotsopp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerProjects))
	ts2 := httptest.NewServer(http.HandlerFunc(HandlerWebhookId))

	defer ts.Close()
	defer ts2.Close()

	s := Webhook{"", "testurl", "EUR", "NOK", float64(1.55), float64(5.40)}

	json, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
		return
	}

	id, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = http.Get(ts2.URL + "/" + string(id))
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
		return
	}

	client := &http.Client{}

	del, err := http.NewRequest(http.MethodDelete, ts2.URL+"/"+string(id), nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
		return
	}

	_, err = client.Do(del)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
		return
	}

}

func Test_handlerLatest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerLatest))
	defer ts.Close()

	ma := make(map[string]interface{})

	ma["baseCurrency"] = "EUR"
	ma["targetCurrency"] = "NOK"

	json, err := json.Marshal(ma)
	if err != nil {
		fmt.Println(err)
	}

	_, err = http.Post(ts.URL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
		return
	}
}

func Test_handlerAverage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerAverage))
	defer ts.Close()

	ma := make(map[string]interface{})

	ma["baseCurrency"] = "EUR"
	ma["targetCurrency"] = "NOK"

	json, err := json.Marshal(ma)
	if err != nil {
		fmt.Println(err)
	}

	_, err = http.Post(ts.URL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}
}

func Test_handlerEvaluation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandlerEvaluation))
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}
}
