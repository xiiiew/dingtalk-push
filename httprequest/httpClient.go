package httprequest

import (
	"bytes"
	"crypto/tls"
	loggerHttpRequest "dingtalk-push/utils/log/httpRequest"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func HttpGet(urlPath string, args map[string]string) ([]byte, error) {
	respBody := make(chan []byte)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go func() {
		select {
		case body := <-respBody:
			b := []byte{}
			if body != nil {
				b = body
			}

			loggerHttpRequest.InfoWithFields(
				urlPath,
				loggerHttpRequest.Fields{
					"method": "GET",
					"args":   args,
					"resp":   string(b),
				},
			)
		}
		wg.Done()
	}()

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}}
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		respBody <- nil
		return nil, err
	}
	params := url.Values{}
	for k, v := range args {
		params.Set(k, v)
	}
	req.URL.RawQuery = params.Encode()
	resp, err := client.Do(req)
	if err != nil {
		respBody <- nil
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respBody <- nil
		return nil, err
	}
	respBody <- body
	return body, err
}

func HttpPost(urlPath string, args map[string]string) ([]byte, error) {
	respBody := make(chan []byte)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go func() {
		select {
		case body := <-respBody:
			b := []byte{}
			if body != nil {
				b = body
			}

			loggerHttpRequest.InfoWithFields(
				urlPath,
				loggerHttpRequest.Fields{
					"method": "GET",
					"args":   args,
					"resp":   string(b),
				},
			)
		}
		wg.Done()
	}()

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, DisableKeepAlives: true}}
	params := url.Values{}
	for k, v := range args {
		params.Set(k, v)
	}
	req, err := http.NewRequest("POST", urlPath, strings.NewReader(params.Encode()))
	if err != nil {
		respBody <- nil
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		respBody <- nil
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respBody <- nil
		return nil, err
	}
	respBody <- body
	return body, err
}

func HTTPPostJsonWithClient(urlPath string, params interface{}, client *http.Client) (_ []byte, err error) {
	respBody := make(chan []byte)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	go func() {
		select {
		case body := <-respBody:
			b := []byte{}
			if body != nil {
				b = body
			}

			loggerHttpRequest.InfoWithFields(
				urlPath,
				loggerHttpRequest.Fields{
					"method": "POST",
					"args":   params,
					"resp":   string(b),
				},
			)
		}
		wg.Done()
	}()

	bytesParams, err := json.Marshal(params)
	if err != nil {
		return
	}
	reader := bytes.NewReader(bytesParams)

	req, err := http.NewRequest("POST", urlPath, reader)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp == nil {
		return nil, errors.New("response is empty")
	}

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respBody <- bytesBody
	return bytesBody, nil
}
