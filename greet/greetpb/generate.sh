export PATH=$PATH:$(go env GOPATH)/bin
#protoc -I greet/greetpb --go_out=greet/greetpb greet/greetpb/greet.proto
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:greet/greetpb/