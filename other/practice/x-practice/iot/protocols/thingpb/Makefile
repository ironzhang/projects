all_proto = $(wildcard *.proto)
all_pb_go = $(all_proto:.proto=.pb.go)

all: $(all_pb_go)

$(all_pb_go): $(all_proto)
	protoc --go_out=. $(all_proto)

clean:
	rm -rf $(all_pb_go)

