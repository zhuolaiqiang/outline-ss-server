// Copyright 2018 Jigsaw Operations LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"io"

	"github.com/Jigsaw-Code/outline-sdk/transport"
)

type ProxyMetrics struct {
	ClientProxy int64 // 1. 上行入站: 客户端发给代理的数据量 (加密的)
	ProxyTarget int64 // 2. 上行出站: 代理发给目标网站的数据量 (解密后)
	TargetProxy int64 // 3. 下行入站: 目标网站发回给代理的数据量 (明文)
	ProxyClient int64 // 4. 下行出站: 代理发回给客户端的数据量 (加密后)
}

type measuredConn struct {
	transport.StreamConn // 1. 嵌入原始连接 (继承所有方法)
	io.WriterTo          // 显式声明实现 WriterTo 接口
	io.ReaderFrom        // 显式声明实现 ReaderFrom 接口

	readCount  *int64 // 指向外部计数器的指针 (统计读了多少字节)
	writeCount *int64 // 指向外部计数器的指针 (统计写了多少字节)
}

func (c *measuredConn) Read(b []byte) (int, error) {
	n, err := c.StreamConn.Read(b)
	*c.readCount += int64(n)
	return n, err
}

func (c *measuredConn) WriteTo(w io.Writer) (int64, error) {
	n, err := io.Copy(w, c.StreamConn)
	*c.readCount += n
	return n, err
}

func (c *measuredConn) Write(b []byte) (int, error) {
	n, err := c.StreamConn.Write(b)
	*c.writeCount += int64(n)
	return n, err
}

func (c *measuredConn) ReadFrom(r io.Reader) (n int64, err error) {
	if rf, ok := c.StreamConn.(io.ReaderFrom); ok {
		// Prefer ReadFrom if we are calling ReadFrom. Otherwise io.Copy will try WriteTo first.
		n, err = rf.ReadFrom(r)
	} else {
		n, err = io.Copy(c.StreamConn, r)
	}
	*c.writeCount += n
	return n, err
}

func MeasureConn(conn transport.StreamConn, bytesSent, bytesReceived *int64) transport.StreamConn {
	return &measuredConn{StreamConn: conn, writeCount: bytesSent, readCount: bytesReceived}
}
