# goshnix

![](https://github.com/swarvanusg/goshnix/blob/master/doc/goshnix%20(1).png)

A powerful UNIX utility for performing Golang API over SSH. It abstracts the underlying method for Goshnix supported API, so that your program logic stays same

[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/swarvanusg/goshnix)

### Version
1.0.0

### Get It
```go
go get github.com/swarvanusg/goshnix
```


### Use it  
##### Initiate Goshnix Client
```go
goshclient, err := goshnix.Init("<host_ip>", "<port>", "<uname>", "<pass>")
```

##### Use Goshnix Client as a Std Lib
Check you environment variable
```go
envval, err := goshclient.Getenv("<key>")
```
Get a file stat
```go 
fileinfo, _ := goshclient.Stat("<filepath>")
// Check if its a dir (as of std lib)
if fileinfo.IsDir() {
  // ...
}
```
Get the list of directory entries
```go
var fileinfos []os.FileInfo
fileinfos, _ := goshclient.ReadDir("<dirpath>")
```
For full listing of supported API see the godoc.

##### Use Goshnix API to Understand Error
Goshnix returned error could be communication error as well as command execution error. Goshnix provides API to understand proper error sub-type 
```go
filecontent, err := goshclient.ReadFile(filename string)
if err != nil {
  // Check error type
  switch {
  case Isssherror(err):
    // ssh error handling logic
  case Iscmderror(err):
    // cmd error handling logic
  }
}
```

### How it works
For each of the API, one or multiple UNIX commands gets executed over ssh to the targeted host.

### License
MIT
