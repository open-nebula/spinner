package spinner

import (
  "github.com/gorilla/websocket"
  "net/http"
  "log"
)

func Spin() {
  messenger := NewMessenger()
  return func(w http.ResponseWriter, r *http.Request) {
    conn, err := messenger.Upgrader.Upgrade(w,r,nil)
    if err != nil {
      log.Println(err)
      return
    }
    NewClient(messenger, conn)
  }
}
