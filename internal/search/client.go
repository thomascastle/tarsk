package search

import "github.com/elastic/go-elasticsearch/v7"

func NewClient() (*elasticsearch.Client, error) {
	client, e := elasticsearch.NewDefaultClient()
	if e != nil {
		return nil, e
	}

	response, e := client.Info()
	if e != nil {
		return nil, e
	}

	defer func() {
		response.Body.Close()
	}()

	return client, nil
}
