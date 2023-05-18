// curl -s -XGET "http://localhost:9200/places/_doc/0"  -- проверяем через _doc а не через place
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	res             *esapi.Response
	countSuccessful uint64
	err             error
)

const (
	csvFileName     = "data.csv"
	indexName       = "places"
	mappingFileName = "schema.json"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	mapping, err := readMappingFromFile()
	if err != nil {
		log.Fatal(err)
	}
	createIndexandMapping(es, mapping)
	data, err := parseCsvFile()
	if err != nil {
		log.Fatal(err)
	}
	loadDataIntoElastic(es, data)
}

func readMappingFromFile() (string, error) {
	res, err := os.ReadFile(mappingFileName)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func createIndexandMapping(es *elasticsearch.Client, mapping string) {
	if res, err = es.Indices.Delete(
		[]string{indexName}, es.Indices.Delete.WithIgnoreUnavailable(true),
	); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err) // удаляем если уже существует такой индекс
	}
	defer res.Body.Close()
	res, err = es.Indices.Create( // создаем новый индекс и добавляем маппинг
		indexName,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	defer res.Body.Close()
}

func loadDataIntoElastic(es *elasticsearch.Client, datas []Data) {
	bi, err := esutil.NewBulkIndexer(
		esutil.BulkIndexerConfig{
			Index:         indexName,        // The default index name
			Client:        es,               // The Elasticsearch client
			NumWorkers:    8,                // The number of worker goroutines
			FlushBytes:    10000,            // The flush threshold in bytes
			FlushInterval: 30 * time.Second, // The periodic flush interval
		},
	)
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}
	for _, a := range datas {
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", a.ID, err)
		}
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: a.ID,
				Body:       bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
