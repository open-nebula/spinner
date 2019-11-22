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
  responses     map[uuid.UUID]*uuid.UUID  // pair[taskId, whereFrom(userId)]
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
  // read channel in socket connection
  read := (*c.socket).Reader()
  // write channel in socket connection
  write := (*c.socket).Writer()
  for {
    select {
    // response from this captain received
    case response, ok := <- read:
      if !ok {return}
      resp, ok := response.(*spinresp.Response)
      if !ok {return}
      if resp.Id == nil {break}
      if identifier, ok := c.responses[*resp.Id]; ok {
        if resp.Code <= 0 {delete(c.responses, *resp.Id)}  // wrong response code
        go func() {
          if !c.handler.Requester.SendMessage(identifier, resp) {  // send response message back to user
            log.Printf("Failed: %+v\n", resp)
          }
        }()
      }
    // one task scheduled on this captain
    case data, ok := <- c.spinup:
      if !ok {break}
      task, ok := data.(*Task)
      if !ok {break}
      // give task id
      if task.Config.Id == nil {
        identifier := uuid.New()
        task.Config.Id = &identifier
      }
      // store the task info in map before send it to captain
      c.responses[*task.Config.Id] = task.From
      write <- task.Config  // send to the captain (actual sending will be handled in socket package)
    }
  }
}

// Register client with messenger and accept read/writes.
func (c *client) Register() {
  var resp spinresp.Response
  // start read and write routine
  (*c.socket).Start(resp)
  // Bug: should be c.handler.Clients -> returns a captain instance
  c.self = c.handler.Requester.MakeInstance(c.spinup)
  // push new instance (captain) into register queue
  c.handler.Register <- c.self
  go c.Run()
  log.Println("Client registered")

}

// Close the client connection
func (c *client) Quit() {close(c.quit)}
