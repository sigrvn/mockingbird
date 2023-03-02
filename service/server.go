package service

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
)

type dict map[string]any
type HandlerFunc = func(http.ResponseWriter, *http.Request)

type MockServer struct {
    config   Config
}

func NewMockServer(config Config) MockServer {
    for name, target := range config.Targets {
        handler := createHandlerFunc(name, target)
        http.HandleFunc("/" + name, handler)
    }

    m := MockServer{
        config:   config,
    }
    return m
}

func (s *MockServer) Listen() {
    fmt.Println("=> Running mock service...")
    port := ":" + strconv.Itoa(int(s.config.Port))
    fmt.Println("Listening on port:", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func createHandlerFunc(name string, target MockTarget) HandlerFunc {
    fmt.Println(" * Creating handler for endpoint: ", name)
    handler := func(w http.ResponseWriter, r *http.Request) {
        log.Print("got request for /", name)
        accept := false
        for _, method := range target.Methods {
            if r.Method == method {
                accept = true
                break
            }
        }

        w.Header().Set("Content-Type", "application/json")
        resp := json.NewEncoder(w)

        if !accept {
            w.WriteHeader(http.StatusMethodNotAllowed)
            err := resp.Encode(dict{
                "statusCode": http.StatusMethodNotAllowed,
                "error": "invalid method " + r.Method,
            })
            if err != nil {
                panic(err)
            }
            return
        }

        delay := time.Duration(target.Delay) * time.Millisecond
        time.Sleep(delay)

        err := resp.Encode(dict{
            "statusCode": http.StatusOK,
            "response": target.Response,
        })
        if err != nil {
            panic(err)
        }
    }
    return handler
}
