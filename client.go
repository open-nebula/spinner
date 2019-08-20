package spinner

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/gorilla/websocket"
  "time"
  "log"
)

// Single client (Captain) connection
type Client interface {
  // Accept reading from the client
  Read()
  // Accept writing from the client
  Write()
  // Accept reading and writing simultaneously in the background
  BackgroundRun()
  // Quit the client
  Quit()
}

type client struct {
  messenger     Messenger
  conn          *websocket.Conn
  spinup        chan *dockercntrl.Config
  quit          chan struct{}
}

const (
  pingPeriod = 20
  writeWait = 10
)

type clientPool   map[*client]bool
type clientChan   chan *client

func NewClient(m Messenger, conn *websocket.Conn) Client {
  return &client{
    messenger: m,
    conn: conn,
    spinup: make(chan *dockercntl.Config),
    quit: make(chan struct{}),
  }
}

func (c *client) Read() {
  defer func(){
    c.messenger.unregister <- c
    c.conn.Close()
  }()
  for {
  }
}

func (c *client) Write() {
  ticker := time.NewTicker(pingPeriod)
  defer func(){
    ticker.Stop()
    c.conn.Close()
  }()
  for {
    select {
    case config, ok := <- c.spinup:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      err := c.conn.WriteJSON(config)
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

func (c *client) BackgroundRun() {go c.Read(); go c.Write()}
func (c *client) Quit() {close(c.quit)}
