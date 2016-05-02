package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "github.com/ghodss/yaml"
    "net/http"
    "log"
    "time"
    "os"
    "gopkg.in/mgo.v2"
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
    config_file, _ := ioutil.ReadFile("/etc/relevy/config.yaml")
    var v Config
    yaml.Unmarshal(config_file, &v)
    mongo_db = v.Mongo_db
    mongo_passwd = v.Mongo_passwd
    mongo_user = v.Mongo_user
    mongo_authdb = v.Mongo_authdb
    mongo_addr = v.Mongo_addr
    jsonstats = v.Jsonstats
    return
}

func json_grab(url string) (jsonstats_json []byte) {
    resp, err := http.Get(url)
    //Bomb out if htt.Get fails
    if err != nil {
        log.Fatal(err)
    }
    jsonstats_json, _ = ioutil.ReadAll(resp.Body)
    return jsonstats_json
}

func main() {
    // For loop to keep it running forever
    for {
        //Read Config, load values
        //mongo_db,_,_,_,mongo_addr,jsonstats := config()
        mongo_db,mongo_passwd,mongo_user,mongo_authdb,mongo_addr,jsonstats := config()
        //mongo_db,_,_,_,jsonstats := config()

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

        // Initialize values, a string interface?
        values := make(map[string]interface{})
        //Load Json, if jsonstats is passed from config above, else move on without jsonstats
        if jsonstats != "" {
            //jsonfile, _ := ioutil.ReadFile("/tmp/sample.json")
            jsonfile := json_grab(jsonstats)
            //Unpack Json so we can add things to it
            json.Unmarshal(jsonfile, &values)
        }
        //Unpack stuff from yaml into values as well
        yamlfile, _ := ioutil.ReadFile("/etc/relevy/info.yaml")
        y2, _ := yaml.YAMLToJSON(yamlfile)
        json.Unmarshal(y2, &values)
        // Hostname and time
        hostname, _ := os.Hostname()
        values["_id"] = &hostname
        values["Updated"] = time.Now()
        // Add tnto Mongo
        coll := mongoSession.DB(mongo_db).C("relevy")
        //coll.UpsertId(values["Updated"], values)
        _, err2 := coll.UpsertId(&hostname, values)
        fmt.Println(err2)
        //coll.Insert(values)
       // Pack it all back up
        //b, _ := json.Marshal(values)
        //Print it
        //fmt.Println(string(b))

        //time.Sleep(1 * time.Minute)
        time.Sleep(1 * time.Second)
    }
}
