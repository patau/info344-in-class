package main

// import "os"
// import "log"
// import "net/http"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

//Struct: series of fields that defines an object that you want to manipulate as one unit
//Basically a Java object without the methods
type zip struct {
	Zip   string `json:"zip"`
	City  string `json:"city"`
	State string `json:"state"`
}

type zipSlice []*zip              //Slice of zips
type zipIndex map[string]zipSlice //Seattle --> slice of zips

//Capital Function --> exported; lowercase --> not exported
//In packages, capital b/c functions in packages need to be exported in order to be used

//* returns reference --> increases efficiency and ability to modify
//:= --> declare and set variable
func helloHandler(w http.ResponseWriter, r *http.Request) {
	//w = response writer to write back to client
	//Slice of bites?: In Java, arraylists; similar to this - can access and change
	name := r.URL.Query().Get("name") //r = request object

	//Add headers before sending response
	w.Header().Add("Content-Type", "text/plain")

	w.Write([]byte("hello " + name))
}

//HTTP handlers must have these two parameters
//Passed an address instead of the actual object b/c more efficient; 64bts passed instead of some larger object
//zi zipIndex = "this"
func (zi zipIndex) zipsForCityHandler(w http.ResponseWriter, r *http.Request) {
	//URL: /zips/city/seattle
	_, city := path.Split(r.URL.Path)
	lcity := strings.ToLower(city) //So we can still read it if user passed it in lowercase

	//Lets the client know that he's getting data in json/charset
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(zi[lcity]); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		log.Fatal("please set ADDR envrionment variable")
	}

	f, err := os.Open("./zips.json")
	if err != nil {
		log.Fatal("error opening zips file: " + err.Error())
	}

	zips := make(zipSlice, 0, 43000)
	decoder := json.NewDecoder(f)
	//& --> give me the address of this so I can pass it to something else
	if err := decoder.Decode(&zips); err != nil {
		log.Fatal("error decoding zips json: " + err.Error())
	}
	fmt.Printf("loaded %d zips\n", len(zips))

	zi := make(zipIndex)
	for _, z := range zips {
		lower := strings.ToLower(z.City)
		zi[lower] = append(zi[lower], z)
	}

	fmt.Printf("there are %d zips in Seattle\n", len(zi["seattle"]))

	//Params: resource path, pointer to function (when you put in the function's name)
	//By going to localhost:8000/hello --> calls helloHandler function, where w.Write is called
	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/zips/city/", zi.zipsForCityHandler)

	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
