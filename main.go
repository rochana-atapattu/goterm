package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/user/goterm/lib"
	"net/http"
)

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

	b := t.GetServers(i.Instanceid)
	cmd := t.CreateProxyCmd(b)

	go func() {
		t.CreateTerm(cmd)
	}()

	fmt.Fprint(w, "Welcome to my website!")
}
