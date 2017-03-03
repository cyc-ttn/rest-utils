package model;

import (
  "log"
  "reflect"
  "github.com/jinzhu/gorm"
)

var objects []interface{}
var count int

// Model type
type Model interface{
  GetID() uint
}

func init(){
  objects = make([]interface{}, 0, 10)
  count = 0
}

// Generate uses the ORM capability of gorm
// to migrate
func Generate( db *gorm.DB ){
  log.Println("Model: Generating models...")
  for _, object := range objects{
    log.Printf("Model: >> Generating model for %s", reflect.TypeOf(object).Elem().Name() )
    db.AutoMigrate(object)
  }
}

// Add adds an object to the object list
func Add( obj interface{} ){
  length := cap(objects)

  if( count >= length ){
    newObj := make([]interface{}, length, length + 10)
    copy(newObj, objects)
    objects = newObj
  }
  
  objects = append(objects, obj)
  count = count + 1
}
