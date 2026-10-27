package controller

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/calebportodev/golang/hare-st/src/internal/app"
	"github.com/calebportodev/golang/hare-st/src/internal/storage"
)

func TestHandleConnection(t *testing.T) {
	app.Storage = &storage.StorageMap{}

	serverConn, clientConn := net.Pipe()
	c := &Controller{Conn: serverConn}

	done := make(chan struct{})
	go func() {
		c.HandleConnection()
		close(done)
	}()

	client := bufio.NewReader(clientConn)

	write := func(s string) {
		clientConn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		if _, err := clientConn.Write([]byte(s)); err != nil {
			t.Fatalf("failed to write: %v", err)
		}
	}
	readLine := func() string {
		clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		line, err := client.ReadString('\n')
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		return line
	}

	t.Run("valid put", func(t *testing.T) {
		write("put chave int 42\n")
		got := readLine()
		want := "status success\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("valid get", func(t *testing.T) {
		write("get chave\n")
		got := readLine()
		want := "data: int 42\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("valid clean", func(t *testing.T) {
		write("clean chave\n")
		got := readLine()
		want := "status success\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		write("invalidverb chave\n")
		got := readLine()
		if got == "" {
			t.Fatal("expected an error response, got empty string")
		}
	})

	clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return after connection close")
	}
}

func TestResponse(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	c := &Controller{Conn: serverConn}

	go c.Response("hello")

	clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 16)
	n, err := clientConn.Read(buf)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}
	got := string(buf[:n])
	want := "hello\n"
	if got != want {
		t.Errorf("Response() wrote %q, want %q", got, want)
	}
	clientConn.Close()
	serverConn.Close()
}
