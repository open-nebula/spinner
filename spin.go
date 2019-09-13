package spinner

import (
  // "github.com/gorilla/websocket"
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
    conn, err := messenger.upgrader.Upgrade(w,r,nil)
    if err != nil {
      log.Println(err)
      return
    }
    client := NewClient(messenger, conn)
    client.Register()
  }
}

// Send a container to the messenger to pass on
func spin(messenger *messenger) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    conn, err := messenger.upgrader.Upgrade(w,r,nil)
    if err != nil {
      log.Println(err)
      return
    }
    requester := NewRequester(messenger, conn)
    requester.Register()
  }
}
