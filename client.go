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
  handler       *Handler
  socket        *comms.Socket
  self          *comms.Instance
  spinup        chan interface{}
  responses     map[uuid.UUID]*uuid.UUID
  quit          chan struct{}
}

// Create new Client interface of client struct
func NewClient(h *Handler, socket *comms.Socket) Client {
  return &client{
    handler: h,
    socket: socket,
    spinup: make(chan interface{}),
    responses: make(map[uuid.UUID]*uuid.UUID),
    quit: make(chan struct{}),
    self: nil,
  }
}

// Get messages from the client
func (c *client) Run() {
  defer func(){
    c.handler.Unregister <- c.self
    (*c.socket).Close()
  }()
  read := (*c.socket).Reader()
  write := (*c.socket).Writer()
  for {
    select {
    case response, ok := <- read:
      if !ok {return}
      resp, ok := response.(*spinresp.Response)
      if !ok {return}
      if resp.Id == nil {break}
      if identifier, ok := c.responses[*resp.Id]; ok {
        if resp.Code <= 0 {delete(c.responses, *resp.Id)}
        go func() {
          if !c.handler.Requester.SendMessage(identifier, resp) {
            log.Printf("Failed: %+v\n", resp)
          }
        }()
      }
    case data, ok := <- c.spinup:
      if !ok {break}
      task, ok := data.(*Task)
      if !ok {break}
      if task.Config.Id == nil {
        identifier := uuid.New()
        task.Config.Id = &identifier
      }
      c.responses[*task.Config.Id] = task.From
      write <- task.Config
    }
  }
}

// Register client with messenger and accept read/writes.
func (c *client) Register() {
  var resp spinresp.Response
  (*c.socket).Start(resp)
  c.self = c.handler.Requester.MakeInstance(c.spinup)
  c.handler.Register <- c.self
  go c.Run()
  log.Println("Client registered")

}

// Close the client connection
func (c *client) Quit() {close(c.quit)}
