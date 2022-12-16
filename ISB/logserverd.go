package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/zeromq/goczmq.v4"
)

type StringArrFlag []string

func (i *StringArrFlag) String() string {
	return strings.Join(*i, ",")
}

func Send(sock *goczmq.Sock, data []byte, flag int) {
	err := sock.SendFrame(data, flag)
	for err != nil {
		// fmt.Println("!!! SafeSendFrame ERROR", err)
		err = sock.SendFrame(data, flag)
	}
}

func Recv(sock *goczmq.Sock) [][]byte {
	for {
		msg, err := sock.RecvMessage()
		if err == nil {
			return msg
		}
	}
}

func (i *StringArrFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var Addr = flag.String("addr", "0.0.0.0:5552", "log service address")
var logFilenames StringArrFlag

func main() {
	flag.Var(&logFilenames, "f", "the path to log file")
	flag.Parse()

	if len(logFilenames) == 0 {
		fmt.Println("need at least one log file, specify by '-f'")
		return
	}

	allExist := true
	for _, filename := range logFilenames {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			allExist = false
			fmt.Printf("file not exist: %s\n", filename)
		}
	}
	if !allExist {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, logFilename := range logFilenames {
			endpoint := "ipc://" + logFilename + ".endpoint"
			fmt.Println(endpoint)
			logFlushSock, err := goczmq.NewReq(endpoint)
			for err != nil {
				logFlushSock, err = goczmq.NewReq(endpoint)
			}
			defer logFlushSock.Destroy()
			Send(logFlushSock, []byte("flush"), goczmq.FlagNone)
			msg := Recv(logFlushSock)
			fmt.Println(logFilename, "flush", msg)
		}

		// logBuf := new(bytes.Buffer)
		for _, logFilename := range logFilenames {
			log, err := ioutil.ReadFile(logFilename)
			if err != nil {
				fmt.Printf("read log file failed: %s\n", logFilename)
				continue
			}
			w.Write(log)
			w.Write([]byte("\n"))
		}

		w.WriteHeader(http.StatusOK)
		// w.Write(logBuf.Bytes())
	})

	fmt.Printf("log server listen at: %s\n", *Addr)
	if err := http.ListenAndServe(*Addr, mux); err != nil {
		panic(err)
	}
}
