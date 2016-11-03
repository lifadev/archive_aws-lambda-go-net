//
// Copyright 2016 Alsanium, SAS. or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package net

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eawsy/aws-lambda-go/service/lambda/runtime"
)

var (
	laddr, raddr  *net.TCPAddr
	conn          *lambdaConn
	inc, outc     chan struct{}
	inbuf, outbuf bytes.Buffer
	lsnr          net.Listener
)

type lambdaConn struct{}

func (c *lambdaConn) Read(b []byte) (int, error) {
	return inbuf.Read(b)
}

func (c *lambdaConn) Write(b []byte) (int, error) {
	return outbuf.Write(b)
}

func (c *lambdaConn) Close() error {
	outc <- struct{}{}
	return nil
}

func (c *lambdaConn) LocalAddr() net.Addr {
	return laddr
}

func (c *lambdaConn) RemoteAddr() net.Addr {
	return raddr
}

func (c *lambdaConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *lambdaConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *lambdaConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type lambdaListener struct{}

func (l *lambdaListener) Accept() (net.Conn, error) {
	<-inc
	return conn, nil
}

func (l *lambdaListener) Close() error {
	return nil
}

func (l *lambdaListener) Addr() net.Addr {
	return laddr
}

type gtwInput struct {
	HTTPMethod            string                 `json:"httpMethod"`
	Path                  string                 `json:"path"`
	QueryStringParameters map[string]string      `json:"queryStringParameters"`
	Headers               map[string]string      `json:"headers"`
	Body                  string                 `json:"body"`
	RequestContext        map[string]interface{} `json:"requestContext"`
	Resource              interface{}            `json:"resource"`
	PathParameters        interface{}            `json:"pathParameters"`
	StageVariables        map[string]interface{} `json:"stageVariables"`
}

type gtwOutput struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// Listener returns a network listener corresponding to Lambda function that is executing.
func Listener() net.Listener {
	return lsnr
}

func init() {
	laddr = &net.TCPAddr{IP: net.IPv4zero}
	raddr = &net.TCPAddr{IP: net.IPv4zero}
	conn = &lambdaConn{}
	inc = make(chan struct{})
	outc = make(chan struct{})
	lsnr = &lambdaListener{}

	runtime.HandleFunc(func(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
		inbuf.Reset()
		outbuf.Reset()

		var ingtw gtwInput

		err := json.Unmarshal(evt, &ingtw)
		if err != nil {
			return nil, err
		}

		raddr.IP = net.ParseIP(ingtw.RequestContext["identity"].(map[string]interface{})["sourceIp"].(string))

		u, _ := url.Parse(ingtw.Path)
		q := u.Query()
		for k, v := range ingtw.QueryStringParameters {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()

		req, _ := http.NewRequest(ingtw.HTTPMethod, u.String(), strings.NewReader(ingtw.Body))

		for k, v := range ingtw.Headers {
			req.Header.Set(k, v)
		}
		inctx, _ := json.Marshal(ingtw.RequestContext)
		invars, _ := json.Marshal(ingtw.StageVariables)

		req.Header.Set("X-Amz-ApiGtw-Ctx", base64.StdEncoding.EncodeToString(inctx))
		req.Header.Set("X-Amz-ApiGtw-Vars", base64.StdEncoding.EncodeToString(invars))

		req.Write(&inbuf)

		inc <- struct{}{}
		<-outc

		res, _ := http.ReadResponse(bufio.NewReader(&outbuf), req)

		outgtw := &gtwOutput{}
		outgtw.StatusCode = res.StatusCode
		outgtw.Headers = make(map[string]string)
		for k := range res.Header {
			outgtw.Headers[k] = res.Header.Get(k)
		}
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		outgtw.Body = string(body)

		return outgtw, nil
	})
}
