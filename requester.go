package spinner

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "log"
)

// Single container requestor interface
type Requester interface {
  // Accept reading from the client
  Run()
  // Enter client into messenger system
  Register()
  // Quit the client
  Quit()
}

type requester struct {
  handler       *Handler
  socket        *comms.Socket
  responses     spinresp.ResponseChan
  quit          chan struct{}
  self          *comms.Instance  // user instance
}

// Create new Client interface of client struct
func NewRequester(h *Handler, socket *comms.Socket) Requester {
  return &requester{
    handler: h,
    socket: socket,
    quit: make(chan struct{}),
    self: nil,
  }
}

// Get messages from the requester
func (r *requester) Run() {
  defer func(){
    (*r.socket).Close()
  }()
  read := (*r.socket).Reader()
  log.Println("Requester running")
  for {
    select {
    case config, ok := <- read:
      if !ok {return}
      dockerconfig, ok := config.(*dockercntrl.Config)
      if !ok {return}
      log.Println(dockerconfig)
      log.Printf("%T\n", dockerconfig)
      go func() {
        log.Printf("Sending: %+v\n", dockerconfig)
        success := r.handler.SendTask(r.self, dockerconfig)
        if !success {
          log.Printf("Failure: %+v\n", dockerconfig)
        } else {
          log.Printf("Success: %+v\n", dockerconfig)
        }
      }()
    }
  }
}

// Register requester
func (r *requester) Register() {
  log.Println("Requester")
  var dockerconfig dockercntrl.Config
  (*r.socket).Start(dockerconfig)
  log.Println("Started")
  r.self = r.handler.Requester.MakeInstance((*r.socket).Writer())
  log.Println(r.self)
  r.handler.Requester.Register <- r.self
  log.Println("Registered")
  go r.Run()
}

// Close the client connection
func (r *requester) Quit() {close(r.quit)}
