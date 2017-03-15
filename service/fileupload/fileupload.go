package fileupload

import (
  "encoding/base64"
  "errors"
  "io/ioutil"
  "log"

  "luxe.technology/rest-utils/objectmanager"

  "github.com/nu7hatch/gouuid"
)

var (
  // ErrNoSuchFile thrown when the file to upload to doesn't exist
  ErrNoSuchFile = errors.New("No such file exists")

  // ErrBufferOverflow thrown when buffer overflows
  ErrBufferOverflow = errors.New("File is too large")

  // ErrBufferIncomplete thrown when buffer is incomplete but tries to write
  ErrBufferIncomplete = errors.New("Buffer is still awaiting data")
)

// File defines a file for upload
type File struct{
  Buffer      []byte
  Extension   string
  Length      uint
}

// Service is a service to help with file uploads
type Service struct {
  files       map[ string ]*File
}

func init(){
  service := &Service{
    files: make(map[string]*File),
  }

  objectmanager.Set("FileUploadService", service)
}

// Initialize initializes a file upload
func (s * Service) Initialize( length uint, extension string, data []byte ) (string, bool) {
  id, err := uuid.NewV4()
  if err != nil { return "", false }

  idAsStr := id.String()
  log.Printf("ID: %s", idAsStr)

  lengthOfdata := len(data)

  f := &File{
    Extension: extension,
    Length: length,
  }

  if uint(lengthOfdata) == length {
    // Try to write
    f.Buffer = data
    filepath, err := f.write(idAsStr)
    if err != nil { return "", false }
    return filepath, true
  }

  buf := make([]byte, lengthOfdata, length)
  copy( buf, data )

  f.Buffer = buf
  s.files[ idAsStr ] = f

  return idAsStr, false
}

// Upload adds data to an existing file upload
func (s * Service) Upload( name string, data []byte ) (string, error) {

  currFile, ok := s.files[name]
  if !ok {
    return "", ErrNoSuchFile
  }

  length := len(currFile.Buffer)
  capacity := cap(currFile.Buffer)

  lengthOfData := len(data)
  total := length + lengthOfData

  if total > capacity{
    return "", ErrBufferOverflow
  }

  currFile.Buffer = currFile.Buffer[0:length+lengthOfData]
  copy( currFile.Buffer[length:], data)

  filepath, err := currFile.write(name)
  if err != nil { return "", nil }
  delete(s.files, name)

  return filepath, nil
}

// Remove removes an existing upload (cancels)
func (s * Service) Remove(name string){
  delete(s.files, name)
}

func (f * File) write(name string) (string, error){
  if len(f.Buffer) != int(f.Length) {
      return "", ErrBufferIncomplete
  }

  contents, err := base64.StdEncoding.DecodeString( string( f.Buffer ) )
  if err != nil { return "", err }

  filepath := "products/images/" + name + "." + f.Extension
  log.Printf("Filepath: %s", filepath)

  err = ioutil.WriteFile("./" + filepath, contents, 0644)
  if err != nil { return "", err }

  return filepath, nil
}
