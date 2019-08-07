include ~/gopath/src/github.com/golang/protobuf/Make.protobuf

all_proto = $(wildcard *.proto)
all_pb_go = $(all_proto:.proto=.pb.go)

all: $(all_pb_go)

clean:
	rm -rf $(all_pb_go)

