package spinner

import (
  "github.com/gorilla/websocket"
  "time"
  "log"
)

type Client struct {
  Messenger     *Messenger
  Conn          *websocket.Conn
}

const (
  pingPeriod = 20
  writeWait = 10
)

type ClientPool   map[*Client]bool
type ClientChan   chan *Client

func NewClient(m *Messenger, conn *websocket.Conn) *Client {
  client := &Client{
    Messenger: m,
    Conn: conn,
  }
  client.Read()
  client.Write()
  return client
}

func (c *Client) Read() {
  go func(){
    defer func(){
      c.Messenger.Unregister <- c
      c.Conn.Close()
    }()
  }()
}

func (c *Client) Write() {
  go func(){
    ticker := time.NewTicker(pingPeriod)
    defer func(){
      ticker.Stop()
      c.Conn.Close()
    }()

    for {
      select {
      case <- ticker.C:
        c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
        if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
          log.Println(err)
          return
        }
      }
    }
  }()
}
