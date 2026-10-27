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
	parser := protocol.Parser{}

	switch request.Verb {
	case "get":
		// Handle GET request
		data, ok := h.StorageMap.Get(request.Key)
		return parser.ParseGetResponse(data, ok)
	case "clean":
		// Handle CLEAN request
		h.StorageMap.Clean(request.Key)
		return parser.ParseCleanResponse(request)
	case "put":
		// Handle PUT request
		h.StorageMap.Put(request.Key, storage.StorageValue{DataType: request.DataType, Data: request.Data})
		return parser.ParsePutResponse(request)
	}
	return "error", fmt.Errorf("Unknown verb: %s", request.Verb)
}

