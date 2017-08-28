package controller

import (
  "errors"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"
)

// IDFromRequest extracts the ID from the request
func IDFromRequest(req * http.Request) (uint,error) {
  vars := mux.Vars(req)
  idAsStr, ok := vars["id"]
  if !ok { return 0, errors.New("ID is required") }

  id, err := strconv.Atoi(idAsStr)
  if err != nil { return 0, errors.New("ID must be an integer") }

  return uint(id), nil
}
