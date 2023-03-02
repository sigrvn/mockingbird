package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/sigrvn/mockingbird/service"
    "gopkg.in/yaml.v2"
)

var config service.Config

func main() {
    //liveReload := flag.Bool("live-reload", false, "Reload the server on config change")
    flag.Parse()

    if len(flag.Args()) == 0 {
        fmt.Println("mockingbird error: no config file was provided")
        os.Exit(1)
    }

    configFile := flag.Args()[0]
    // TODO: Allow the config filename to be named something different and be passed in via cmdline args
    configData, err := os.ReadFile(configFile)
    if err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    if err = yaml.Unmarshal(configData, &config); err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    if err = config.Verify(); err != nil {
        fmt.Printf("mockingbird error: %s\n", err.Error())
        os.Exit(1)
    }

    server := service.NewMockServer(config)
    server.Listen()
}
