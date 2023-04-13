package main

import (
	"flag"
	"fmt"
)

const (
	GROUP_VERSION = "app.grafana.com/v1"
	CRD_NAME      = "test-crds"
	NAMESPACE     = "default"

	TOTAL         = 100_000
	TEST_ID_COUNT = 100
	PAGE_SIZE     = 500
)

func main() {
	group_version := flag.String("group-version", GROUP_VERSION, "Group version")
	crd_name := flag.String("crd-name", CRD_NAME, "CRD name")
	namespace := flag.String("namespace", NAMESPACE, "Namespace")

	total_items := flag.Int("total", TOTAL, "Total number of objects in the namespace")
	test_id_count := flag.Int("test-id-count", TEST_ID_COUNT, "Number of ids to search for")
	page_size := flag.Int("page-size", PAGE_SIZE, "Number of items to request per page")

	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	base_url, token, err := getHostAndToken()
	if err != nil {
		panic(err)
	}

	next := ""
	items := [][]byte{}

	prefix := *crd_name + "-"
	ids := generateIds(prefix, *total_items, *test_id_count)

	// TODO: not sure if if this is a bad idea, but this request is reused for all pages
	req, err := newRequest(base_url, *group_version, *namespace, *crd_name, token, next, *page_size)
	if err != nil {
		panic(err)
	}

	for {
		body, err := fetchPage(req, next)
		if err != nil {
			panic(err)
		}

		var matchedItems [][]byte
		ids, matchedItems, next = filterResponse(body, ids)
		items = append(items, matchedItems...)
		if next == "" || len(ids) == 0 {
			break
		}
	}

	fmt.Printf("Found %d items", len(items))
}
