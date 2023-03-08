package proto

import "time"

type Message struct {
	ClientSend time.Time
	ServerRecv time.Time
	ServerSend time.Time
	ClientRecv time.Time
}
