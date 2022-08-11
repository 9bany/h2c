package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/posener/h2conn"
)

var ctx context.Context
var cancel context.CancelFunc

const url = "https://staging.chatbot.iviet.com:443/connect"

type Listener struct{}

func (listener *Listener) Write(data []byte) (n int, err error) {

	if len(data) == 0 {
		return 0, nil
	}
	return listener.readData(data)
}

func (listener *Listener) readData(data []byte) (int, error) {
	dataString := string(data)
	log.Println(dataString)
	cancel()
	return len(data), nil
}

func main() {
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// go catchSignal(cancel)q
	d := &h2conn.Client{Method: http.MethodGet, Header: http.Header{
		"device-id":       {"5D2D7002-9F0C-4DBC-BE72-D564643108F1"},
		"device-type":     {"ios"},
		"timezone":        {"Asia/Ho_Chi_Minh"},
		"version-info":    {"2.4.2"},
		"meta":            {"eyJldmVudCI6eyJoZWFkZXIiOnsiZGlhbG9nUmVxdWVzdElkIjoiRUY1RUI0MjMtNTNBRi00NjRFLUJDMDYtQjk4NUQ0QTdCMENGIiwibmFtZXNwYWNlIjoiU3BlZWNoUmVjb2duaXplciIsInJhd1NwZWVjaCI6Ik3hu58gxJHDqG4gcGjDsm5nIG5n4bunIiwibmFtZSI6IlJlY29nbml6ZSIsIm1lc3NhZ2VJZCI6Im1lc3NhZ2VJZC0zMDdEMUMwNzdGQ0M0RDQyODRBMkVDMTJDMzlBRDU2NyJ9fX0="},
		"olli-session-id": {"5D2D7002-9F0C-4DBC-BE72-D564643108F1"},
		"Authorization":   {"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2NTUxMDU0MjUsIm5iZiI6MTY1NTEwNTQyNSwianRpIjoiZjAxNzdkYzgtMzVmZC00NDI1LTg0NzQtNTZiNGI1M2MzM2UxIiwiaWRlbnRpdHkiOiJ7XCJzdWJcIjogNzM1LCBcIm5hbWVcIjogXCJcIiwgXCJlbWFpbFwiOiBudWxsLCBcInJvbGVcIjogMSwgXCJzdGF0dXNcIjogMSwgXCJkZXZpY2VfaWRcIjogXCJcIiwgXCJkZWZhdWx0X2xhbmd1YWdlXCI6IFwidmktVk5cIiwgXCJleHByaXJhdGlvblwiOiA4NjQwMCwgXCJwaG9uZV9udW1iZXJcIjogXCIrODQzNzg3ODA4NDNcIiwgXCJjYWxsaW5nX25hbWVcIjogXCJcIn0iLCJmcmVzaCI6ZmFsc2UsInR5cGUiOiJhY2Nlc3MifQ.NVLxRK-dKYjskDKqrRVm2pwBO0AIorMoCWfTjjBvHKg"},
	}}

	conn, resp, err := d.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Initiate conn: %s", err)
	}
	defer conn.Close()

	// Check server status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad status code: %d", resp.StatusCode)
	}

	defer log.Println("Exited")
                
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(conn, "pong")
		}
	}()

	// Loop until user terminates
	fmt.Println("Echo session starts, press ctrl-C to terminate.")
	for ctx.Err() == nil {
		_, err := io.Copy(&Listener{}, conn)
		if err != nil {
			log.Fatalf("Failed receiving message: %v", err)
		}
	}
}

func catchSignal(cancel context.CancelFunc) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Cancelling due to interrupt")
	cancel()
}
