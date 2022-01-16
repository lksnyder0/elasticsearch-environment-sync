package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

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
	app := cli.NewApp()
	app.Name = "elastiSync"
	app.Usage = "Sync configuration objects between elasticsearch clusters."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yml",
			Usage: "Load configuration from `FILE`",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Enable debug logging",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Error logging only",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("verbose") && c.Bool("quiet") {
			log.Fatal("--verbose and --quiet are mutually exclusive")
		} else if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else if c.Bool("quiet") {
			log.SetLevel(log.ErrorLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		config := getConfigObj(c.String("config"))
		_, _ = getElasticClient(config.Src_elastic)
		_, _ = getElasticClient(config.Dst_elastic)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
}

func getConfigObj(path string) sync_conf {
	var config sync_conf
	// Config parsing
	conf_cont, err := os.ReadFile(path)
	handleErr(err)
	err = yaml.Unmarshal(conf_cont, &config)
	handleErr(err)

	return config
}

func getElasticClient(conf elasticsearch.Config) (*elasticsearch.Client, map[string]interface{}) {
	es, err := elasticsearch.NewClient(conf)
	handleErr(err)
	es_info := getClusterInfo(*es)
	logNameAndInfo(es_info)
	return es, es_info
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
	log.Infof("Successfully connnected to %s", c_info["cluster_name"])
}
