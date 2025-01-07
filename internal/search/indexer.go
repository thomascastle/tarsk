package search

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/thomascastle/tarsk/internal/data"
)

type TaskIndexer struct {
	client *elasticsearch.Client
	index  string
}

func NewTaskIndexer(client *elasticsearch.Client) *TaskIndexer {
	return &TaskIndexer{
		client: client,
		index:  "tasks",
	}
}

func (i *TaskIndexer) Delete(ctx context.Context, id string) error {
	request := esapi.DeleteRequest{
		Index:      i.index,
		DocumentID: id,
	}

	response, e := request.Do(ctx, i.client)
	if e != nil {
		return e
	}
	defer response.Body.Close()

	if response.IsError() {
		return e
	}

	io.Copy(io.Discard, response.Body)

	return nil
}

func (i *TaskIndexer) Index(ctx context.Context, task data.Task) error {
	var buf bytes.Buffer

	if e := json.NewEncoder(&buf).Encode(task); e != nil {
		return e
	}

	request := esapi.IndexRequest{
		Index:      i.index,
		Body:       &buf,
		DocumentID: task.ID,
		Refresh:    "true",
	}

	response, e := request.Do(ctx, i.client)
	if e != nil {
		return e
	}
	defer response.Body.Close()

	if response.IsError() {
		return e
	}

	io.Copy(io.Discard, response.Body)

	return nil
}
