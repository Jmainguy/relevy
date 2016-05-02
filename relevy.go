package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "github.com/ghodss/yaml"
    "net/http"
    "log"
)

type Config struct {
    Mongo_db string `json:"mongo_db"`
    Mongo_passwd string `json:"mongo_passwd"`
    Mongo_user string `json:"mongo_user"`
    Mongo_authdb string `json:"mongo_authdb"`
    Jsonstats string `json:"jsonstats"`
}

func config()  (mongo_db,mongo_passwd,mongo_user,mongo_authdb,jsonstats string){
    config_file, _ := ioutil.ReadFile("/etc/relevy/config.yaml")
    var v Config
    yaml.Unmarshal(config_file, &v)
    mongo_db = v.Mongo_db
    mongo_passwd = v.Mongo_passwd
    mongo_user = v.Mongo_user
    mongo_authdb = v.Mongo_authdb
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

//func mongo_stuff(mongo_db,mongo_passwd,mongo_user,mongo_authdb string) (

func main() {
    //Read Config, load values
    _,_,_,_,jsonstats := config()
    //mongo_db,mongo_passwd,mongo_user,mongo_authdb,jsonstats := config()
    mongo_db,_,_,_,jsonstats := config()
    // Initialize values, a string interface?
    values := make(map[string]interface{})
    //Load Json, if jsonstats is passed from config above, else move on without jsonstats
    if jsonstats != "" {
        fmt.Println(jsonstats)
        //jsonfile, _ := ioutil.ReadFile("/tmp/sample.json")
        jsonfile := json_grab(jsonstats)
        //Unpack Json so we can add things to it
        json.Unmarshal(jsonfile, &values)
    }
    //Unpack stuff from yaml into values as well
    yamlfile, _ := ioutil.ReadFile("/etc/relevy/info.yaml")
    y2, _ := yaml.YAMLToJSON(yamlfile)
    //fmt.Println(string(y2))
    json.Unmarshal(y2, &values)
    //Pack it all back up
    b, _ := json.Marshal(values)
    //Print it
    fmt.Println(string(b))
    message := mongo_db + "\n" + jsonstats
    fmt.Println(string(message))
}
