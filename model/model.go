package model;

import (
  "log"
  "github.com/jinzhu/gorm"
)

var objects []interface{}
var count int

func init(){
  objects = make([]interface{}, 0, 10)
  count = 0
}

// Generate uses the ORM capability of gorm
// to migrate
func Generate( db *gorm.DB ){
  log.Println("Generating models...")
  db.AutoMigrate(objects...)
}

// Add adds an object to the object list
func Add( obj interface{} ){
  length := cap(objects)

  if( count >= length ){
    newObj := make([]interface{}, 0, length + 10)
    copy(newObj, objects)
    objects = newObj
  }

  objects = append(objects, obj)
  count = count + 1
}
