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

// This function is for testing porpouse and validate the steps to the integrations
func elasticTest(metric []byte) {
	//fmt.Printf("%s\n", metric)
	fmt.Printf("End Metric\n")
}

func toLogstash() {

}

func toElastic(metric []byte) {
	cfg := loadESconfig()
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error en la conexion contra elasticsearch: %s\n", err)
	}

	// debug line
	//fmt.Println(string(metric))

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
	defer res.Body.Close()
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

func debugjsonReportStruct(j jsonReportStruct) {

	var doc elasticStruct = elasticStruct{}

	doc.TimeStamp = j.CollectionTime.TimeStamp
	doc.Epoch = j.CollectionTime.Epoch
	//doc.ObjectList.ObjectType.Tags = make(map[string]string)
	//doc.ObjectList.ObjectType.Metric = make(map[string]float64)
	//doc.ObjectList.Object = nil

	for _, value := range j.Points {

		/*
			for range value.ObjectType {
				fmt.Println(value.ObjectType)
			}
		*/

		//fmt.Println("Objeto ", value.ObjectType)

		if value.ObjectType == "queueManager" || value.ObjectType == "qmgr" {

			//var queueManager = make(map[string]string)

			var objectType = make(map[string]ObjectType)
			var tags = make(map[string]string)
			var metrics = make(map[string]float64)

			for tag, value := range value.Tags {
				//fmt.Println(tag, value)
				tags[tag] = value
			}

			for metric, value := range value.Metric {
				//fmt.Println(metric, value)
				metrics[metric] = value
			}

			objectType["queueManager"] = ObjectType{
				Tags:   tags,
				Metric: metrics,
			}

			doc.ObjectList.ObjectType = objectType

			//queueManager.Tags = map[interface{}]interface{}
			//doc.ObjectList.Object = append(doc.ObjectList.Object, objectType)

		}

		/*
			if value.ObjectType == "queue" {
				doc.ObjectList.ObjectType = value.ObjectType
				for tag, value := range value.Tags {
					//fmt.Println(tag, value)
					doc.ObjectList.Tags[tag] = value
				}
				for metricName, metricValue := range value.Metric {
					//fmt.Println(metricName, metricValue)
					doc.ObjectList.Metric[metricName] = metricValue
				}
			}
		*/
	}

	//fmt.Println(doc)
	jsonDoc, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Printf("%s\n", jsonDoc)

}
