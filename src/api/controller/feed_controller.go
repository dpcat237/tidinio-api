package controller

import (
	"net/http"
	"github.com/tidinio/src/core/component/controller"
	"net/rpc"
	"log"
)

type addFeed struct {
	FeedUrl string `json:"feed_url"`
}

//Holds arguments to be passed to service Arith in RPC call
type Args struct {
	A, B int
}

//Representss service Arith with method Multiply
type Arith int

//Result of RPC call is of this type
type Result int


func AddFeed(w http.ResponseWriter, r *http.Request) {
	data := addFeed{}
	_, err := common_controller.GetAuthContent(w, r, &data)
	if err != nil {
		return
	}

	/*feed, err := feed_handler.AddFeed(user.ID, data.FeedUrl)
	if err != nil {
		common_controller.ReturnPreconditionFailed(w, "Wrong url")
	}
	common_controller.ReturnJson(w, feed)*/

	log.Println("tut: AddFeed a")
	sendCommand()
	log.Println("tut: AddFeed b")
}

func sendCommand() {
	//make connection to rpc server
	client, err := rpc.DialHTTP("tcp", ":2378")
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}
	//make arguments object
	args := &Args{
		A: 2,
		B: 3,
	}
	//this will store returned result
	var result Result
	//call remote procedure with args
	err = client.Call("Arith.Multiply", args, &result)
	if err != nil {
		log.Fatalf("error in Arith", err)
	}
	//we got our result in result
	log.Printf("%d*%d=%d\n", args.A, args.B, result)
}

