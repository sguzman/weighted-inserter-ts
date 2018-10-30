package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/imroc/req"
    _ "github.com/lib/pq"
    "math/rand"
    "os"
    "strings"
)

const (
    defaultHost = "192.168.1.63"
    defaultPort = "30002"
    defaultDBHost = "192.168.1.63"
    defaultDBPort = "30000"
)

type Connection struct {
    Host string
    Port string
}

func conn() Connection {
    host := os.Getenv("HOST")
    port := os.Getenv("PORT")

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

func connStr() string {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")

    if len(host) == 0 || len(port) == 0 {
        return fmt.Sprintf("user=postgres dbname=youtube host=%s port=%s sslmode=disable", defaultDBHost, defaultDBPort)
    } else {
        return fmt.Sprintf("user=postgres dbname=youtube host=%s port=%s sslmode=disable", host, port)
    }
}

func connection() *sql.DB {
    db, err := sql.Open("postgres", connStr())
    if err != nil {
        panic(err)
    }

    return db
}

func getKey() string {
    rawKey := os.Getenv("API_KEY")
    splitKeys := strings.Split(rawKey, "|")

    return splitKeys[rand.Intn(len(splitKeys))]
}

func get(cs []string) []byte {
    key := getKey()
    url := "https://www.googleapis.com/youtube/v3/channels"
    partStr := "statistics"
    idStr := strings.Join(cs, ",")

    param := req.Param{
        "part":  partStr,
        "id": idStr,
        "key": key,
    }

    r, err := req.Get(url, param)
    if err != nil {
        panic(err)
    }

    str, err := r.ToBytes()
    if err != nil {
        panic(err)
    }

    return str
}

type PageInfo struct {
    Results int `json:"totalResults"`
    Per int `json:"resultsPerPage"`
}

type StatsType struct {
    ViewCount string `json:"viewCount"`
    CommentCount string `json:"commentCount"`
    SubscriberCount string `json:"subscriberCount"`
    HiddenSubscriberCount bool `json:"hiddenSubscriberCount"`
    VideoCount string `json:"videoCount"`
}

type ItemType struct {
    Kind string `json:"kind"`
    Etag string `json:"etag"`
    Id string `json:"id"`
    Statistics StatsType `json:"statistics"`
}

type DataType struct {
    Kind string `json:"kind"`
    Etag string `json:"etag"`
    Token string `json:"nextPageToken"`
    Info PageInfo `json:"pageInfo"`
    Items []ItemType `json:"items"`
}

func getData(cs []string) DataType {
    jsonBytes := get(cs)
    fmt.Println(string(jsonBytes))

    var data DataType
    err := json.Unmarshal(jsonBytes, &data)
    if err != nil {
        panic(err)
    }

    return data
}

func insert(ds DataType) {
    db := connection()
    defer func() {
        err := db.Close()
        if err != nil {
            panic(err)
        }
    }()

    sqlInsert := "INSERT INTO youtube.entities.chan_stats (serial, subs, videos, views) VALUES ($1, $2, $3, $4)"

    for i := range ds.Items {
        d := ds.Items[i]

        serial := d.Id
        subs := d.Statistics.SubscriberCount
        videos := d.Statistics.VideoCount
        views := d.Statistics.ViewCount

        _, err := db.Exec(sqlInsert, serial, subs, videos, views)
        if err != nil {
            panic(err)
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

func stuff() {
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

    datas := getData(data.Serials)
    insert(datas)
}

func main() {
    for {
        stuff()
    }
}
