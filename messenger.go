package spinner

import (
  "github.com/gorilla/websocket"
)

type Messenger struct {
  Clients       ClientPool

  Register      ClientChan
  Unregister    ClientChan

  Upgrader      websocket.Upgrader
  quit          chan struct{}
}

func NewMessenger() *Messenger {
  messenger := &Messenger{
    Clients: make(ClientPool),
    Register: make(ClientChan),
    Unregister: make(ClientChan),
    Upgrader: websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: 1024,
    },
    quit: make(chan struct{}),
  }
  messenger.Run()
  return messenger
}

func (m *Messenger) Run() {
  go func() {
    for {
      select {
      case client := <- m.Register:
        m.Clients[client] = true
      case client := <- m.Unregister:
        if _, ok := m.Clients[client]; ok {
          delete(m.Clients, client)
        }
      case <- m.quit:
        return
      }
    }
  }()
}

func (m *Messenger) Close() {close(m.quit)}
func (m *Messenger) Open() {m.quit = make(chan struct{})}
