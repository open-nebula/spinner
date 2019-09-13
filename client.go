package spinner

import (
  "github.com/open-nebula/spinner/spinresp"
  "github.com/gorilla/websocket"
  "github.com/google/uuid"
  "time"
  "log"
)

// Single client (Captain) connection
type Client interface {
  // Accept reading from the client
  Read()
  // Accept writing from the client
  Write()
  // Enter client into messenger system
  Register()
  // Quit the client
  Quit()
}

type client struct {
  messenger     *messenger
  conn          *websocket.Conn
  spinup        chan message
  responses     map[uuid.UUID]spinresp.ResponseChan
  quit          chan struct{}
}

type clientPool   map[*client]bool
type clientChan   chan *client

// Create new Client interface of client struct
func NewClient(m *messenger, conn *websocket.Conn) Client {
  return &client{
    messenger: m,
    conn: conn,
    spinup: make(chan message),
    responses: make(map[uuid.UUID]spinresp.ResponseChan),
    quit: make(chan struct{}),
  }
}

// Get messages from the client
func (c *client) Read() {
  defer func(){
    c.messenger.unregister <- c
    c.conn.Close()
  }()
  for {
    var resp spinresp.Response
    err := c.conn.ReadJSON(&resp)
    if err != nil {
      log.Println(err)
      return
    }
    if resp.Id != nil {return}
    if respchan, ok := c.responses[*resp.Id]; ok {
      respchan <- resp
      if resp.Code < 0 {delete(c.responses, *resp.Id); return}
    } else {return}
  }
}

// Send messages to the client.
// Currently messages are only of the type dockercntrl.Config
// as a request to spin up a container of that type.
func (c *client) Write() {
  ticker := time.NewTicker(pingPeriod)
  defer func(){
    ticker.Stop()
    c.conn.Close()
  }()
  for {
    select {
    case message, ok := <- c.spinup:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      if message.config.Id == nil {
        identifier := uuid.New()
        message.config.Id = &identifier
      }
      err := c.conn.WriteJSON(message.config)
      if err != nil {
        log.Println(err)
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      c.responses[*message.config.Id] = message.response
    case <- ticker.C:
      // c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      // if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
      //   log.Println(err)
      //   return
      // }
    case <- c.quit:
      return
    }
  }
}

// Register client with messenger and accept read/writes.
func (c *client) Register() {
  c.messenger.register <- c
  go c.Read()
  go c.Write()
  log.Println("Client registered")
}

// Close the client connection
func (c *client) Quit() {close(c.quit)}
