package service

import (
    "fmt"
    "net/http"
)

type MockTarget struct {
    Methods  []string       `yaml:"methods"`
    Delay    int            `yaml:"delay"`
    Response dict `yaml:"response"`
}

type Config struct {
    Port    uint16
    Targets map[string]MockTarget `yaml:"targets"`
}

var (
    validMethods = map[string]uint8{
        http.MethodGet:     0,
        http.MethodHead:    0,
        http.MethodPost:    0,
        http.MethodPut:     0,
        http.MethodPatch:   0,
        http.MethodDelete:  0,
        http.MethodConnect: 0,
        http.MethodOptions: 0,
        http.MethodTrace:   0,
    }
)

func (c *Config) Verify() error {
    for name, target := range c.Targets {
        for _, method := range target.Methods {
            if _, isValid := validMethods[method]; !isValid {
                return fmt.Errorf("invalid method '%s' for target '%s'", method, name)
            }
        }
    }

    return nil
}
