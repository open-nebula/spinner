package spinresp

import (
  "github.com/google/uuid"
)

type Response struct {
  Id      *uuid.UUID    `json:"id"`
  Code    int           `json:"code"`
  Data    struct{}      `json:"data"`
}

type ResponseChan chan Response

const (
  NoCaptainsAvailable = -6
)