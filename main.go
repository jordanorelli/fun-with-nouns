package main

import (
    "fmt"
    "github.com/simonz05/godis"
    "io/ioutil"
    "strconv"
    "net/http"
)

var redis *godis.Client

func nounKey(id int) string {
    return "noun_" + strconv.Itoa(id)
}

func nounFile(id int) ([]byte, error) {
    filename := "/projects/noun_fountain/svg/" + strconv.Itoa(id) + ".svg"
    return ioutil.ReadFile(filename)
}

func nounHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[6:])
    if err != nil {
        http.Error(w, "Unable to parse integer from request url.",
                   http.StatusInternalServerError)
        return
    }
    key := nounKey(id)
    fmt.Println("Fetching key " + key)
    elem, err := redis.Get(key)
    if err != nil {
        fmt.Println("Unable to read svg " + strconv.Itoa(id) + " from redis.")
        http.Error(w, "Unable to read svg from redis",
                   http.StatusInternalServerError)
        return
    }
    if elem == nil {
        http.Error(w, "Got a null element.  What's that all about?",
                  http.StatusInternalServerError)
        return
    }
    w.Write(elem)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Index!");
    http.ServeFile(w, r, "/projects/noun_fountain/index.html")
}

func main() {
    redis = godis.New("", 0, "")

    for i := 0; i < 1535; i++ {
        body, err := nounFile(i)
        if err != nil {
            fmt.Println("Unable to read svg " + strconv.Itoa(i))
            continue
        }
        if err := redis.Set(nounKey(i), body); err != nil {
            fmt.Println("Unable to write svg " + strconv.Itoa(i) + " to redis.")
        }
    }

    http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("/projects/noun_fountain/js"))))

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/noun/", nounHandler)

    http.ListenAndServe(":8080", nil)
}
