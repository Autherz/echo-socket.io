package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
)

type Test2 struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

type VideosProcess struct {
	MaterialId     string  `json:"materialid"`
	MaterialPackId string  `json:"materialpackid"`
	UserId         string  `json:"userid"`
	Process        float64 `json:"process"`
	Status         bool    `json:"status"`
}

func main() {
	e := echo.New()

	pt := polling.Default

	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			pt,
			wt,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected: ", s.ID())
		return nil
	})

	server.OnEvent("/", "subscribe", func(s socketio.Conn, room string) {
		s.SetContext(room)
		server.JoinRoom(room, s)
	})

	server.OnEvent("/", "unsubscribe", func(s socketio.Conn, room string) {
		fmt.Println("Leave Room : ", room)
		server.LeaveRoom(room, s)

		server.OnDisconnect("/", func(s socketio.Conn, msg string) {
			fmt.Println("closedddddd", msg)
		})

		fmt.Println(server.Rooms())

	})

	server.OnEvent("/", "transcode", func(s socketio.Conn, msg string) {
		fmt.Println("connected: ", s.ID())
		s.SetContext(msg)
		var v VideosProcess
		json.Unmarshal([]byte(msg), &v)
		fmt.Println(v)
		server.BroadcastToRoom(v.MaterialPackId, "message", msg)
	})

	server.OnEvent("/", "send", func(s socketio.Conn, msg string) {
		fmt.Println("connected: ", s.ID())
		s.SetContext(msg)
		var test Test2
		json.Unmarshal([]byte(msg), &test)
		fmt.Println(test)

		server.BroadcastToRoom(test.Room, "message", test.Message)

	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	e.GET("/socket.io/", echo.WrapHandler(server))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1333"))
}
