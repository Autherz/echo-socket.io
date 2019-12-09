package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ali-IoT-Lab/socketio-client-go"
)

type Test struct {
	MaterialId string  `json:"materialid"`
	Process    float64 `json:"process"`
}

type Test2 struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

func main() {
	var Header http.Header = map[string][]string{
		"moja":     {"ccccc, asdasdasdasd"},
		"terminal": {"en-esadasdasdwrw"},
		"success":  {"dasdadas", "wdsadaderew"},
	}
	fmt.Println("-------------------request.header-------------------------")
	fmt.Println(Header)
	// s, err := socketio.Socket("ws://172.18.0.5:3001")
	s, err := socketio.Socket("ws://localhost:1333")
	if err != nil {
		panic(err)
	}
	s.Connect(Header)

	// s.Emit("subscribe", "roomTwo")
	s.Emit("subscribe", "roomOne")

	s.On("message", func(args ...interface{}) {
		fmt.Println("servver message!")
		fmt.Println(args[0])
	})

	test_temp2, _ := json.Marshal(&Test2{
		// Room:    "roomTwo",
		Room:    "roomOne",
		Message: "Test2",
	})

	s.Emit("send", string(test_temp2))

	// i := 0.0
	for {
		// i++
		// test_temp, _ := json.Marshal(&Test{
		// 	MaterialId: "xxxxxxxxxxxxxxxxxxxzzzzzzzzzzzzzzz",
		// 	Process:    i,
		// })

		// s.Emit("transcode", string(test_temp))
		// s.Emit("messgae", string(test_temp))
		// time.Sleep(2 * time.Second)
		select {}
	}
}
