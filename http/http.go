package toolshttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GoHttpClient struct {
	client  *http.Client
	method  string
	url     string
	headers []map[string]string
	body    io.Reader
}

type HttpClient interface {
	SetClient(client *http.Client) *GoHttpClient
	SetMethod(method string) *GoHttpClient
	SetUrl(url string) *GoHttpClient
	SetHeader(key, val string) *GoHttpClient
	SetBody(body io.Reader) *GoHttpClient
	Do() (int, *bytes.Buffer, error)
}

func NewHttpClient() *GoHttpClient {
	return &GoHttpClient{}
}

func (ghc *GoHttpClient) SetClient(client *http.Client) *GoHttpClient {
	ghc.client = client
	return ghc
}

func (ghc *GoHttpClient) SetMethod(method string) *GoHttpClient {
	ghc.method = method
	return ghc
}

func (ghc *GoHttpClient) SetUrl(url string) *GoHttpClient {
	ghc.url = url
	return ghc
}

func (ghc *GoHttpClient) SetHeader(key, val string) *GoHttpClient {
	ghc.headers = append(ghc.headers, map[string]string{key: val})
	return ghc
}

func (ghc *GoHttpClient) SetBody(body io.Reader) *GoHttpClient {
	ghc.body = body
	return ghc
}

func (ghc *GoHttpClient) Do() (int, *bytes.Buffer, error) {
	if ghc.client == nil {
		ghc.client = &http.Client{}
	}

	if len(ghc.url) == 0 {
		return 0, nil, errors.New(fmt.Sprintf("url error %s", ghc.url))
	}

	req, err := http.NewRequest(ghc.method, ghc.url, ghc.body)
	if err != nil {
		return 0, nil, err
	}

	if cap(ghc.headers) > 0 {
		for _, header := range ghc.headers {
			for key, val := range header {
				req.Header.Set(key, val)
			}
		}
	}
	resp, err := ghc.client.Do(req)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("request error %v", err))
	}

	defer resp.Body.Close()

	var buffer [512]byte
	response := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		response.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return 0, nil, errors.New(fmt.Sprintf("read response body error %v", err))
		}
	}
	return resp.StatusCode, response, nil
}
