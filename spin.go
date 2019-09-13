package spinner

import (
  // "github.com/gorilla/websocket"
  "net/http"
  "log"
)

// On request adds client through the messenger
func join(messenger Messenger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    conn, err := messenger.Upgrader.Upgrade(w,r,nil)
    if err != nil {
      log.Println(err)
      return
    }
    client := NewClient(messenger, conn)
    client.Register()
  }
}

// Send a container to the messenger to pass on
func spin(messenger Messenger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    conn, err := messenge.r Upgrader.Upgrade(w,r,nil)
    if err != nil {
      log.Println(err)
      return
    }
    requester := NewRequester(messenger, conn)
    requeter.Register()
  }
}
