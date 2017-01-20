//
// Copyright 2017 Alsanium, SAS. or its affiliates. All rights reserved.
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
	"io"
	"net"
	"time"
)

// LambdaListener is an AWS Lambda network listener.
type LambdaListener struct {
	conn chan net.Conn
}

// Accept implements the Accept method in the Go net.Listener interface.
// It waits for the next call of the AWS Lambda function and returns a generic
// net.Conn.
func (l *LambdaListener) Accept() (net.Conn, error) {
	return <-l.conn, nil
}

// Close mocks the Close method in the Go net.Listener interface.
func (l *LambdaListener) Close() error {
	return nil
}

// Addr mocks the Addr method in the Go net.Listener interface.
func (l *LambdaListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4zero}
}

// Handle is the glue betwwen Go net.Listener world and AWS Lambda world.
// It provides a mechanism similar to standard Go TCP, UDP, etc. listeners for
// waiting for the next connection aka the next AWS Lambda call.
func (l *LambdaListener) Handle(addr net.Addr, req io.Reader, res io.Writer) {
	conn := &LambdaConn{addr: addr, req: req, res: res, done: make(chan struct{})}
	l.conn <- conn
	<-conn.done
}

// ListenLambda creates, initializes and returns an AWS Lambda listener.
func ListenLambda() *LambdaListener {
	return &LambdaListener{make(chan net.Conn)}
}

// LambdaConn is an implementation of the Go net.Conn interface for AWS Lambda
// mocked network connections.
type LambdaConn struct {
	addr net.Addr
	req  io.Reader
	res  io.Writer
	done chan struct{}
}

// Read implements the Go net.Conn Read method.
// It reads the request prepared by the AWS Lambda handler.
func (c *LambdaConn) Read(b []byte) (n int, err error) {
	return c.req.Read(b)
}

// Write implements the Go net.Conn Write method.
// It writes the response which will be consumed by the AWS Lambda handler.
func (c *LambdaConn) Write(b []byte) (n int, err error) {
	return c.res.Write(b)
}

// Close closes the connection by signaling the end of the request processing to
// the AWS Lambda handler.
func (c *LambdaConn) Close() error {
	close(c.done)
	return nil
}

// LocalAddr returns the mocked local network address.
// The Addr returned is shared by all invocations of LocalAddr, so do not modify
// it.
func (c *LambdaConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4zero}
}

// RemoteAddr returns the remote network address given by the AWS Lambda
// handler.
// The Addr returned is shared by all invocations of RemoteAddr, so do not
// modify it.
func (c *LambdaConn) RemoteAddr() net.Addr {
	return c.addr
}

// SetDeadline mocks the Go net.Conn SetDeadline method.
func (c *LambdaConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline mocks the Go net.Conn SetDeadline method.
func (c *LambdaConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline mocks the Go net.Conn SetDeadline method.
func (c *LambdaConn) SetWriteDeadline(t time.Time) error {
	return nil
}
