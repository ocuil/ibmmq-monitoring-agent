package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
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
	fmt.Printf("%s\n", metric)
	fmt.Printf("End Metric")
}

func toLogstash() {

}

func toElastic() {
	cfg := loadESconfig()
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error en la conexion contra elasticsearch: ", err, "\n")
	}

	fmt.Printf(es)

}
