package spinner

import (
  // "github.com/gorilla/websocket"
  "github.com/open-nebula/comms"
  "net/http"
  "log"
)

// On request adds client through the messenger
func join(handler *Handler) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    socket, err := comms.AcceptSocket(w,r)
    if err != nil {
      log.Println(err)
      return
    }
    client := NewClient(handler, &socket)
    client.Register()
  }
}

// Send a container to the messenger to pass on
func spin(handler *Handler) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    socket, err := comms.AcceptSocket(w,r)
    if err != nil {
      log.Println(err)
      return
    }
    requester := NewRequester(handler, &socket)
    requester.Register()
  }
}
