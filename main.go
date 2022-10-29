package main

import (
	"log"
	"os"

	"carrierleads.com/internal/crawler"
	"carrierleads.com/internal/dao"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBUrl        string `yaml:"db_url"`
	DOTWatermark int    `yaml:"dot_watermark"`
	BucketSize   int    `yaml:"bucket_size"`
}

func main() {
	config := readConfig()
	dao := dao.Instance(config.DBUrl)
	err := crawler.CrawlSafer(config.DOTWatermark, config.BucketSize, 50, dao)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() (c Config) {
	f, err := os.ReadFile("config.yaml") // just pass the file name
	if err != nil {
		log.Fatalf("failed to find config file %v", err)
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return
}
