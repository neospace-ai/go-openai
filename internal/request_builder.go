package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type RequestBuilder interface {
	Build(ctx context.Context, method, url string, body any, header http.Header) (*http.Request, error)
}

type HTTPRequestBuilder struct {
	marshaller Marshaller
}

func NewRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		marshaller: &JSONMarshaller{},
	}
}

func (b *HTTPRequestBuilder) Build(
	ctx context.Context,
	method string,
	url string,
	body any,
	header http.Header,
) (req *http.Request, err error) {
	var bodyReader io.Reader
	if body != nil {
		if v, ok := body.(io.Reader); ok {
			fmt.Printf("NEOSPACE CUSTOM OPENAI IO.READER: %v", v)
			bodyReader = v
		}

		var reqBytes []byte
		reqBytes, err = b.marshaller.Marshal(body)
		if err != nil {
			return
		}
		bodyReader = bytes.NewBuffer(reqBytes)
	}

	req, err = http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return
	}
	if header != nil {
		req.Header = header
	}
	return
}
