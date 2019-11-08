package spinner

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "github.com/mitchellh/mapstructure"
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
  messenger     *messenger
  socket        *comms.Socket
  responses     spinresp.ResponseChan
  quit          chan struct{}
}

// Create new Client interface of client struct
func NewRequester(m *messenger, socket *comms.Socket) Requester {
  return &requester{
    messenger: m,
    socket: socket,
    quit: make(chan struct{}),
  }
}

// Get messages from the requester
func (r *requester) Run() {
  defer func(){
    (*r.socket).Close()
  }()
  read := (*r.socket).Reader()
  write := (*r.socket).Writer()
  for {
    select {
    case config, ok := <- read:
      log.Println("incoming")
      log.Println(config)
      if !ok {return}
      // dockerconfig, ok := config.(dockercntrl.Config)
      var dockerconfig dockercntrl.Config
      mapstructure.Decode(config, &dockerconfig)
      log.Println("inc2")
      log.Println(dockerconfig)
      dockerconfig.Cmd = []string{"echo", "hello"}
      // if !ok {log.Println("fail"); return}
      log.Println(dockerconfig)
      respchan := r.messenger.ContainerConnect(&dockerconfig)
      select{
      case resp, ok := <- respchan:
        if !ok {return}
        write <- resp
      }
    }
  }
}

// Register requester
func (r *requester) Register() {
  go r.Run()
  log.Println("requester")
}

// Close the client connection
func (r *requester) Quit() {close(r.quit)}
