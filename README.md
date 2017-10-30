# GSSTrainingSystem
This is a light weight LMS that serves 'courses' for training purposes

## Installing assets
To install the assets and dependencies, from the 'assets' directory run:

  `bower install`
  
##Running the service
The simple thing to do is to Clone this repository and run 
`./server` To run the pre-built executable.

To re-compile the server, you'll need to have a golang environment setup. 
(see https://golang.org/doc/install for a getting started guide)
Once you have Golang setup locally, clone this repository into your GOPATH, and run 
`go run main.go` This will compile and run the server. If you just want to rebuild the executable, run 
`go build -o server`