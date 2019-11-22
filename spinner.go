// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
  "github.com/gorilla/mux"
  "github.com/phayes/freeport"
  "github.com/google/uuid"
  "net/http"
  "log"
  "strconv"
  "io/ioutil"
  "encoding/json"
  "time"
  "bytes"
)

type GeoIP struct {
	Ip         string  `json:"ip"`
	Lat        float32 `json:"latitude"`
	Lon        float32 `json:"longitude"`
}
type Spinner struct {
  Id          *uuid.UUID
  Ip          string
  Port        int
  LastUpdate  time.Time
}

// Server for the Nebula Spinner
type Server interface {
  // Given a port of 0, assigns a free port to the server.
  Run(port int)
}

type server struct {
  router    *mux.Router
  handler   *Handler
}

// Produces a new Server interface of struct server
func New() Server {
  router := mux.NewRouter().StrictSlash(true)
  handler := NewHandler()
  router.HandleFunc("/join", join(handler)).Name("Join")
  router.HandleFunc("/spin", spin(handler)).Name("Spin")
  handler.Start()
  return &server{
    router: router,
    handler: handler,
  }
}

// Runs the spinner server.
// If given a port value of 0, then finds a free port.
func (s *server) Run(port int) {
  var (
  	err      error
  	geo      GeoIP
  	response *http.Response
  	body     []byte
  )
  // get self IP and location info
  response, err = http.Get("http://api.ipstack.com/check?access_key=0bbaa9ccd131225ec08fa2c02c0a3260")
	if err != nil {
		log.Println(err)
	}
  body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
  err = json.Unmarshal(body, &geo)
	if err != nil {
		log.Println(err)
	}
  response.Body.Close()
  // start ping routine (pass self ip and port)
  go s.Ping(geo.Ip, port)

  if port == 0 {
    port, err = freeport.GetFreePort()
    if err != nil {log.Println(err); return}
  }
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), s.router))
}

func (s *server) Ping(ip string, port int) {
  ip = "no"  // hard code
  port = 1  // hard code
  // get id for spinner
  id := uuid.New()
  // periodically call register Api (ping)
  for {
    requestBody, err := json.Marshal(Spinner{
      Id:         &id,
      Ip:         "c2a13350.ngrok.io",
      Port:       80,
      LastUpdate: time.Now(),
    })
    if err!= nil {log.Println(err); return}

    resp, err := http.Post("http://c062a166.ngrok.io/register", "application.json", bytes.NewBuffer(requestBody))
    if err!= nil {log.Println(err); return}
    // no response needed
    resp.Body.Close()

    time.Sleep(3 * time.Second)
  }
}
