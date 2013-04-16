# Foodle - The USDA Nutrient Database Search Engine

This is a sample app for a tech talk given at the [AngularJS NYC Meetup](http://www.meetup.com/angularjs-nyc) in April 2013.

The repository is tagged with with varous `preso/step-N` tags, which match up with the various step slides in the presentation.

**TODO: Post link to slides.**

## Running the Server

The backend server is written in [Go](http://golang.org), which the only requirement for building the app.

1. Download and install the [Go runtime](http://golang.org/doc/install).
2. Set up your GOPATH, e.g.:
  * `export GOPATH=$HOME/go/`
  * `mkdir $GOPATH`
3. Have Go clone the repository for you:
  * `go get -d github.com/rsesek/usda-ndb`
4. Go to the work directory and build the server:
  * `cd $GOPATH/src/github.com/rsesek/usda-ndb`
  * `go build`
5. Run the server:
  * `./usda-ndb`

The server by default runs on port 8077, but it can be changed with the `-port=8077` flag to the binary.
