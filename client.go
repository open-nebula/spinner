package spinner

import (
  "github.com/open-nebula/spinner/spinresp"
  "github.com/google/uuid"
  "github.com/open-nebula/comms"
  "log"
)

// Single client (Captain) connection
type Client interface {
  Run()
  // Enter client into messenger system
  Register()
  // Quit the client
  Quit()
}

type client struct {
  messenger     *messenger
  socket        *comms.Socket
  spinup        chan message
  responses     map[uuid.UUID]spinresp.ResponseChan
  quit          chan struct{}
}

type clientPool   map[*client]bool
type clientChan   chan *client

// Create new Client interface of client struct
func NewClient(m *messenger, socket *comms.Socket) Client {
  return &client{
    messenger: m,
    socket: socket,
    spinup: make(chan message),
    responses: make(map[uuid.UUID]spinresp.ResponseChan),
    quit: make(chan struct{}),
  }
}

// Get messages from the client
func (c *client) Run() {
  defer func(){
    c.messenger.unregister <- c
    (*c.socket).Close()
  }()
  read := (*c.socket).Reader()
  write := (*c.socket).Writer()
  for {
    select {
    case response, ok := <- read:
      if !ok {return}
      resp, ok := response.(spinresp.Response)
      if !ok {return}
      if resp.Id != nil {return}
      if respchan, ok := c.responses[*resp.Id]; ok {
        respchan <- resp
        if resp.Code < 0 {delete(c.responses, *resp.Id); return}
      } else {return}
    case message, ok := <- c.spinup:
      if !ok {return}
      if message.config.Id == nil {
        identifier := uuid.New()
        message.config.Id = &identifier
      }
      write <- message.config
      c.responses[*message.config.Id] = message.response
    }
  }
}

// Register client with messenger and accept read/writes.
func (c *client) Register() {
  c.messenger.register <- c
  go c.Run()
  log.Println("Client registered")
}

// Close the client connection
func (c *client) Quit() {close(c.quit)}
