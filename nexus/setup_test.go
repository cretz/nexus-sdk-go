package nexus

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testTimeout = time.Second * 5
const testService = "Ser/vic e"
const getResultMaxTimeout = time.Millisecond * 300

func setupSerializer(t *testing.T, handler Handler, serializer Serializer) (ctx context.Context, client *HTTPClient, teardown func()) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)

	httpHandler := NewHTTPHandler(HandlerOptions{
		GetResultTimeout: getResultMaxTimeout,
		Handler:          handler,
		Serializer:       serializer,
	})

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	client, err = NewHTTPClient(HTTPClientOptions{
		BaseURL:    fmt.Sprintf("http://%s/", listener.Addr().String()),
		Service:    testService,
		Serializer: serializer,
	})
	require.NoError(t, err)

	go func() {
		// Ignore for test purposes
		_ = http.Serve(listener, httpHandler)
	}()

	return ctx, client, func() {
		cancel()
		listener.Close()
	}
}

func setup(t *testing.T, handler Handler) (ctx context.Context, client *HTTPClient, teardown func()) {
	return setupSerializer(t, handler, nil)
}

func setupForCompletion(t *testing.T, handler CompletionHandler, serializer Serializer) (ctx context.Context, callbackURL string, teardown func()) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)

	httpHandler := NewCompletionHTTPHandler(CompletionHandlerOptions{
		Handler:    handler,
		Serializer: serializer,
	})

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	callbackURL = fmt.Sprintf("http://%s/callback?a=b", listener.Addr().String())

	go func() {
		// Ignore for test purposes
		_ = http.Serve(listener, httpHandler)
	}()

	return ctx, callbackURL, func() {
		cancel()
		listener.Close()
	}
}
