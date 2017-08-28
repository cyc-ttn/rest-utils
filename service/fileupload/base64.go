package fileupload

import (
  "encoding/base64"
  "io/ioutil"
  "log"
)

// Base64FileUploadHelper defines a file for upload
type Base64FileUploadHelper struct{
  Buffer      []byte
  Extension   string
  Length      uint
}

// InitializeBase64FileUpload inits a standard file upload
func InitializeBase64FileUpload(length uint, extension string, data []byte) (string, bool, error ){

  // Get the service.
  service := Init()

  // Create the File
  f := &Base64FileUploadHelper{
    Extension: extension,
    Length: length,
  }

  lengthOfdata := len(data)

  if uint(lengthOfdata) == length {
    f.Buffer = data
  }else{
    buf := make([]byte, lengthOfdata, length)
    copy( buf, data )
    f.Buffer = buf
  }

  return service.Initialize(f, data)
}

// Append adds data to the buffer
func (f * Base64FileUploadHelper) Append(data []byte) error {
  lengthOfData := uint(len(data))
  length := f.GetUploaded()
  total := lengthOfData + length

  if total > f.Length {
    return ErrBufferOverflow
  }

  f.Buffer = f.Buffer[0:total]
  copy( f.Buffer[length:], data)
  return nil
}

// Decode decodes the file from base64 to bytes
func (f * Base64FileUploadHelper) Decode() ([]byte, error) {
  return base64.StdEncoding.DecodeString( string( f.Buffer ) )
}

// Write writes data to the file
func (f * Base64FileUploadHelper) Write(name string) (string, error){
  if len(f.Buffer) != int(f.Length) {
      return "", ErrBufferIncomplete
  }

  contents, err := f.Decode()
  if err != nil { return "", err }

  filepath := name + "." + f.Extension
  log.Printf("Filepath: %s", filepath)

  err = ioutil.WriteFile("./" + filepath, contents, 0644)
  if err != nil { return "", err }

  return filepath, nil
}

// GetLength returns the total length of the file
func (f * Base64FileUploadHelper) GetLength() uint { return f.Length }

// GetUploaded returns number of bytes already uploaded
func (f * Base64FileUploadHelper) GetUploaded() uint{ return uint(len(f.Buffer)) }

// IsComplete - returns true if buffer is complete
func (f * Base64FileUploadHelper) IsComplete() bool { return f.Length == f.GetUploaded() }
