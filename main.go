package main

import (
	"encoding/json"
	"fmt"
	"github.com/rochana-atapattu/goterm/termlib"
	"io/ioutil"
	"net/http"
)

type instance struct {
	Instanceid string `json:"instanceid"`
}

func main() {
	http.HandleFunc("/", TestHandler)
	http.ListenAndServe(":8080", nil)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {

	i := instance{}

	fmt.Println("got request")
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &i)
	if err != nil {
		fmt.Println(err)
		return
	}

	termlib.StartPty(i.Instanceid)

	fmt.Fprint(w, "Welcome to my website!")
	return
}
