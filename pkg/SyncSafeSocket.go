package pkg

import (
	"github.com/gorilla/websocket"
	"sync"
)

// SyncSafeSocket is a wrapper around a websocket.Conn that provides thread-safe
// methods for sending and receiving messages. It uses a mutex to ensure that
// only one goroutine can access the websocket at a time
type SyncSafeSocket struct {
	Socket *websocket.Conn
	Lock   *sync.Mutex
}

func NewSyncSafeSocket(socket *websocket.Conn) *SyncSafeSocket {
	return &SyncSafeSocket{
		Socket: socket,
		Lock:   &sync.Mutex{},
	}
}

func (s *SyncSafeSocket) SendMessage(messageType int, data []byte) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Socket == nil {
		return websocket.ErrCloseSent
	}

	return s.Socket.WriteMessage(messageType, data)
}

func (s *SyncSafeSocket) ReadMessage() (int, []byte, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Socket == nil {
		return 0, nil, websocket.ErrCloseSent
	}

	return s.Socket.ReadMessage()
}

func (s *SyncSafeSocket) Close() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Socket == nil {
		return websocket.ErrCloseSent
	}

	err := s.Socket.Close()
	if err != nil {
		return err
	}
	s.Socket = nil // Prevent further use after closing
	return nil
}

func (s *SyncSafeSocket) WriteJSON(v interface{}) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Socket == nil {
		return websocket.ErrCloseSent
	}

	err := s.Socket.WriteJSON(v)
	if err != nil {
		return err
	}

	return nil
}
