package spinner

import (
  "github.com/gorilla/websocket"
)

// Relay messages amongst connections
type Messenger interface {
  // Open the messenger service. Is open by default.
  Open()
  // Operate the messenger service
  Run()
  // Close and quit the messenger service. Can be re-opened and restarted.
  Close()
}

type messenger struct {
  clients       ClientPool

  register      ClientChan
  unregister    ClientChan

  upgrader      websocket.Upgrader
  quit          chan struct{}
}

func NewMessenger() Messenger {
  return &messenger{
    clients: make(ClientPool),
    register: make(ClientChan),
    unregister: make(ClientChan),
    upgrader: websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: 1024,
    },
    quit: make(chan struct{}),
  }
}

func (m *messenger) Run() {
  for {
    select {
    case client := <- m.register:
      m.clients[client] = true
    case client := <- m.unregister:
      if _, ok := m.clients[client]; ok {
        delete(m.clients, client)
      }
    case <- m.quit:
      return
    }
  }
}

func (m *messenger) Close() {close(m.quit)}
func (m *messenger) Open() {m.quit = make(chan struct{})}
