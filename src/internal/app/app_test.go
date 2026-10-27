package app

import (
	"testing"

	"github.com/calebportodev/golang/hare-st/src/internal/protocol"
	"github.com/calebportodev/golang/hare-st/src/internal/storage"
)

func newHandler() *RequestHandler {
	return &RequestHandler{StorageMap: &storage.StorageMap{}}
}

func TestHandleRequestPut(t *testing.T) {
	h := newHandler()

	response, err := h.HandleRequest(protocol.ParsedRawRequest{
		Verb:     "put",
		Key:      "chave",
		DataType: "int",
		Data:     "42",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response != "status success" {
		t.Errorf("HandleRequest(put) = %q, want %q", response, "status success")
	}

	value, ok := h.StorageMap.Get("chave")
	if !ok {
		t.Fatal("expected key to be stored")
	}
	if value.DataType != "int" || value.Data != "42" {
		t.Errorf("stored value = %+v, want DataType=int Data=42", value)
	}
}

func TestHandleRequestGet(t *testing.T) {
	h := newHandler()
	h.StorageMap.Put("chave", storage.StorageValue{DataType: "int", Data: "42"})

	t.Run("existing key", func(t *testing.T) {
		response, err := h.HandleRequest(protocol.ParsedRawRequest{Verb: "get", Key: "chave"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := "data: int 42"
		if response != want {
			t.Errorf("HandleRequest(get) = %q, want %q", response, want)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		response, err := h.HandleRequest(protocol.ParsedRawRequest{Verb: "get", Key: "inexistente"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		want := "data: nil nil"
		if response != want {
			t.Errorf("HandleRequest(get) = %q, want %q", response, want)
		}
	})
}

func TestHandleRequestClean(t *testing.T) {
	h := newHandler()
	h.StorageMap.Put("chave", storage.StorageValue{DataType: "int", Data: "42"})

	response, err := h.HandleRequest(protocol.ParsedRawRequest{Verb: "clean", Key: "chave"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response != "status success" {
		t.Errorf("HandleRequest(clean) = %q, want %q", response, "status success")
	}

	_, ok := h.StorageMap.Get("chave")
	if ok {
		t.Fatal("expected key to be removed")
	}
}

func TestHandleRequestUnknownVerb(t *testing.T) {
	h := newHandler()

	_, err := h.HandleRequest(protocol.ParsedRawRequest{Verb: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown verb, got nil")
	}
}
