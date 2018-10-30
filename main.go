package main

import (
    "encoding/json"
    "fmt"
    "github.com/imroc/req"
    "os"
    "strings"
)

const (
    defaultHost = "192.168.1.63"
    defaultPort = "30002"
)

type Connection struct {
    Host string
    Port string
}

func conn() Connection {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")

    if len(host) == 0 || len(port) == 0 {
        return Connection{
            Host: defaultHost,
            Port: defaultPort,
        }
    } else {
        return Connection{
            Host: host,
            Port: port,
        }
    }
}

type Json struct {
    Serials []string `json:"serials"`
}

func (that Json) String() string {
    return fmt.Sprintf(`
{
    "serials": [
    %s
    ]
}`,
        strings.Join(that.Serials, ",\n\t\t"))
}

func main() {
    for {
        c := conn()
        url := fmt.Sprintf("http://%s:%s/", c.Host, c.Port)

        r, err := req.Get(url)
        if err != nil {
            panic(err)
        }

        var data Json
        err = json.Unmarshal(r.Bytes(), &data)
        if err != nil {
            panic(err)
        }

        fmt.Println(data)
    }
}
