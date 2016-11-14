package utp

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/missinggo"
	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	writeTimeout = 1 * time.Second
	initialLatency = 10 * time.Millisecond
	packetReadTimeout = 2 * time.Second
}

func TestUTPPingPong(t *testing.T) {
	defer goroutineLeakCheck(t)()
	s, err := NewSocket("udp", "localhost:0")
	require.NoError(t, err)
	defer s.Close()
	pingerClosed := make(chan struct{})
	go func() {
		defer close(pingerClosed)
		b, err := Dial(s.Addr().String())
		require.NoError(t, err)
		defer b.Close()
		n, err := b.Write([]byte("ping"))
		require.NoError(t, err)
		require.EqualValues(t, 4, n)
		buf := make([]byte, 4)
		b.Read(buf)
		require.EqualValues(t, "pong", buf)
		log.Printf("got pong")
	}()
	a, err := s.Accept()
	require.NoError(t, err)
	defer a.Close()
	log.Printf("accepted %s", a)
	buf := make([]byte, 42)
	n, err := a.Read(buf)
	require.NoError(t, err)
	require.EqualValues(t, "ping", buf[:n])
	log.Print("got ping")
	n, err = a.Write([]byte("pong"))
	require.NoError(t, err)
	require.Equal(t, 4, n)
	log.Print("waiting for pinger to close")
	<-pingerClosed
}

