package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	vault "github.com/sosedoff/ansible-vault-go"
)

func GetPassword() string {
	file, err := vault.DecryptFile("./conf/secrets.json", ".lab.")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return file
}

func loadVars() {
	file, err := os.Open("conf/vars.json")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func loadSecrets() {
	getSecrets := json.NewDecoder(strings.NewReader(GetPassword()))
	err := getSecrets.Decode(&secrets)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func loadESconfig() elasticsearch.Config {
	transport := http.DefaultTransport
	tlsClientConfig := &tls.Config{InsecureSkipVerify: true}
	transport.(*http.Transport).TLSClientConfig = tlsClientConfig
	cfg := elasticsearch.Config{
		Addresses: []string{
			configuration.EShost,
		},
		Username: secrets.USER_elastic,
		Password: secrets.PWD_elastic,
	}
	return cfg
}

func toLogstash() {

}

func toElastic(doc elasticStruct) {
	metric, _ := json.MarshalIndent(doc, "", "  ")

	fmt.Printf("%s\n", metric)

	cfg := loadESconfig()
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error en la conexion contra elasticsearch: %s\n", err)
	}

	currentTime := time.Now()
	fecha := currentTime.Format("2006.01.02")
	index := "ibmmq.v1." + fecha
	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(string(metric)),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	res.Body.Close()
}

type ObjectType struct {
	Tags   map[string]string  `json:"tags"`
	Metric map[string]float64 `json:"metrics"`
}

type elasticStruct struct {
	TimeStamp  string `json:"timeStamp"`
	Epoch      int64  `json:"epoch"`
	ObjectList struct {
		ObjectType map[string]ObjectType `json:"objectType"`
	} `json:"objectList"`
}

func createJSON(j jsonReportStruct) {

	var doc elasticStruct = elasticStruct{}

	doc.TimeStamp = j.CollectionTime.TimeStamp
	doc.Epoch = j.CollectionTime.Epoch

	for _, value := range j.Points {

		objectType := make(map[string]ObjectType)
		tags := make(map[string]string)
		metrics := make(map[string]float64)

		for tag, value := range value.Tags {
			tags[tag] = value
		}
		for metric, value := range value.Metric {
			metrics[metric] = value
		}
		objectType[value.ObjectType] = ObjectType{
			Tags:   tags,
			Metric: metrics,
		}
		doc.ObjectList.ObjectType = objectType

		if configuration.SendTO == "elastic" {
			go toElastic(doc)
		}

	}

}
