package main

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"

	"PBQuerier/tutorial"
)

// KSlice is the key slice.
var KSlice = []interface {
	Unmarshal([]byte) error
}{
	// Add new proto type here only.
	&tutorial.Person{},
	&tutorial.AddressBook{},
}

// KMap is the key map.
var KMap = map[string]interface {
	Unmarshal([]byte) error
}{}

// KString is the key string.
var KString = []string{}

func init() {
	for _, v := range KSlice {
		t := reflect.TypeOf(v)
		KString = append(KString, t.Elem().Name())
		KMap[t.Elem().Name()] = v
	}
	sort.Strings(KString)
}

func main() {
	log.SetFlags(log.Lshortfile)

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8901", nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {

	urlVal := req.URL.Query()

	if urlVal.Encode() == "" {
		execute(w, "")
		return
	}

	d := KMap[urlVal.Get("key")]

	atob, err := base64.StdEncoding.DecodeString(urlVal.Get("proto"))
	if err != nil {
		log.Println(err)
	}

	proto, err := strconv.Unquote("\"" + string(atob) + "\"")
	if err != nil {
		log.Println(err)
	}

	err = d.Unmarshal([]byte(proto))
	if err != nil {
		log.Println(err)
	}

	marshalled, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Println(err)
	}

	execute(w, string(marshalled))
}

// Execute writes the response.
func execute(w http.ResponseWriter, out string) {

	d := struct {
		Title string
		Items []string
		Out   string
	}{
		Title: "My page",
		Items: KString,
		Out:   out,
	}

	T, err := template.ParseFiles("query.html")
	if err != nil {
		log.Println(err)
	}

	err = T.Execute(w, d)
	if err != nil {
		log.Println(err)
	}
}
