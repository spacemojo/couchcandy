package couchcandy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

// defaultGet is default method with explicit "GET" and &http.Client{}
func defaultGet(url string) (*http.Response, error) {
	return defaultMethod(http.MethodGet, url, &http.Client{})
}

// defaultDelete is default method with explicit "DELETE" and &http.Client{}
func defaultDelete(url string) (*http.Response, error) {
	return defaultMethod(http.MethodDelete, url, &http.Client{})
}

// defaultMethod is for GET and DELETE statements
func defaultMethod(method, url string, client CandyHTTPClient) (*http.Response, error) {

	request, requestError := http.NewRequest(method, url, nil)
	if requestError != nil {
		return nil, requestError
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func defaultPostJSON(url, body string) (*http.Response, error) {
	return defaultJSONWithBody(http.MethodPost, url, body, &http.Client{})
}

func defaultPutJSON(url, body string) (*http.Response, error) {
	return defaultJSONWithBody(http.MethodPut, url, body, &http.Client{})
}

func defaultJSONWithBody(method, url, body string, client CandyHTTPClient) (*http.Response, error) {

	bodyJSON := strings.NewReader(body)
	request, requestError := http.NewRequest(method, url, bodyJSON)
	if requestError != nil {
		return nil, requestError
	}

	request.Header.Add(HeaderContentType, JSONContentType)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func defaultPutBytes(url, contentType string, body []byte) (*http.Response, error) {
	return defaultBytesWithBody(http.MethodPut, url, contentType, body, &http.Client{})
}

func defaultBytesWithBody(method, url, contentType string, body []byte, client CandyHTTPClient) (*http.Response, error) {

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add(HeaderContentType, contentType)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil

}

func readBytesWithBody(url, contentType string, body []byte, handler func(string, string, []byte) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url, contentType, body)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}

func readJSONWithBody(url, body string, handler func(str, bd string) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url, body)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}

func readJSON(url string, handler func(str string) (*http.Response, error)) ([]byte, error) {

	res, err := handler(url)
	if err != nil {
		return nil, err
	}

	page, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return page, nil

}
