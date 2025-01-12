package data

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Search struct {
	client *elasticsearch.Client
	index  string
}

func NewSearch(client *elasticsearch.Client) *Search {
	return &Search{
		client: client,
		index:  "tasks",
	}
}

func (i *Search) Query(ctx context.Context, params SearchParams) (SearchResults, error) {
	if params.IsZero() {
		return SearchResults{}, nil
	}

	should := make([]interface{}, 0, 3)

	if params.Description != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"description": *params.Description,
			},
		})
	}
	if params.Done != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"done": *params.Done,
			},
		})
	}
	if params.Priority != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"priority": *params.Priority,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	query["sort"] = []interface{}{
		"_score",
		map[string]interface{}{"due_at": "asc"},
	}
	query["from"] = params.From
	query["size"] = params.Size

	var buf bytes.Buffer

	if e := json.NewEncoder(&buf).Encode(query); e != nil {
		return SearchResults{}, e
	}

	request := esapi.SearchRequest{
		Index: []string{i.index},
		Body:  &buf,
	}

	response, e := request.Do(ctx, i.client)
	if e != nil {
		return SearchResults{}, e
	}
	defer response.Body.Close()

	if response.IsError() {
		return SearchResults{}, e
	}

	var results struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source Task `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if e := json.NewDecoder(response.Body).Decode(&results); e != nil {
		return SearchResults{}, e
	}

	hits := make([]Task, len(results.Hits.Hits))

	for i, hit := range results.Hits.Hits {
		hits[i].Description = hit.Source.Description
		hits[i].Done = hit.Source.Done
		hits[i].DueAt = hit.Source.DueAt
		hits[i].ID = hit.Source.ID
		hits[i].Priority = hit.Source.Priority
		hits[i].StartedAt = hit.Source.StartedAt
	}

	return SearchResults{
		Tasks: hits,
		Total: results.Hits.Total.Value,
	}, nil
}

type SearchParams struct {
	Description *string
	Done        *bool
	From        int64
	Priority    *Priority
	Size        int64
}

func (p SearchParams) IsZero() bool {
	return p.Description == nil && p.Done == nil && p.Priority == nil
}

type SearchResults struct {
	Tasks []Task
	Total int64
}
