package main

import (
	"net/rpc"
	"log"
	"net"
	"net/http"
)

//Holds arguments to be passed to service Arith in RPC call
type Args struct {
	A, B int
}

//Representss service Arith with method Multiply
type Arith int

//Result of RPC call is of this type
type Result int

func main() {
	arith := new(Arith)
	err := rpc.Register(arith)
	if err != nil {
		log.Fatalf("Format of service Arith isn't correct. %s", err)
	}
	rpc.HandleHTTP()
	//start listening for messages on port 1234
	l, e := net.Listen("tcp", ":2378")
	if e != nil {
		log.Fatalf("Couldn't start listening on port 1234. Error %s", e)
	}
	log.Println("Serving RPC handler")
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}

//This procedure is invoked by rpc and calls rpcexample.Multiply which stores product of args.A and args.B in result pointer
func (t *Arith) Multiply(args Args, result *Result) error {
	return Multiply(args, result)
}

//stores product of args.A and args.B in result pointer
func Multiply(args Args, result *Result) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result(args.A * args.B)

	return nil
}
