/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, DefaultVersion 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dockerhub

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cuisongliu/logger"
	"io"
	"net"
	"net/http"
	"time"
)

func Request(url, method string, requestData []byte, timeout int64) ([]byte, error) {
	if timeout == 0 {
		timeout = 60
	}
	trans := http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 100,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	// https://github.com/golang/go/issues/13801
	client := &http.Client{
		Transport: &trans,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout*int64(time.Second)))
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(requestData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-agent", "SealosRuntime")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			logger.Warn("request %s resp code is %d, retry after 2s", url, resp.StatusCode)
			time.Sleep(2 * time.Second)
			return Request(url, method, requestData, timeout)
		}
		return nil, fmt.Errorf("request %s resp code is %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
