package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/zopsmart/gofr/pkg/gofr"
	"github.com/zopsmart/gofr/pkg/gofr/request"
	"github.com/zopsmart/gofr/pkg/gofr/responder"
)

func TestWSHandler(t *testing.T) {
	var (
		conn     *websocket.Conn
		err      error
		upgrader websocket.Upgrader
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
		if err != nil {
			t.Error(err)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	req, _ := request.NewMock("GET", server.URL+"/", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-Websocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "wehkjeh21-sdjk210-wsknb")

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("Hi")); err != nil {
		t.Errorf("could not send message over ws connection %v", err)
	}

	ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), gofr.New())
	ctx.WebSocketConnection = ws

	_ = conn.Close()

	got, err := WSHandler(ctx)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	if got != nil {
		t.Errorf("got: %v", got)
	}
}