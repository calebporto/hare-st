package server

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	go Start("127.0.0.1:18379")

	var conn net.Conn
	var err error
	for i := 0; i < 20; i++ {
		conn, err = net.DialTimeout("tcp", "127.0.0.1:18379", 100*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write([]byte("put chave string ola\n")); err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	want := "status success\n"
	if line != want {
		t.Errorf("got %q, want %q", line, want)
	}
}

func TestStartInvalidAddress(t *testing.T) {
	// An address with an invalid port should fail to listen.
	Start("invalid-address")
}
