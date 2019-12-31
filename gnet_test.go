// Copyright 2019 Andy Pan. All rights reserved.
// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gnet

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/panjf2000/gnet/pool/bytebuffer"
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/valyala/bytebufferpool"
)

func TestCodecServe(t *testing.T) {
	// start a server
	// connect 10 clients
	// each client will pipe random data for 1-3 seconds.
	// the writes to the server will be random sizes. 0KB - 1MB.
	// the server will echo back the data.
	// waits for graceful connection closing.
	t.Run("poll", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9991", false, false, 10, false, new(LineBasedFrameCodec))
			})
			t.Run("1-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9992", false, false, 10, false, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("1-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9993", false, false, 10, false, NewFixedLengthFrameCodec(64))
			})
			t.Run("1-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9994", false, false, 10, false, nil)
			})
			t.Run("N-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9995", true, false, 10, false, new(LineBasedFrameCodec))
			})
			t.Run("N-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9996", true, false, 10, false, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("N-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9997", true, false, 10, false, NewFixedLengthFrameCodec(64))
			})
			t.Run("N-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9998", true, false, 10, false, nil)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9991", false, true, 10, false, new(LineBasedFrameCodec))
			})
			t.Run("1-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9992", false, true, 10, false, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("1-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9993", false, true, 10, false, NewFixedLengthFrameCodec(64))
			})
			t.Run("1-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9994", false, true, 10, false, nil)
			})
			t.Run("N-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9995", true, true, 10, false, new(LineBasedFrameCodec))
			})
			t.Run("N-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9996", true, true, 10, false, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("N-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9997", true, true, 10, false, NewFixedLengthFrameCodec(64))
			})
			t.Run("N-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9998", true, true, 10, false, nil)
			})
		})
	})
	t.Run("poll-reuseport", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9991", false, false, 10, true, new(LineBasedFrameCodec))
			})
			t.Run("1-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9992", false, false, 10, true, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("1-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9993", false, false, 10, true, NewFixedLengthFrameCodec(64))
			})
			t.Run("1-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9994", false, false, 10, true, nil)
			})
			t.Run("N-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9995", true, false, 10, true, new(LineBasedFrameCodec))
			})
			t.Run("N-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9996", true, false, 10, true, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("N-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9997", true, false, 10, true, NewFixedLengthFrameCodec(64))
			})
			t.Run("N-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9998", true, false, 10, true, nil)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9991", false, true, 10, true, new(LineBasedFrameCodec))
			})
			t.Run("1-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9992", false, true, 10, true, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("1-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9993", false, true, 10, true, NewFixedLengthFrameCodec(64))
			})
			t.Run("1-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9994", false, true, 10, true, nil)
			})
			t.Run("N-loop-LineBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9995", true, true, 10, true, new(LineBasedFrameCodec))
			})
			t.Run("N-loop-DelimiterBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9996", true, true, 10, true, NewDelimiterBasedFrameCodec('|'))
			})
			t.Run("N-loop-FixedLengthFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9997", true, true, 10, true, NewFixedLengthFrameCodec(64))
			})
			t.Run("N-loop-LengthFieldBasedFrameCodec", func(t *testing.T) {
				testCodecServe("tcp", ":9998", true, true, 10, true, nil)
			})
		})
	})
}

type testCodecServer struct {
	*EventServer
	network      string
	addr         string
	multicore    bool
	async        bool
	nclients     int
	started      int32
	connected    int32
	disconnected int32
	codec        ICodec
	workerPool   *goroutine.Pool
}

