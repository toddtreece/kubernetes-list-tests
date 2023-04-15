package main

import (
	"encoding/json"
	"io"
)

type metadata struct {
	Name string `json:"name"`
}
type item struct {
	Metadata metadata        `json:"metadata"`
	Spec     json.RawMessage `json:"spec"`
}

type listmeta struct {
	Next string `json:"continue"`
}

type response struct {
	Items    []item   `json:"items"`
	Metadata listmeta `json:"metadata"`
}

func filterResponse(body io.ReadCloser, ids []string) (filteredIds []string, matchedItems [][]byte, next string) {
	res := response{}
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		panic(err)
	}

	next = res.Metadata.Next

	for _, id := range ids {
		found := false
		for _, item := range res.Items {
			if item.Metadata.Name == id {
				found = true
				matchedItems = append(matchedItems, item.Spec)
			}
		}
		if !found {
			filteredIds = append(filteredIds, id)
		}
	}

	return filteredIds, matchedItems, next
}
