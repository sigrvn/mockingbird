package service

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
)

type dict map[string]any
type HandlerFunc = func(http.ResponseWriter, *http.Request)

type MockService struct {
    config Config
    logger *log.Logger
}

func NewMockService(config Config) MockService {
    const loggerPrefix = "[mockingbird] | "
    logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags | log.Lmsgprefix)

    m := MockService{config, logger}
    m.logger.Println("Creating Mock Service")
    m.logger.Println("=> Generating handlers for API(s)...")

    for name, target := range config.Targets {
        handler := m.createHandlerFunc(name, target)
        http.HandleFunc("/" + name, handler)
    }

    return m
}

func (m *MockService) Mock() {
    m.logger.Println("=> Running mock service...")
    port := ":" + strconv.Itoa(int(m.config.Port))
    m.logger.Println("Listening on port:", port)
    m.logger.Fatal(http.ListenAndServe(port, nil))
}

func (m *MockService) createHandlerFunc(name string, target MockTarget) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m.logger.Print("got request for /", name)
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

        err := resp.Encode(target.Response)
        if err != nil {
            panic(err)
        }
    }
}
