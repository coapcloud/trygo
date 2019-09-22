package main

import (
	"context"
	"fmt"
	"log"
	"time"

	coap "github.com/go-ocf/go-coap"
	flag "github.com/spf13/pflag"
)

var port = flag.IntP("port", "p", 5683, "coap port to listen on")

func handleA(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message in handleA: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	w.SetContentFormat(coap.TextPlain)
	log.Printf("Transmitting from A")
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if _, err := w.WriteWithContext(ctx, []byte("hello world")); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func handleB(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message in handleB: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	resp := w.NewResponse(coap.Content)
	resp.SetOption(coap.ContentFormat, coap.TextPlain)
	resp.SetPayload([]byte("good bye!"))
	log.Printf("Transmitting from B %#v", resp)
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if err := w.WriteMsgWithContext(ctx, resp); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func main() {
	flag.Parse()

	mux := coap.NewServeMux()
	mux.Handle("/a", coap.HandlerFunc(handleA))
	mux.Handle("/b", coap.HandlerFunc(handleB))

	fmt.Printf("serving coap requests on %d\n", *port)
	log.Fatal(coap.ListenAndServe("udp", fmt.Sprintf(":%d", *port), mux))
}
