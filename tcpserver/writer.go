package tcpserver

import (
	"github.com/iotames/qrbridge/tcpserver/model"
)

// webSocketWriter 用于 WebSocket 连接，自动封帧并通过 channel 发送
type webSocketWriter struct {
	user *User
}

func (w *webSocketWriter) Write(p []byte) (n int, err error) {
	framed := model.WebSocketPack(p)
	w.user.ReceiveDataToSend(framed)
	return len(p), nil
}

// rawTCPWriter 用于普通 TCP 连接，直接通过 channel 发送原始数据
type rawTCPWriter struct {
	user *User
}

func (w *rawTCPWriter) Write(p []byte) (n int, err error) {
	w.user.ReceiveDataToSend(p)
	return len(p), nil
}
