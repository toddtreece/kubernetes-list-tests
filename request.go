package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func newRequest(base_url, group_version, namespace, crd_name, token, next string, page_size int) (*http.Request, error) {
	u, err := url.Parse(base_url)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("/apis/%s/namespaces/%s/%s", group_version, namespace, crd_name)
	query := u.Query()
	query.Add("limit", fmt.Sprintf("%d", page_size))
	if next != "" {
		query.Add("continue", next)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func fetchPage(req *http.Request, next string) (io.ReadCloser, error) {
	if next != "" {
		query := req.URL.Query()
		query.Set("continue", next)
		req.URL.RawQuery = query.Encode()
	}
	fmt.Printf("Requesting %s\n", req.URL.String())
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
