package spinner

import (
  "github.com/gorilla/mux"
  "github.com/phayes/freeport"
  "net/http"
  "log"
  "strconv"
)

type Spinner struct {
  Router    *mux.Router
}

func New() *Spinner {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/spin", Spin()).Name("Spin")

  return &Spinner{
    Router: router,
  }
}

func (s *Spinner) Run(port int) {
  var err error
  if port == 0 {
    port, err = freeport.GetFreePort()
    if err != nil {log.Println(err); return}
  }
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), s.Router))
}
