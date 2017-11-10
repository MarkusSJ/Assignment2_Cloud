package fotsopp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

// Webhook struct for new webhook
type Webhook struct {
	Id             bson.ObjectId `json:"-" bson:"_id"`
	WebHookUrl     string        `json:"webhookURL"`
	BaseCurrency   string        `json:"baseCurrency"`
	TargetCurrency string        `json:"targetCurrency"`
	MinTrigger     float64       `json:"minTriggerValue"`
	MaxTrigger     float64       `json:"maxTriggerValue"`
}

// Currency struct for currency exchange
type Currency struct {
	Base  string                 `json:"base"`
	Date  string                 `json:"date"`
	Rates map[string]interface{} `json:"rates"`
}

var Durl = "mongodb://admin:123@ds233895.mlab.com:33895/assignment2"

func CheckTrigger() {

	var re []Webhook
	cu := &Currency{}

	session, err := mgo.Dial(Durl)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := session.DB("assignment2").C("webhook")
	defer session.Close()

	err = c.Find(nil).Sort("-$natural").All(&re)
	if err != nil {
		fmt.Println(err)
		return
	}
	c = session.DB("assignment2").C("currency")

	err = c.Find(nil).Sort("-$natural").One(&cu)
	if err != nil {
		fmt.Println(err)
		return
	}

	cu.Rates["EUR"] = float64(1)
	var target float64
	var base float64

	for i := range re {
		base = cu.Rates[re[i].BaseCurrency].(float64)
		target = cu.Rates[re[i].TargetCurrency].(float64)
		rate := target / base

		if rate > re[i].MaxTrigger || rate < re[i].MinTrigger {
			invokeWebhook(re[i], rate)
		}
	}
}

func invokeWebhook(ur Webhook, cur float64) {
	ma := make(map[string]interface{})

	ma["baseCurrency"] = ur.BaseCurrency
	ma["targetCurrency"] = ur.TargetCurrency
	ma["currentRate"] = cur
	ma["minTriggerValue"] = ur.MaxTrigger
	ma["maxTriggerValue"] = ur.MaxTrigger

	resp, err := json.Marshal(ma)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := http.Post(ur.WebHookUrl, "application/json", bytes.NewBuffer(resp))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(res.StatusCode)
	}
}

func HandlerProjects(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	ur := &Webhook{}

	session, err := mgo.Dial(Durl)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	c := session.DB("assignment2").C("webhook")
	defer session.Close()

	if r.Method == http.MethodPost {
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		json.Unmarshal(res, &ur)
		ur.Id = bson.NewObjectId()
		err = c.Insert(ur)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Write([]byte(ur.Id.Hex()))
	}
}

func HandlerWebhookId(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	ur := &Webhook{}

	argid := r.URL.Path
	s := strings.Split(argid, "/")
	urlid := s[len(s)-1]

	session, err := mgo.Dial(Durl)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	c := session.DB("assignment2").C("webhook")
	defer session.Close()

	if r.Method == http.MethodGet {
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = c.Find(bson.M{"_id": bson.ObjectIdHex(urlid)}).One(ur)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(ur)
	}

	if r.Method == http.MethodDelete {
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = c.Remove(bson.M{"_id": bson.ObjectIdHex(urlid)})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}

func GetContent(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func HandlerLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Header.Add(w.Header(), "content-type", "application/json")

		ur := &Webhook{}
		cu := &Currency{}

		session, err := mgo.Dial(Durl)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		c := session.DB("assignment2").C("currency")
		defer session.Close()

		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		json.Unmarshal(res, &ur)

		err = c.Find(nil).Sort("$natural: 1").One(&cu)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cu.Rates["EUR"] = float64(1)

		base := cu.Rates[ur.BaseCurrency].(float64)
		target := cu.Rates[ur.TargetCurrency].(float64)

		rate := target / base
		fmt.Fprintf(w, "%f", rate)

	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Using ticker as showed on https://gobyexample.com/tickers

func HandlerAverage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ur := &Webhook{}

		http.Header.Add(w.Header(), "content-type", "application/json")

		session, err := mgo.Dial(Durl)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		c := session.DB("assignment2").C("currency")
		defer session.Close()

		var cur []Currency

		err = c.Find(nil).Sort("-$natural").Limit(3).All(&cur)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		json.Unmarshal(res, &ur)

		var base float64
		var target float64

		for i := range cur {
			cur[i].Rates["EUR"] = float64(1)
			base += cur[i].Rates[ur.BaseCurrency].(float64)
			target += cur[i].Rates[ur.TargetCurrency].(float64)
		}

		rate := target / base
		fmt.Fprintf(w, "%f", rate)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func HandlerEvaluation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		var re []Webhook
		cu := &Currency{}

		session, err := mgo.Dial(Durl)
		if err != nil {
			fmt.Println(err)
			return
		}
		c := session.DB("assignment2").C("webhook")
		defer session.Close()

		err = c.Find(nil).Sort("-$natural").All(&re)
		if err != nil {
			fmt.Println(err)
			return
		}
		c = session.DB("assignment2").C("currency")

		err = c.Find(nil).Sort("-$natural").One(&cu)
		if err != nil {
			fmt.Println(err)
			return
		}
		cu.Rates["EUR"] = float64(1)
		fmt.Println(cu.Rates)

		var target float64
		var base float64

		for i := range re {
			base = cu.Rates[re[i].BaseCurrency].(float64)
			target = cu.Rates[re[i].TargetCurrency].(float64)
			rate := target / base

			invokeWebhook(re[i], rate)
		}

	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
