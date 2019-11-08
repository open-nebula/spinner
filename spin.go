package spinner

import (
  // "github.com/gorilla/websocket"
  "github.com/open-nebula/comms"
  "net/http"
  "log"
)

const (
  pingPeriod = 20
  writeWait = 10
)

// On request adds client through the messenger
func join(messenger *messenger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    socket, err := comms.AcceptSocket(w,r)
    if err != nil {
      log.Println(err)
      return
    }
    client := NewClient(messenger, &socket)
    client.Register()
  }
}

// Send a container to the messenger to pass on
func spin(messenger *messenger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    socket, err := comms.AcceptSocket(w,r)
    if err != nil {
      log.Println(err)
      return
    }
    requester := NewRequester(messenger, &socket)
    requester.Register()
  }
}
