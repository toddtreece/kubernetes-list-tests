package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	GROUP_VERSION = "app.grafana.com/v1"
	CRD_NAME      = "test-crds"
	NAMESPACE     = "default"

	TOTAL         = 20_000
	TEST_ID_COUNT = 10
	PAGE_SIZE     = 500
)

type options struct {
	base_url string
	token    string

	time         bool
	cpu_profile  bool
	heap_profile bool

	group_version string
	crd_name      string
	name_prefix   string
	namespace     string

	total_items   int
	test_id_count int
	page_size     int
}

type runner func(*options) ([][]byte, error)

func main() {
	base_url, token, err := getHostAndToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	opt := &options{
		base_url: base_url,
		token:    token,
	}

	flag.BoolVar(&opt.cpu_profile, "cpu-profile", false, "Profile the program")
	flag.BoolVar(&opt.time, "time", false, "Output time stats")
	flag.BoolVar(&opt.heap_profile, "heap-profile", false, "Profile the program's heap")

	flag.StringVar(&opt.group_version, "group-version", GROUP_VERSION, "Group version")
	flag.StringVar(&opt.crd_name, "crd-name", CRD_NAME, "CRD name")
	flag.StringVar(&opt.name_prefix, "name-prefix", CRD_NAME+"-", "Name prefix")
	flag.StringVar(&opt.namespace, "namespace", NAMESPACE, "Namespace")

	flag.IntVar(&opt.total_items, "total", TOTAL, "Total number of objects in the namespace")
	flag.IntVar(&opt.test_id_count, "test-id-count", TEST_ID_COUNT, "Number of ids to search for")
	flag.IntVar(&opt.page_size, "page-size", PAGE_SIZE, "Number of items to request per page")

	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var r runner = run

	if opt.time {
		r = timed(r)
	}

	if opt.heap_profile {
		r = heap_profile(r)
	}

	if opt.cpu_profile {
		r = cpu_profile(r)
	}

	_, err = r(opt)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run(o *options) ([][]byte, error) {
	next := ""
	items := [][]byte{}

	ids := generateIds(o.name_prefix, o.total_items, o.test_id_count)

	// TODO: not sure if if this is a bad idea, but this request is reused for all pages
	req, err := newRequest(o.base_url, o.group_version, o.namespace, o.crd_name, o.token, next, o.page_size)
	if err != nil {
		return nil, err
	}

	for {
		body, err := fetchPage(req, next)
		if err != nil {
			return nil, err
		}

		var matchedItems [][]byte
		ids, matchedItems, next = filterResponse(body, ids)
		items = append(items, matchedItems...)

		if next == "" || len(ids) == 0 {
			return items, nil
		}
	}
}

func timed(fn runner) runner {
	return func(o *options) ([][]byte, error) {
		start := time.Now()
		items, err := fn(o)
		if err != nil {
			return nil, err
		}
		elapsed := time.Since(start)
		fmt.Printf("%s to fetch %d items\n", elapsed, len(items))
		return items, nil
	}
}

func heap_profile(fn runner) runner {
	return func(o *options) ([][]byte, error) {
		start, err := os.Create("out/heap-start.prof")
		if err != nil {
			return nil, err
		}
		defer start.Close()

		end, err := os.Create("out/heap-end.prof")
		if err != nil {
			return nil, err
		}
		defer end.Close()

		err = pprof.WriteHeapProfile(start)
		if err != nil {
			return nil, err
		}
		defer pprof.WriteHeapProfile(end)

		items, err := fn(o)
		return items, err
	}
}

func cpu_profile(fn runner) runner {
	return func(o *options) ([][]byte, error) {
		cpuProfileFile, err := os.Create("out/cpu.prof")
		if err != nil {
			return nil, err
		}
		runtime.SetCPUProfileRate(1000)

		err = pprof.StartCPUProfile(cpuProfileFile)
		if err != nil {
			return nil, err
		}
		defer func() {
			pprof.StopCPUProfile()
			cpuProfileFile.Close()
		}()
		items, err := fn(o)
		return items, err
	}
}
