package spinner

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/comms"
  "github.com/google/uuid"
  "log"
)
// task scheduling request
type Request struct {
  Success  chan bool  // scheduled or not
  Task     *Task  // task meta data
}

type Task struct {
  Config  *dockercntrl.Config
  From    *uuid.UUID
}

type Handler struct {
  clients         *comms.Messenger
  clientMetaData  map[uuid.UUID]int
  Requester       *comms.Messenger
  Register        chan *comms.Instance
  Unregister      chan *comms.Instance
  Request         chan *Request
}

func NewHandler() *Handler {
  h := &Handler{
    clients: comms.NewMessenger(),  // all captains'messenger
    clientMetaData: make(map[uuid.UUID]int),  // captains data in handler (there is also one in messenger)
    Requester: comms.NewMessenger(),  // ?
    Register: make(chan *comms.Instance),  // captain register queue in handler (there is also one in messenger)
    Unregister: make(chan *comms.Instance),  // unregister queue in handler (there is also one in messenger)
    Request: make(chan *Request),  // task queue in handler (there is also one in messenger)
  }
  h.clients.Start()  // start captains'messenger loop
  h.Requester.Start()
  return h
}

func (h *Handler) run() {
  defer func() {
    log.Println("Handler Complete")
  }()
  for {
    log.Println("Handler Action")
    select {
    // new captain registered
    case client := <- h.Register:
      h.clientMetaData[*client.Id] = 0  // initialize this new captain
      h.clients.Register <- client  // push into messenger's register queue
    case client := <- h.Unregister:
      delete(h.clientMetaData, *client.Id)
      h.clients.Unregister <- client
    // new task is scheduled: select a random captain from map "clientMetaData"
    case request := <- h.Request:
      // Round-Robin, extract away to Schedule type
      log.Printf("Round Robin Scheduling\n")
      minimum := 0
      var chosen *uuid.UUID
      for k,v := range h.clientMetaData {
        if chosen == nil || v < minimum {
          chosen = &k; minimum = v
        }
      }
      if chosen == nil {request.Success <- false; break}
      log.Printf("Chosen: %+v\n", chosen)
      // after random selection: push the task into messenger's scheduling queue
      h.clients.Message <- &comms.Message{
        Success: request.Success,
        Reciever: chosen,
        Data: request.Task,
      }
    }
  }
}

func (h *Handler) Start() {go h.run()}
// after receive task from user, push the task into the Request queue in handler waiting for scheduling
func (h *Handler) SendTask(from *comms.Instance, task *dockercntrl.Config) bool {
  response := make(chan bool)
  req := &Request{
    Success: response,
    Task: &Task{
      From: from.Id,
      Config: task,
    },
  }
  h.Request <- req
  // block until the succ bool value is returned (after this task is pushed into the spinup queue!)
  status := <- response
  return status
}
