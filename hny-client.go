package main

import (
	"encoding/json"
	"net/http"
)

type HoneycombClient struct {
	baseUrl string
	client  *http.Client
}

type honeycombTransport struct {
	apiKey string
}

func (t *honeycombTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Honeycomb-Team", t.apiKey)
	return http.DefaultTransport.RoundTrip(req)
}

func NewHoneycombClient(apiKey string) *HoneycombClient {
	return &HoneycombClient{
		baseUrl: "https://api.honeycomb.io",
		client:  &http.Client{Transport: &honeycombTransport{apiKey: apiKey}},
	}
}

func (c *HoneycombClient) GetBoard(boardId string) (*HoneycombBoard, error) {

	url := c.baseUrl + "/1/boards/" + boardId
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var board HoneycombBoard
	if err := json.NewDecoder(resp.Body).Decode(&board); err != nil {
		return nil, err
	}

	return &board, nil
}

func (c *HoneycombClient) GetQuery(dataset string, queryId string) (*HoneycombQuery, error) {

	url := c.baseUrl + "/1/queries/" + dataset + "/" + queryId
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var query HoneycombQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, err
	}

	return &query, nil
}

func (c *HoneycombClient) GetQueryAnnotation(dataset string, queryAnnotationId string) (*HoneycombQueryAnnotation, error) {
	url := c.baseUrl + "/1/query_annotations/" + dataset + "/" + queryAnnotationId
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queryAnnotation HoneycombQueryAnnotation
	if err := json.NewDecoder(resp.Body).Decode(&queryAnnotation); err != nil {
		return nil, err
	}

	return &queryAnnotation, nil
}
