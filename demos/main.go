package main

import (
    "log"
    "net/http"
    "github.com/dlutxx/goblin"
)

var settings goblin.Settings = goblin.Settings{
    "routes": map[string]interface{}{
        `hi/`: goblin.View(hello),
    },
}

func hello(res *goblin.ResponseWriter, req *http.Request) {
    log.Println("OK")
    res.Write([]byte("hello"))
}

func main() {
    app, err := goblin.NewApp(settings)
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.ListenAndServe(":8096"))
}
