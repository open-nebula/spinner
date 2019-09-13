package spinner

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/gorilla/websocket"
  "time"
  "log"
)

// Single container requestor interface
type Requester interface {
  // Accept reading from the client
  Read()
  // Accept writing from the client
  Write()
  // Enter client into messenger system
  Register()
  // Quit the client
  Quit()
}

type requester struct {
  messenger     Messenger
  conn          *websocket.Conn
  responses     spinresp.ResponseChan
  quit          chan struct{}
}

const (
  pingPeriod = 20
  writeWait = 10
)

// Create new Client interface of client struct
func NewRequester(m Messenger, conn *websocket.Conn) Requester {
  return &requester{
    messenger: m,
    conn: conn,
    quit: make(chan struct{}),
  }
}

// Get messages from the requester
func (r *requester) Read() {
  defer func(){
    c.conn.Close()
  }()
  for {
    var config dockercntrl.Config
    err := r.conn.ReadJSON(&config)
    if err != nil {
      log.Println(err)
      return
    }
    respchan := r.messenger.ContainerConnect(&config)
    r.responses = respchan
    go r.Write() // TODO: Must be idempotent
  }
}

// Send messages to the requester.
func (r *requester) Write() {
  ticker := time.NewTicker(pingPeriod)
  defer func(){
    ticker.Stop()
    c.conn.Close()
  }()
  for {
    select {
    case resp, ok := <- r.responses:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      err := c.conn.WriteJSON(resp)
      if err != nil {
        log.Println(err)
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
    case <- ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
        log.Println(err)
        return
      }
    case <- c.quit:
      return
    }
  }
}

// Register requester
func (r *requester) Register() {
  go c.Read()
}

// Close the client connection
func (r *requester) Quit() {close(c.quit)}
