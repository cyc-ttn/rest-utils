package fileupload

import (
  "errors"
  "log"
  "sync"

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

var service * Service
var serviceOnce sync.Once

// FileHelper - interface used by Service
type FileHelper interface{
  GetLength()   uint
  GetUploaded() uint
  Write(string) (string, error)
  Append( []byte ) error
  IsComplete()  bool
}

// Service is a service to help with file uploads
type Service struct {
  files       map[ string ]FileHelper
  mu          sync.RWMutex
}

// Init - Initiailizes the Service
func Init() * Service{
  serviceOnce.Do(func() {
    service = &Service{
      files: make(map[string]FileHelper),
    }
  })

  return service
}

// Initialize initializes a file upload
func (s * Service) Initialize( f FileHelper, data []byte ) (string, bool, error) {
  id, err := uuid.NewV4()
  if err != nil { return "", false, err }

  idAsStr := id.String()
  log.Printf("ID: %s", idAsStr)

  lengthOfData := f.GetLength()
  uploadedLength := f.GetUploaded()

  s.mu.Lock()
  s.files[idAsStr] = f
  s.mu.Unlock()

  return idAsStr, lengthOfData == uploadedLength, nil
}

// Upload adds data to an existing file upload
func (s * Service) Upload( name string, data []byte ) (error) {

  s.mu.RLock()
  currFile, ok := s.files[name]
  s.mu.RUnlock()

  if !ok {
    return ErrNoSuchFile
  }

  return currFile.Append(data)
}

// AttemptWrite will try to write the file only if the file is totally available.
// Returns ("", nil) if file is incomplete
// Returns ("", error) if there is an error
// Returns ([filepath], nil) if successfully written.
func (s * Service) AttemptWrite(id string, filename string) (string, error) {
  s.mu.RLock()
  currFile, ok := s.files[id]
  s.mu.RUnlock()

  if !ok {
    return "", ErrNoSuchFile
  }

  if currFile.IsComplete() {
    filepath, err := currFile.Write(filename)
    if err != nil { return "", nil }
    delete(s.files, id)

    return filepath, nil
  }

  return "", nil
}

// Remove removes an existing upload (cancels)
func (s * Service) Remove(name string){
  delete(s.files, name)
}