func (s *testCodecServer) OnOpened(c Conn) (out []byte, action Action) {
	c.SetContext(c)
	atomic.AddInt32(&s.connected, 1)
	out = []byte("sweetness\r\n")
	if c.LocalAddr() == nil {
		panic("nil local addr")
	}
	if c.RemoteAddr() == nil {
		panic("nil local addr")
	}
	return
}
func (s *testCodecServer) OnClosed(c Conn, err error) (action Action) {
	if c.Context() != c {
		panic("invalid context")
	}

	atomic.AddInt32(&s.disconnected, 1)
	if atomic.LoadInt32(&s.connected) == atomic.LoadInt32(&s.disconnected) &&
		atomic.LoadInt32(&s.disconnected) == int32(s.nclients) {
		action = Shutdown
	}

	return
}
func (s *testCodecServer) React(frame []byte, c Conn) (out []byte, action Action) {
	if s.async {
		if frame != nil {
			data := append([]byte{}, frame...)
			_ = s.workerPool.Submit(func() {
				c.AsyncWrite(data)
			})
		}
		return
	}
	out = frame
	return
}
func (s *testCodecServer) Tick() (delay time.Duration, action Action) {
	if atomic.LoadInt32(&s.started) == 0 {
		for i := 0; i < s.nclients; i++ {
			go func() {
				startCodecClient(s.network, s.addr, s.multicore, s.async, s.codec)
			}()
		}
		atomic.StoreInt32(&s.started, 1)
	}
	delay = time.Second / 5
	return
}

var (
	n            = 0
	fieldLengths = []int{1, 2, 3, 4, 8}
)

func testCodecServe(network, addr string, multicore, async bool, nclients int, reuseport bool, codec ICodec) {
	var err error
	fieldLength := fieldLengths[n]
	if codec == nil {
		encoderConfig := EncoderConfig{
			ByteOrder:                       binary.BigEndian,
			LengthFieldLength:               fieldLength,
			LengthAdjustment:                0,
			LengthIncludesLengthFieldLength: false,
		}
		decoderConfig := DecoderConfig{
			ByteOrder:           binary.BigEndian,
			LengthFieldOffset:   0,
			LengthFieldLength:   fieldLength,
			LengthAdjustment:    0,
			InitialBytesToStrip: fieldLength,
		}
		codec = NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
	}
	n++
	if n > 4 {
		n = 0
	}
	ts := &testCodecServer{network: network, addr: addr, multicore: multicore, async: async, nclients: nclients, codec: codec, workerPool: goroutine.Default()}
	if reuseport {
		err = Serve(ts, network+"://"+addr, WithMulticore(multicore), WithTicker(true),
			WithTCPKeepAlive(time.Minute*5), WithCodec(codec), WithReusePort(true))
	} else {
		err = Serve(ts, network+"://"+addr, WithMulticore(multicore), WithTicker(true),
			WithTCPKeepAlive(time.Minute*5), WithCodec(codec))
	}
	if err != nil {
		panic(err)
	}
}

func startCodecClient(network, addr string, multicore, async bool, codec ICodec) {
	rand.Seed(time.Now().UnixNano())
	c, err := net.Dial(network, addr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	rd := bufio.NewReader(c)
	msg, err := rd.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	if string(msg) != "sweetness\r\n" {
		panic("bad header")
	}
	duration := time.Duration((rand.Float64()*2+1)*float64(time.Second)) / 8
	start := time.Now()
	for time.Since(start) < duration {
		//data := []byte("Hello, World")
		data := make([]byte, 1024)
		rand.Read(data)
		encodedData, _ := codec.Encode(nil, data)
		if _, err := c.Write(encodedData); err != nil {
			panic(err)
		}
		data2 := make([]byte, len(encodedData))
		if _, err := io.ReadFull(rd, data2); err != nil {
			panic(err)
		}
		if string(encodedData) != string(data2) && !async {
			panic(fmt.Sprintf("mismatch %s/multi-core:%t: %d vs %d bytes\n", network, multicore, len(encodedData), len(data2)))
		}
	}
}

func TestServe(t *testing.T) {
	// start a server
	// connect 10 clients
	// each client will pipe random data for 1-3 seconds.
	// the writes to the server will be random sizes. 0KB - 1MB.
	// the server will echo back the data.
	// waits for graceful connection closing.
	t.Run("poll", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				testServe("tcp", ":9991", false, false, false, 10)
			})
			t.Run("N-loop", func(t *testing.T) {
				testServe("tcp", ":9992", false, true, false, 10)
			})
		})
		t.Run("tcp-async", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				testServe("tcp", ":9991", false, false, true, 10)
			})
			t.Run("N-loop", func(t *testing.T) {
				testServe("tcp", ":9992", false, true, true, 10)
			})
		})
		t.Run("udp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				testServe("udp", ":9991", false, false, false, 10)
			})
			t.Run("N-loop", func(t *testing.T) {
				testServe("udp", ":9992", false, true, false, 10)
			})
		})
		//t.Run("unix", func(t *testing.T) {
		//	t.Run("1-loop", func(t *testing.T) {
		//		testServe("unix", "socket9991", false, false, false, 10)
		//	})
		//	t.Run("N-loop", func(t *testing.T) {
		//		testServe("unix", "socket9992", false, true, false, 10)
		//	})
		//})
	})

	t.Run("poll-reuseport", func(t *testing.T) {
		t.Run("tcp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				testServe("tcp", ":9991", true, false, false, 10)
			})
			t.Run("N-loop", func(t *testing.T) {
				testServe("tcp", ":9992", true, true, false, 10)
			})
		})
		t.Run("udp", func(t *testing.T) {
			t.Run("1-loop", func(t *testing.T) {
				testServe("udp", ":9991", true, false, false, 10)
			})
			t.Run("N-loop", func(t *testing.T) {
				testServe("udp", ":9992", true, true, false, 10)
			})
		})
		//t.Run("unix", func(t *testing.T) {
		//	t.Run("1-loop", func(t *testing.T) {
		//		testServe("unix", "socket9991", true, false, false, 10)
		//	})
		//	t.Run("N-loop", func(t *testing.T) {
		//		testServe("unix", "socket9992", true, true, false, 10)
		//	})
		//})
	})
}

