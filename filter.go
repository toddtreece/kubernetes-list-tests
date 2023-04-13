package main

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

func filterResponse(body io.ReadCloser, ids []string) (filterdIds []string, matchedItems [][]byte, next string) {
	filterdIds = ids
	iter := jsoniter.Parse(jsoniter.ConfigFastest, body, 4096)
	for key := iter.ReadObject(); key != ""; key = iter.ReadObject() {
		switch key {
		case "items":
			matches, newIds := getMatches(iter, filterdIds)
			filterdIds = newIds
			matchedItems = append(matchedItems, matches...)
		case "metadata":
			next = getNextID(iter)
		default:
			iter.Skip()
		}
	}
	return filterdIds, matchedItems, next
}

func getMatches(iter *jsoniter.Iterator, ids []string) ([][]byte, []string) {
	found := []string{}
	items := [][]byte{}
	for iter.ReadArray() {
		raw := iter.SkipAndReturnBytes()
		item := jsoniter.ParseBytes(jsoniter.ConfigFastest, raw)
		for key := item.ReadObject(); key != ""; key = item.ReadObject() {
			switch key {
			case "metadata":
				name := getName(item)
				for _, id := range ids {
					if name == id {
						found = append(found, id)
						items = append(items, raw)
					}
				}
			default:
				item.Skip()
			}
		}
	}
	newids := []string{}
	for _, id := range ids {
		matched := false
		for _, foundId := range found {
			if id == foundId {
				matched = true
			}
		}
		if !matched {
			newids = append(newids, id)
		}
	}
	return items, newids
}

func getName(iter *jsoniter.Iterator) string {
	for key := iter.ReadObject(); key != ""; key = iter.ReadObject() {
		if key == "name" {
			return iter.ReadString()
		}
		iter.Skip()
	}
	return ""
}

func getNextID(iter *jsoniter.Iterator) string {
	for key := iter.ReadObject(); key != ""; key = iter.ReadObject() {
		if key == "continue" {
			return iter.ReadString()
		}
		iter.Skip()
	}
	return ""
}
