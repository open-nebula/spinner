// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
  "github.com/gorilla/mux"
  "github.com/phayes/freeport"
  "net/http"
  "log"
  "strconv"
)

// Server for the Nebula Spinner
type Server interface {
  // Given a port of 0, assigns a free port to the server.
  Run(port int)
}

type server struct {
  router    *mux.Router
  messenger Messenger
}

// Produces a new Server interface of struct server
func New() Server {
  router := mux.NewRouter().StrictSlash(true)
  messenger := NewMessenger()
  router.HandleFunc("/join", join(messenger)).Name("Join")
  router.HandleFunc("/spin", spin(messenger)).Name("Spin")
  go messenger.Run()
  return &server{
    router: router,
    messenger: messenger,
  }
}

// Runs the spinner server.
// If given a port value of 0, then finds a free port.
func (s *server) Run(port int) {
  var err error
  if port == 0 {
    port, err = freeport.GetFreePort()
    if err != nil {log.Println(err); return}
  }
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), s.router))
}
