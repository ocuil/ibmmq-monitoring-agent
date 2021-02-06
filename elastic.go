package main

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	vault "github.com/sosedoff/ansible-vault-go"
)

func GetPassword() string {
	file, err := vault.DecryptFile("./conf/secrets.json", ".lab.")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return file
}

func loadESconfig() elasticsearch.Config {
	cfg := elasticsearch.Config{
		Addresses: []string{
			configuration.EShost,
		},
		Username: User,
		Password: Password,
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

}
