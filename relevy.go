package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "log"
    "time"
    "os"
    "gopkg.in/mgo.v2"
    "strings"
    "github.com/ghodss/yaml"
)

type Config struct {
    Mongo_db string `json:"mongo_db"`
    Mongo_passwd string `json:"mongo_passwd"`
    Mongo_user string `json:"mongo_user"`
    Mongo_authdb string `json:"mongo_authdb"`
    Mongo_addr string `json:"mongo_addr"`
    Jsonstats string `json:"jsonstats"`
}

func config()  (mongo_db,mongo_passwd,mongo_user,mongo_authdb,mongo_addr,jsonstats string){
    var v Config
    config_file, err := ioutil.ReadFile("/etc/relevy/config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    yaml.Unmarshal(config_file, &v)
    mongo_db = v.Mongo_db
    mongo_passwd = v.Mongo_passwd
    mongo_user = v.Mongo_user
    mongo_authdb = v.Mongo_authdb
    mongo_addr = v.Mongo_addr
    jsonstats = v.Jsonstats
    return
}

func replace_key(record map[string]interface {}, key string, value interface {}) {
    delete(record, key)
    new_key := strings.Replace(key, ".", "-", -1)
    record[new_key] = value
}

func main() {
    // For loop to keep it running forever
    for {
        // Read Config, load values
        mongo_db,mongo_passwd,mongo_user,mongo_authdb,mongo_addr,jsonstats := config()
        // We need this object to establish a session to our MongoDB.
        mongoDBDialInfo := &mgo.DialInfo{
          Addrs:    []string{mongo_addr},
          Timeout:  60 * time.Second,
          Database: mongo_authdb,
          Username: mongo_user,
          Password: mongo_passwd,
        }

        // Create a session which maintains a pool of socket connections
        // to our MongoDB.
        mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
          log.Fatalf("CreateSession: %s\n", err)
        }

        // Initialize values, a map of strings
        var values map[string]interface{}
        // Load Json, if jsonstats is passed from config above, else move on without jsonstats
        if jsonstats != "" {
            resp, err := http.Get(jsonstats)
            // Bomb out if http.Get fails
            if err != nil {
                log.Fatal(err)
            }
            json.NewDecoder(resp.Body).Decode(&values)
            resp.Body.Close()
        }
        // Unpack stuff from yaml into values as well
        data, err := ioutil.ReadFile("/etc/relevy/info.yaml")
        if err != nil {
            log.Fatal(err)
        }
        yaml.Unmarshal(data, &values)
        // Hostname and time
        hostname, err := os.Hostname()
        if err != nil {
            log.Fatal(err)
        }
        values["_id"] = &hostname
        values["Updated"] = time.Now()
        // Remove dots from keys
        for key, value := range values {
            if strings.ContainsAny(key, ".") {
                replace_key(values, key, value)
            }
            // If a map inside a map, remove dots from their keys as well
            if rec, ok := value.(map[string]interface{}); ok {
                for skey, sval := range rec {
                    if strings.ContainsAny(skey, ".") {
                        replace_key(rec, skey, sval)
                    }
                }
            }
        }
        // Request a socket connection
        sessionCopy := mongoSession.Copy()
        // Close session whn goroutine exits
        // Add into Mongo
        coll := sessionCopy.DB(mongo_db).C("relevy")
        _ ,err = coll.UpsertId(&hostname, values)
        if err != nil {
            log.Fatal(err)
        }

        // Close session
        mongoSession.Close()
        sessionCopy.Close()
        fmt.Println("Loop complete, everything is working")
        time.Sleep(5 * time.Second)
    }
}
