all:
	cp ../TL-Proto/id/service.pb.go ./proto/id/
	GO111MODULE=off go build -o TL-ID-Service main.go
