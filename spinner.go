package spinner

import (
  "github.com/gorilla/mux"
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
