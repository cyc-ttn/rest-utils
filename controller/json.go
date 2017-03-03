package controller

// Helps with JSON type requests and responses

import (
  "encoding/json"
  "net/http"
)

// JSONFromRequest returns a map from a JSON value from a request
func JSONFromRequest( r * http.Request ) (map[string]interface{}, error) {

  decoder := json.NewDecoder(r.Body)

  t := make(map[string]interface{})
  err := decoder.Decode(&t)
  defer r.Body.Close()

  if err != nil {
    return nil, err
  }

  return t, nil
}

// JSONFromRequestToStruct returns a Struct from a JSON value from a request
func JSONFromRequestToStruct( r * http.Request, v interface{}) error{
  decoder := json.NewDecoder(r.Body)

  err := decoder.Decode(&v)
  defer r.Body.Close()

  if err != nil {
    return err
  }

  return nil
}

// JSONHandlerFunc Wraps the input function and passes in data as a JSON variable
func JSONHandlerFunc( f func(http.ResponseWriter,*http.Request,map[string]interface{},error ) ) http.HandlerFunc {

  h := func(w http.ResponseWriter, r * http.Request){
    t, err := JSONFromRequest(r)
    f(w,r,t,err)
  }

  return h
}

// JSONWriter Writes JSON as output
func JSONWriter( w http.ResponseWriter, object interface{} ) error {

  w.Header().Set("Content-Type", "application/json")

  bytes, err := json.Marshal(object)
  if err != nil{
    return err
  }

  w.Write( bytes )
  return nil
}
