


build:
	go build

#
# Only call this to re-build the protobuf generated code
# Should only be necessary if you change the object model
# (add fields, etc)
#
protobuf:
	protoc -I=. -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/gogoproto --gofast_out=. pkg/dump1090/aircraft.proto
