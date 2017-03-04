package objectmanager

import (
  "errors"
  "log"
)

// Object manager
var objects = initManager()

// Errors
var (
  ErrDuplicateObject = errors.New("Object already exists")
  ErrObjectDoesNotExist = errors.New("Object does not exist")
)

func initManager()  map[string]interface{}{
  return make(map[string]interface{})
}

// Set sets an object into the manager for retrieval
func Set(name string, obj interface{}) error{
  log.Printf("Storing object '%s'", name)

  if _, ok := objects[name]; ok {
    return ErrDuplicateObject
  }

  objects[name] = obj
  return nil
}

// Get retrieves an object from the object manager
func Get(name string) (interface{}, error){
  val, ok := objects[name]
  if !ok {
    return nil, ErrObjectDoesNotExist
  }

  return val, nil
}
