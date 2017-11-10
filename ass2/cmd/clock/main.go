package main

import (
	"ass2"
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

func clock() {
	timenow := time.Now().Local()
	s := strings.Split(timenow.String(), " ")
	datenow := s[0]

	cu := &fotsopp.Currency{}

	//ticker := time.NewTicker(time.Hour * 24)
	//for range ticker.C {
	session, err := mgo.Dial(fotsopp.Durl)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := session.DB("assignment2").C("currency")
	err = c.Find(nil).Sort("-$natural").One(&cu)
	if err != nil {
		fmt.Println(err)
		return
	}
	if cu.Date != datenow {
		err = fotsopp.GetContent("http://api.fixer.io/latest", &cu)
		if err != nil {
			fmt.Println(err)
			return
		}
		cu.Date = datenow
		err = c.Insert(cu)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print("Ticker Updated\n")
	} else {
		fmt.Print("Ticker Not Updated\n")
	}
	fotsopp.CheckTrigger()
}

func main() {
	for {
		clock()
		time.Sleep(time.Hour * 24)
	}
}
