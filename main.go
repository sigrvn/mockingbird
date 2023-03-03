package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/sigrvn/mockingbird/service"
    "gopkg.in/yaml.v2"
)

func main() {
    //liveReload := flag.Bool("live-reload", false, "Reload the server on config change")
    flag.Parse()

    if len(flag.Args()) == 0 {
        fmt.Println("mockingbird error: no config file was provided")
        os.Exit(1)
    }

    configFile := flag.Args()[0]
    configData, err := os.ReadFile(configFile)
    if err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    var config service.Config
    if err = yaml.Unmarshal(configData, &config); err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    if err = config.Verify(); err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    m := service.NewMockService(config)
    m.Mock()
}
