package spinner

import (
  "github.com/gorilla/websocket"
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
)

// Relay messages amongst connections
type Messenger interface {
  // Open the messenger service. Is open by default.
  Open()
  // Operate the messenger service
  Run()
  // Close and quit the messenger service. Can be re-opened and restarted.
  Close()
  // Container Connection
  ContainerConnect(*dockercntrl.Config) spinresp.ResponseChan
}

type message struct {
  config    *dockercntrl.Config
  response  spinresp.ResponseChan
}

type messenger struct {
  clients       clientPool

  register      clientChan
  unregister    clientChan

  message        chan message

  upgrader      websocket.Upgrader
  quit          chan struct{}
}

// Create a new Messenger interface of messenger struct
func NewMessenger() *messenger {
  return &messenger{
    clients: make(clientPool),
    register: make(clientChan),
    unregister: make(clientChan),
    message: make(chan message),
    upgrader: websocket.Upgrader{
      ReadBufferSize: 1024,
      WriteBufferSize: 1024,
    },
    quit: make(chan struct{}),
  }
}

// Run an infinite loop of message passing
func (m *messenger) Run() {
  running := make(map[*client]int)
  for {
    select {
    case client := <- m.register:
      m.clients[client] = true
      running[client] = 0
    case client := <- m.unregister:
      if _, ok := m.clients[client]; ok {
        delete(running, client)
        delete(m.clients, client)
      }
    case mes := <- m.message:
      // basic round-robbin scheduling
      minimum := 0
      var chosen *client
      for k,v := range running {
        if chosen == nil || v < minimum {
          chosen = k; minimum = v
        }
      }
      if chosen != nil {
        chosen.spinup <- mes
        running[chosen]++
      } else {
        mes.response <- spinresp.Response{
          Code: spinresp.NoCaptainsAvailable,
          Data: "The spinner found no available captians.",
        }
      }
    case <- m.quit:
      return
    }
  }
}

func (m *messenger) ContainerConnect(config *dockercntrl.Config) spinresp.ResponseChan {
  response := make(spinresp.ResponseChan)
  m.message <- message{
    config: config,
    response: response,
  }
  return response
}

// End the message passing loop
func (m *messenger) Close() {close(m.quit)}

// Allow a closed Messager to pass messages again
func (m *messenger) Open() {m.quit = make(chan struct{})}
