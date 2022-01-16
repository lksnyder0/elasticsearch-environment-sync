package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Structs
type obj_conf struct {
	Type string
	Name string
	Conf map[string]string
}

type sync_conf struct {
	Src_elastic elasticsearch.Config
	Dst_elastic elasticsearch.Config
	Items       []obj_conf
}

// Funcs
func main() {
	var config sync_conf
	app := cli.NewApp()
	app.Name = "elastiSync"

	// Config parsing
	conf_cont, err := os.ReadFile("config.yml")
	handleErr(err)
	err = yaml.Unmarshal(conf_cont, &config)
	handleErr(err)
	log.Printf("Connecting to source cluster at: %s", config.Src_elastic.Addresses)

	src_es, err := elasticsearch.NewClient(config.Src_elastic)
	handleErr(err)
	src_info := getClusterInfo(*src_es)
	logNameAndInfo(src_info)
	dst_es, err := elasticsearch.NewClient(config.Dst_elastic)
	handleErr(err)
	dst_info := getClusterInfo(*dst_es)
	logNameAndInfo(dst_info)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}

func getClusterInfo(es elasticsearch.Client) map[string]interface{} {
	var ret map[string]interface{}
	res, err := es.Info()
	handleErr(err)
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&ret)
	return ret
}

func logNameAndInfo(c_info map[string]interface{}) {
	log.Printf("Successfully connnected to %s", c_info["cluster_name"])
}