type testServer struct {
	*EventServer
	network      string
	addr         string
	multicore    bool
	async        bool
	nclients     int
	started      int32
	connected    int32
	clientActive int32
	disconnected int32
	workerPool   *goroutine.Pool
	bytesList    []*bytebufferpool.ByteBuffer
}

func (s *testServer) OnOpened(c Conn) (out []byte, action Action) {
	c.SetContext(c)
	atomic.AddInt32(&s.connected, 1)
	out = []byte("sweetness\r\n")
	if c.LocalAddr() == nil {
		panic("nil local addr")
	}
	if c.RemoteAddr() == nil {
		panic("nil local addr")
	}
	//fmt.Printf("TCP from remote addr：%s to local addr: %s\n", c.RemoteAddr().String(), c.LocalAddr().String())
	return
}
func (s *testServer) OnClosed(c Conn, err error) (action Action) {
	if c.Context() != c {
		panic("invalid context")
	}

	atomic.AddInt32(&s.disconnected, 1)
	if atomic.LoadInt32(&s.connected) == atomic.LoadInt32(&s.disconnected) &&
		atomic.LoadInt32(&s.disconnected) == int32(s.nclients) {
		action = Shutdown
		for i := range s.bytesList {
			bytebuffer.Put(s.bytesList[i])
		}
		s.workerPool.Release()
	}

	return
}
func (s *testServer) React(frame []byte, c Conn) (out []byte, action Action) {
	if s.async {
		if s.network == "tcp" {
			//bufLen := c.BufferLength()
			buf := bytebuffer.Get()
			_, _ = buf.Write(frame)
			s.bytesList = append(s.bytesList, buf)
			// just for test
			//c.ShiftN(bufLen - 1)
			//
			//c.ShiftN(bufLen)
			//c.ResetBuffer()
			_ = s.workerPool.Submit(
				func() {
					c.AsyncWrite(buf.Bytes())
				})
			return
		}
		if s.network == "udp" {
			_ = s.workerPool.Submit(
				func() {
					c.SendTo(frame)
				})
			return
		}
		return
		//} else if s.multicore {
		//	if s.network == "tcp" {
		//		readSize := 1024 * 1024
		//		n, data := c.ReadN(readSize)
		//		if n == readSize {
		//			out = data
		//		}
		//		return
		//	}
		//	if s.network == "udp" {
		//		out = frame
		//		return
		//	}
		//	return
	} else {
		if s.network == "tcp" {
			out = frame
			return
		}
		//fmt.Printf("UDP from remote addr：%s to local addr: %s\n", c.RemoteAddr().String(), c.LocalAddr().String())
		if s.network == "udp" {
			out = frame
			return
		}
		return
	}
}
func (s *testServer) Tick() (delay time.Duration, action Action) {
	if atomic.LoadInt32(&s.started) == 0 {
		for i := 0; i < s.nclients; i++ {
			atomic.AddInt32(&s.clientActive, 1)
			go func() {
				startClient(s.network, s.addr, s.multicore, s.async)
				atomic.AddInt32(&s.clientActive, -1)
			}()
		}
		atomic.StoreInt32(&s.started, 1)
	}
	if s.network == "udp" && atomic.LoadInt32(&s.clientActive) == 0 {
		action = Shutdown
		return
	}
	delay = time.Second / 5
	return
}