func goroutineLeakCheck(t testing.TB) func() {
	if !testing.Verbose() {
		return func() {}
	}
	numStart := runtime.NumGoroutine()
	return func() {
		var numNow int
		for range iter.N(1) {
			numNow = runtime.NumGoroutine()
			if numNow == numStart {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		// I'd print stacks, or treat this as fatal, but I think
		// runtime.NumGoroutine is including system routines for which we are
		// not provided the stacks, and are spawned unpredictably.
		t.Logf("have %d goroutines, started with %d", numNow, numStart)
	}
}

func TestDialTimeout(t *testing.T) {
	defer goroutineLeakCheck(t)()
	s, _ := NewSocket("udp", "localhost:0")
	defer s.Close()
	conn, err := DialTimeout(s.Addr().String(), 10*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Fatal("expected timeout")
	}
	t.Log(err)
}

func TestMinMaxHeaderType(t *testing.T) {
	require.Equal(t, stSyn, stMax)
}

func TestUTPRawConn(t *testing.T) {
	l, err := NewSocket("udp", "")
	require.NoError(t, err)
	defer l.Close()
	go func() {
		for {
			_, err := l.Accept()
			if err != nil {
				break
			}
		}
	}()
	// Connect a UTP peer to see if the RawConn will still work.
	log.Print("dialing")
	utpPeer := func() net.Conn {
		s, _ := NewSocket("udp", "")
		defer s.Close()
		ret, err := s.Dial(fmt.Sprintf("localhost:%d", missinggo.AddrPort(l.Addr())))
		require.NoError(t, err)
		return ret
	}()
	log.Print("dial returned")
	if err != nil {
		t.Fatalf("error dialing utp listener: %s", err)
	}
	defer utpPeer.Close()
	peer, err := net.ListenPacket("udp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	msgsReceived := 0
	const N = 5000 // How many messages to send.
	readerStopped := make(chan struct{})
	// The reader goroutine.
	go func() {
		defer close(readerStopped)
		b := make([]byte, 500)
		for i := 0; i < N; i++ {
			n, _, err := l.ReadFrom(b)
			if err != nil {
				t.Fatalf("error reading from raw conn: %s", err)
			}
			msgsReceived++
			var d int
			fmt.Sscan(string(b[:n]), &d)
			if d != i {
				log.Printf("got wrong number: expected %d, got %d", i, d)
			}
		}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", missinggo.AddrPort(l.Addr())))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < N; i++ {
		_, err := peer.WriteTo([]byte(fmt.Sprintf("%d", i)), udpAddr)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Microsecond)
	}
	select {
	case <-readerStopped:
	case <-time.After(time.Second):
		t.Fatal("reader timed out")
	}
	if msgsReceived != N {
		t.Fatalf("messages received: %d", msgsReceived)
	}
}

func TestConnReadDeadline(t *testing.T) {
	t.Parallel()
	ls, _ := NewSocket("udp", "localhost:0")
	ds, _ := NewSocket("udp", "localhost:0")
	dcReadErr := make(chan error)
	go func() {
		c, _ := ds.Dial(ls.Addr().String())
		defer c.Close()
		_, err := c.Read(nil)
		dcReadErr <- err
	}()
	c, _ := ls.Accept()
	dl := time.Now().Add(time.Millisecond)
	c.SetReadDeadline(dl)
	_, err := c.Read(nil)
	require.Equal(t, errTimeout, err)
	// The deadline has passed.
	if !time.Now().After(dl) {
		t.FailNow()
	}
	// Returns timeout on subsequent read.
	_, err = c.Read(nil)
	require.Equal(t, errTimeout, err)
	// Disable the deadline.
	c.SetReadDeadline(time.Time{})
	readReturned := make(chan struct{})
	go func() {
		c.Read(nil)
		close(readReturned)
	}()
	select {
	case <-readReturned:
		// Read returned but shouldn't have.
		t.FailNow()
	case <-time.After(time.Millisecond):
	}
	c.Close()
	if err := <-dcReadErr; err != io.EOF {
		t.Fatalf("dial conn read returned %s", err)
	}
	select {
	case <-readReturned:
	case <-time.After(time.Millisecond):
		t.Fatal("read should return after Conn is closed")
	}
}

func connectSelfLots(n int, t testing.TB) {
	defer goroutineLeakCheck(t)()
	s, err := NewSocket("udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for range iter.N(n) {
			c, err := s.Accept()
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
		}
	}()
	dialErr := make(chan error)
	connCh := make(chan net.Conn)
	dialSema := make(chan struct{}, backlog)
	for range iter.N(n) {
		go func() {
			dialSema <- struct{}{}
			c, err := s.Dial(s.Addr().String())
			<-dialSema
			if err != nil {
				dialErr <- err
				return
			}
			connCh <- c
		}()
	}
	conns := make([]net.Conn, 0, n)
	for range iter.N(n) {
		select {
		case c := <-connCh:
			conns = append(conns, c)
		case err := <-dialErr:
			t.Fatal(err)
		}
	}
	for _, c := range conns {
		if c != nil {
			c.Close()
		}
	}
	s.mu.Lock()
	for len(s.conns) != 0 {
		// log.Print(len(s.conns))
		s.event.Wait()
	}
	s.mu.Unlock()
	s.Close()
}

// Connect to ourself heaps.
func TestConnectSelf(t *testing.T) {
	t.Parallel()
	// A rough guess says that at worst, I can only have 0x10000/3 connections
	// to the same socket, due to fragmentation in the assigned connection
	// IDs.
	connectSelfLots(0x1000, t)
}

func BenchmarkConnectSelf(b *testing.B) {
	for range iter.N(b.N) {
		connectSelfLots(2, b)
	}
}

func BenchmarkNewCloseSocket(b *testing.B) {
	for range iter.N(b.N) {
		s, err := NewSocket("udp", "localhost:0")
		if err != nil {
			b.Fatal(err)
		}
		err = s.Close()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRejectDialBacklogFilled(t *testing.T) {
	s, err := NewSocket("udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	errChan := make(chan error, 1)
	dial := func() {
		_, err := s.Dial(s.Addr().String())
		if err != nil {
			errChan <- err
		}
	}
	// Fill the backlog.
	for range iter.N(backlog + 1) {
		go dial()
	}
	s.mu.Lock()
	for len(s.backlog) < backlog {
		s.event.Wait()
	}
	s.mu.Unlock()
	select {
	case <-errChan:
		t.FailNow()
	default:
	}
	// One more connection should cause a dial attempt to get reset.
	go dial()
	err = <-errChan
	if err.Error() != "peer reset" {
		t.FailNow()
	}
	s.Close()
}

// Make sure that we can reset AfterFunc timers, so we don't have to create
// brand new ones everytime they fire. Specifically for the Conn resend timer.
func TestResetAfterFuncTimer(t *testing.T) {
	t.Parallel()
	fired := make(chan struct{})
	timer := time.AfterFunc(time.Millisecond, func() {
		fired <- struct{}{}
	})
	<-fired
	if timer.Reset(time.Millisecond) {
		// The timer should have expired
		t.FailNow()
	}
	<-fired
}

func connPair() (initer, accepted net.Conn) {
	s, err := NewSocket("udp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer s.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		initer, err = Dial(s.Addr().String())
		if err != nil {
			panic(err)
		}
	}()
	accepted, err = s.Accept()
	if err != nil {
		panic(err)
	}
	wg.Wait()
	return
}

// Check that peer sending FIN doesn't cause unread data to be dropped in a
// receiver.
func TestReadFinishedConn(t *testing.T) {
	a, b := connPair()
	defer a.Close()
	defer b.Close()
	mu.Lock()
	originalAPDC := artificialPacketDropChance
	artificialPacketDropChance = 1
	mu.Unlock()
	n, err := a.Write([]byte("hello"))
	require.Equal(t, 5, n)
	require.NoError(t, err)
	n, err = a.Write([]byte("world"))
	require.Equal(t, 5, n)
	require.NoError(t, err)
	mu.Lock()
	artificialPacketDropChance = originalAPDC
	mu.Unlock()
	a.Close()
	all, err := ioutil.ReadAll(b)
	require.NoError(t, err)
	require.EqualValues(t, "helloworld", all)
}

func TestCloseDetachesQuickly(t *testing.T) {
	t.Parallel()
	s, _ := NewSocket("udp", "localhost:0")
	defer s.Close()
	go func() {
		a, _ := s.Dial(s.Addr().String())
		log.Print("close a")
		a.Close()
		log.Print("closed a")
	}()
	b, _ := s.Accept()
	b.Close()
	s.mu.Lock()
	for len(s.conns) != 0 {
		log.Print(len(s.conns))
		s.event.Wait()
	}
	s.mu.Unlock()
}

// Check that closing, and resulting detach of a Conn doesn't close the parent
// Socket. We Accept, then close the connection and ensure it's detached. Then
// Accept again to check the Socket is still functional and unclosed.
func TestConnCloseUnclosedSocket(t *testing.T) {
	t.Parallel()
	s, err := NewSocket("udp", "localhost:0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, s.Close())
	}()
	// Prevents the dialing goroutine from closing its end of the Conn before
	// we can check that it has been registered in the listener.
	dialerSync := make(chan struct{})

	go func() {
		for range iter.N(2) {
			c, err := Dial(s.Addr().String())
			require.NoError(t, err)
			<-dialerSync
			err = c.Close()
			require.NoError(t, err)
		}
	}()
	for range iter.N(2) {
		a, err := s.Accept()
		require.NoError(t, err)
		// We do this in a closure because we need to unlock Server.mu if the
		// test failure exception is thrown. "Do as we say, not as we do" -Go
		// team.
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			require.Len(t, s.conns, 1)
		}()
		dialerSync <- struct{}{}
		require.NoError(t, a.Close())
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			for len(s.conns) != 0 {
				s.event.Wait()
			}
		}()
	}
}

func TestAcceptGone(t *testing.T) {
	s, _ := NewSocket("udp", "localhost:0")
	_, err := DialTimeout(s.Addr().String(), time.Millisecond)
	require.Error(t, err)
	c, _ := s.Accept()
	c.SetReadDeadline(time.Now().Add(time.Millisecond))
	c.Read(nil)
	// select {}
}

func TestPacketReadTimeout(t *testing.T) {
	t.Parallel()
	a, b := connPair()
	_, err := a.Read(nil)
	require.Contains(t, err.Error(), "timeout")
	t.Log(err)
	t.Log(a.Close())
	t.Log(b.Close())
}
