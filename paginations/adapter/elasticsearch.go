package adapter

import (
	"context"
	"encoding/json"
	"log"

	configs "github.com/crowdeco/skeleton/configs"
	elastic "github.com/olivere/elastic/v7"
	paginator "github.com/vcraescu/go-paginator/v2"
)

type (
	elasticsearchAdapter struct {
		context context.Context
		index   string
		query   elastic.Query
	}
)

func NewElasticsearchAdapter(context context.Context, index string, query elastic.Query) paginator.Adapter {
	return &elasticsearchAdapter{
		context: context,
		index:   index,
		query:   query,
	}
}

func (es *elasticsearchAdapter) Nums() (int64, error) {
	result, err := configs.Elasticsearch.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return 0, nil
	}

	return result.TotalHits(), nil
}

func (es *elasticsearchAdapter) Slice(offset, length int, data interface{}) error {
	result, err := configs.Elasticsearch.Search().Index(es.index).IgnoreUnavailable(true).Query(es.query).From(offset).Size(length).Do(es.context)
	if err != nil {
		log.Printf("%s", err.Error())
		return nil
	}

	records := data.(*[]interface{})
	var record interface{}
	for _, hit := range result.Hits.Hits {
		json.Unmarshal(hit.Source, &record)

		*records = append(*records, record)
	}

	data = *records

	return nil
}
