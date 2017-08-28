package controller

import (
  "os"
  "net/http"
  "strings"
)

// ResourceAndStaticFileHandler -
// if the path is a directory, return the static.
// if the path is a file, check to see if the file already exists
func ResourceAndStaticFileHandler(
  testDirs []string,
  static string,
) func(http.ResponseWriter, *http.Request) {

  return func( w http.ResponseWriter, req * http.Request){

    path := req.URL.Path

    file, err := os.Stat(path)
    if err == nil && !file.IsDir() {
      http.ServeFile(w, req, path)
      return
    }

    // Fix the slashes by removing first slash, and last slash if exists.
    if strings.HasSuffix(path, "/") { path = path[1: len(path)-1] }

    // Split into parts
    parts := strings.Split(path, "/")
    length := len(parts)

    //If is a path, return the static file
    if strings.Index(parts[length-1], ".") == -1 {
      http.ServeFile(w, req, static)
      return
    }

    // Find existing files
    _testDirs := make([]string, len(testDirs)+1)
    _testDirs[0] = "./"
    for i, t := range testDirs {
      if strings.HasPrefix(t, "/"){  t = t[1:] }
      if strings.HasPrefix(t, "./"){ t = t[2:] }
      if !strings.HasSuffix(t, "/"){ t = t + "/" }
      _testDirs[i+1] = "./" + t
    }

    // Go through all directories including current directory
    for i:=0; i < length; i++ {
      testPath := strings.Join( parts[i:length], "/")
      for _, d := range _testDirs {
        if _, err := os.Stat("./" + d + testPath); err == nil{
          http.ServeFile(w, req, d + testPath)
          return
        }
      }
    }

    http.ServeFile(w, req, static)
  }
}