func testServe(network, addr string, reuseport, multicore, async bool, nclients int) {
	var err error
	ts := &testServer{network: network, addr: addr, multicore: multicore, async: async, nclients: nclients, workerPool: goroutine.Default()}
	if network == "unix" {
		_ = os.RemoveAll(addr)
		defer os.RemoveAll(addr)
		err = Serve(ts, network+"://"+addr, WithMulticore(multicore), WithTicker(true), WithTCPKeepAlive(time.Minute*5))
	} else {
		if reuseport {
			err = Serve(ts, network+"://"+addr, WithMulticore(multicore), WithReusePort(true), WithTicker(true), WithTCPKeepAlive(time.Minute*5))
		} else {
			err = Serve(ts, network+"://"+addr, WithMulticore(multicore), WithTicker(true), WithTCPKeepAlive(time.Minute*5))
		}
	}
	if err != nil {
		panic(err)
	}
}

func startClient(network, addr string, multicore, async bool) {
	rand.Seed(time.Now().UnixNano())
	c, err := net.Dial(network, addr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	rd := bufio.NewReader(c)
	if network != "udp" {
		msg, err := rd.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		if string(msg) != "sweetness\r\n" {
			panic("bad header")
		}
	}
	duration := time.Duration((rand.Float64()*2+1)*float64(time.Second)) / 8
	start := time.Now()
	for time.Since(start) < duration {
		//sz := rand.Intn(10) * (1024 * 1024)
		sz := 1024 * 1024
		data := make([]byte, sz)
		if network == "udp" {
			n := 1024
			//if sz < 64 {
			//	n = sz
			//}
			data = data[:n]
		}
		if _, err := rand.Read(data); err != nil {
			panic(err)
		}
		if _, err := c.Write(data); err != nil {
			panic(err)
		}
		data2 := make([]byte, len(data))
		if _, err := io.ReadFull(rd, data2); err != nil {
			panic(err)
		}
		if string(data) != string(data2) && !async {
			panic(fmt.Sprintf("mismatch %s/multi-core:%t: %d vs %d bytes\n", network, multicore, len(data), len(data2)))
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func TestTick(t *testing.T) {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	testTick("tcp", ":9991")
	//}()
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	testTick("unix", "socket1")
	//}()
	//wg.Wait()

	testTick("tcp", ":9991")
}

type testTickServer struct {
	*EventServer
	count int
}

func (t *testTickServer) Tick() (delay time.Duration, action Action) {
	if t.count == 25 {
		action = Shutdown
		return
	}
	t.count++
	delay = time.Millisecond * 10
	return
}
func testTick(network, addr string) {
	events := &testTickServer{}
	start := time.Now()
	opts := Options{Ticker: true}
	must(Serve(events, network+"://"+addr, WithOptions(opts)))
	dur := time.Since(start)
	if dur < 250&time.Millisecond || dur > time.Second {
		panic("bad ticker timing")
	}
}

func TestWakeConn(t *testing.T) {
	testWakeConn("tcp", ":9000")
}

type testWakeConnServer struct {
	*EventServer
	network string
	addr    string
	conn    Conn
	wake    bool
}

func (t *testWakeConnServer) OnOpened(c Conn) (out []byte, action Action) {
	t.conn = c
	return
}

func (t *testWakeConnServer) OnClosed(c Conn, err error) (action Action) {
	action = Shutdown
	return
}

func (t *testWakeConnServer) React(frame []byte, c Conn) (out []byte, action Action) {
	out = []byte("Waking up.")
	return
}
func (t *testWakeConnServer) Tick() (delay time.Duration, action Action) {
	if !t.wake {
		t.wake = true
		delay = time.Millisecond * 100
		go func() {
			conn, err := net.Dial(t.network, t.addr)
			must(err)
			defer conn.Close()
			r := make([]byte, 10)
			_, err = conn.Read(r)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(r))
		}()
		return
	}
	t.conn.Wake()
	delay = time.Millisecond * 100
	return
}

func testWakeConn(network, addr string) {
	svr := &testWakeConnServer{network: network, addr: addr}
	must(Serve(svr, network+"://"+addr, WithTicker(true)))
}

func TestShutdown(t *testing.T) {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	testShutdown("tcp", ":9991")
	//}()
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	testShutdown("unix", "socket1")
	//}()
	//wg.Wait()

	testShutdown("tcp", ":9991")
}

type testShutdownServer struct {
	*EventServer
	network string
	addr    string
	count   int
	clients int64
	N       int
}

func (t *testShutdownServer) OnOpened(c Conn) (out []byte, action Action) {
	atomic.AddInt64(&t.clients, 1)
	return
}
func (t *testShutdownServer) OnClosed(c Conn, err error) (action Action) {
	atomic.AddInt64(&t.clients, -1)
	return
}
func (t *testShutdownServer) Tick() (delay time.Duration, action Action) {
	if t.count == 0 {
		// start clients
		for i := 0; i < t.N; i++ {
			go func() {
				conn, err := net.Dial(t.network, t.addr)
				must(err)
				defer conn.Close()
				_, err = conn.Read([]byte{0})
				if err == nil {
					panic("expected error")
				}
			}()
		}
	} else {
		if int(atomic.LoadInt64(&t.clients)) == t.N {
			action = Shutdown
		}
	}
	t.count++
	delay = time.Second / 20
	return
}
func testShutdown(network, addr string) {
	events := &testShutdownServer{network: network, addr: addr, N: 10}
	must(Serve(events, network+"://"+addr, WithTicker(true)))
	if events.clients != 0 {
		panic("did not call close on all clients")
	}
}

type testBadAddrServer struct {
	*EventServer
}

func (t *testBadAddrServer) OnInitComplete(srv Server) (action Action) {
	return Shutdown
}

func TestBadAddresses(t *testing.T) {
	events := new(testBadAddrServer)
	if err := Serve(events, "tulip://howdy"); err == nil {
		t.Fatalf("expected error")
	}
	if err := Serve(events, "howdy"); err == nil {
		t.Fatalf("expected error")
	}
	if err := Serve(events, "tcp://"); err != nil {
		t.Fatalf("expected nil, got '%v'", err)
	}
}

func TestActionError(t *testing.T) {
	testActionError("tcp", ":9991")
}

type testActionErrorServer struct {
	*EventServer
	network, addr string
	action        bool
}

func (t *testActionErrorServer) OnClosed(c Conn, err error) (action Action) {
	action = Shutdown
	return
}

func (t *testActionErrorServer) React(frame []byte, c Conn) (out []byte, action Action) {
	out = frame
	action = Close
	return
}

func (t *testActionErrorServer) Tick() (delay time.Duration, action Action) {
	if !t.action {
		t.action = true
		delay = time.Millisecond * 100
		go func() {
			conn, err := net.Dial(t.network, t.addr)
			must(err)
			defer conn.Close()
			r := make([]byte, 10)
			rand.Read(r)
			_, _ = conn.Write(r)
			_, err = conn.Read(r)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(r))
		}()
		return
	}
	delay = time.Millisecond * 100
	return
}

func testActionError(network, addr string) {
	events := &testActionErrorServer{network: network, addr: addr}
	must(Serve(events, network+"://"+addr, WithTicker(true)))
}
