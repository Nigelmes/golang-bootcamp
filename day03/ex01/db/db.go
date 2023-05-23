package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strings"
)

type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Elasticsearch struct {
	Es *elasticsearch.Client
}

func NewElasticsearch(es *elasticsearch.Client) *Elasticsearch {
	return &Elasticsearch{es}
}

func (e *Elasticsearch) GetPlaces(limit int, offset int) ([]Place, int, error) {
	// Создаем поисковый запрос
	query := map[string]interface{}{
		"size": limit,
		"from": offset,
		// Добавьте необходимые условия фильтрации и сортировки
	}

	// Преобразуем запрос в JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	// Создаем поисковый запрос Elasticsearch
	req := esapi.SearchRequest{
		Index:          []string{"places"}, // Замените на имя вашего индекса
		Body:           strings.NewReader(string(queryJSON)),
		TrackTotalHits: true,
	}

	// Выполняем поиск в Elasticsearch
	res, err := req.Do(context.Background(), e.Es)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	// Проверяем статус ответа
	if res.IsError() {
		return nil, 0, fmt.Errorf("Elasticsearch search request failed: %s", res.String())
	}

	// Читаем ответ Elasticsearch
	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, 0, err
	}

	// Извлекаем список найденных документов (записей)
	hits := resBody["hits"].(map[string]interface{})["hits"].([]interface{})

	// Создаем слайс для хранения результатов
	places := make([]Place, 0, len(hits))

	// Итерируемся по найденным документам и преобразуем их в структуру types.Place
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		placeBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var place Place
		if err := json.Unmarshal(placeBytes, &place); err != nil {
			continue
		}

		places = append(places, place)
	}

	// Извлекаем общее количество записей из результата поиска
	totalHits := int(resBody["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	// Возвращаем список записей, общее количество записей и (или) ошибку
	return places, totalHits, nil
}
