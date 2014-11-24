package main

import (
    "log"
    "net/http"
    "github.com/dlutxx/goblin"
)

var settings goblin.Settings = goblin.Settings{
    "routes": map[string]interface{}{
        `hi/`: goblin.HandlerFromFunc(hello),
    },
    "handle404": goblin.HandlerFromFunc(handle404),
}

func handle404(res *goblin.ResponseWriter, req *http.Request) {
    // res.WriteHeader(http.StatusNotFound)
    n, e := res.WriteString(req.RequestURI + " dose not exist")
    log.Println("handle404", n, e, 1/(n-n))
}

func hello(res *goblin.ResponseWriter, req *http.Request) {
    res.WriteString("hello")
}

func main() {
    app, err := goblin.NewApp(settings)
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.ListenAndServe(":8096"))
}
