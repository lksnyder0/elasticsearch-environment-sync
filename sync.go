package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	log "github.com/sirupsen/logrus"
)

func esResCheck(res esapi.Response) {
	if res.StatusCode != 200 {
		var res_map map[string]interface{}
		json.NewDecoder(res.Body).Decode(&res_map)
		status := int(res_map["status"].(float64))
		reason := res_map["error"].(map[string]interface{})["reason"]
		log.Fatalf("HTTP Status = %d, Error Reason = '%s'", status, reason)
	}
}

func syncILM(src elasticsearch.Client, dst elasticsearch.Client, config obj_conf) {
	body := getIlmBody(src, config.Name)
	putILM(dst, config.Name, body)
}

func getIlmBody(es elasticsearch.Client, id string) io.Reader {
	var req_map map[string]interface{}
	ilm_policy_body := make(map[string]interface{})
	log.Debugf("Getting ILM policy %s", id)
	req := esapi.ILMGetLifecycleRequest{
		Policy: id,
	}
	res, err := req.Do(context.Background(), es.Transport)
	handleErr(err)
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&req_map)
	// Need to clear out the extra fields
	ilm_policy_body["policy"] = req_map[id].(map[string]interface{})["policy"].(map[string]interface{})
	// Then turn it backk into a string
	body_bytes, err := json.Marshal(ilm_policy_body)
	handleErr(err)
	return bytes.NewReader(body_bytes)
}

func putILM(es elasticsearch.Client, id string, body io.Reader) {
	// var req_bytes []byte
	log.Debugf("Putting ILM policy %s", id)
	req := esapi.ILMPutLifecycleRequest{
		Body:   body,
		Policy: id,
	}
	res, err := req.Do(context.Background(), es.Transport)
	handleErr(err)
	esResCheck(*res)
	defer res.Body.Close()
	log.Debug("Success")
}
