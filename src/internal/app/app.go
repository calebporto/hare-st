package app

import (
	"fmt"

	"github.com/calebportodev/golang/hare-st/src/internal/protocol"
	"github.com/calebportodev/golang/hare-st/src/internal/storage"
)

var Storage = &storage.StorageMap{}

type RequestHandler struct {
	StorageMap *storage.StorageMap
}

func (h *RequestHandler) HandleRequest(request protocol.ParsedRawRequest) (string, error) {
	switch request.Verb {
	case "get":
		// Handle GET request
		data, ok := h.StorageMap.Get(request.Key)
		if !ok {
			return "", fmt.Errorf("Key not found")
		}
		return data.Data, nil
	case "clean":
		// Handle CLEAN request
		h.StorageMap.Clean(request.Key)
		return "clean", nil
	case "put":
		// Handle PUT request
		h.StorageMap.Put(request.Key, storage.StorageValue{DataType: request.DataType, Data: request.Data})
		return "put", nil
	}
	return "", fmt.Errorf("Unknown verb: %s", request.Verb)
}

