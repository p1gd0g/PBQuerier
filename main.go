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

var KSlice = []interface {
	Unmarshal([]byte) error
}{
	// Add new proto type here only.
	&tutorial.Person{},
	&tutorial.AddressBook{},
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
		Init(w)
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

	Execute(w, string(marshalled))
}

var KMap = map[string]interface {
	Unmarshal([]byte) error
}{}

var KString = []string{}

type data struct {
	Title string
	Items []string
	Out   string
}

func init() {
	for _, v := range KSlice {
		t := reflect.TypeOf(v)
		KString = append(KString, t.Elem().Name())
		KMap[t.Elem().Name()] = v
	}
	sort.Strings(KString)
}

//Init initiates the web.
func Init(w http.ResponseWriter) {

	d := data{
		Title: "My page",
		Items: KString,
	}

	parseAndExe(w, d)
}

// Execute writes the response.
func Execute(w http.ResponseWriter, out string) {
	d := data{
		Title: "My page",
		Items: KString,
		Out:   out,
	}

	parseAndExe(w, d)
}

func parseAndExe(w http.ResponseWriter, d data) {
	T, err := template.ParseFiles("query.html")
	if err != nil {
		log.Println(err)
	}

	err = T.Execute(w, d)
	if err != nil {
		log.Println(err)
	}
}
