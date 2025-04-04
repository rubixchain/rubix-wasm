package context

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type WasmContext struct {
	baseCtx context.Context
	externalSocketConn *websocket.Conn
}

func (c *WasmContext) WithExternalSocketConn(conn *websocket.Conn) *WasmContext {
	c.externalSocketConn = conn
	return c
}

func (c WasmContext) ExternalSocketConn() *websocket.Conn {
	if c.externalSocketConn == nil {
		return nil
	}
	return c.externalSocketConn
}

func (c WasmContext) Deadline() (deadline time.Time, ok bool) {
	return c.baseCtx.Deadline()
}

func (c WasmContext) Done() <-chan struct{} {
	return c.baseCtx.Done()
}

func (c WasmContext) Err() error {
	return c.baseCtx.Err()
}

func (c WasmContext) Value(key interface{}) interface{} {
	return c.baseCtx.Value(key)
}

var _ context.Context = WasmContext{}

func NewWasmContext() *WasmContext {
	return &WasmContext{
		baseCtx: context.Background(),
		externalSocketConn: nil,
	}
}
