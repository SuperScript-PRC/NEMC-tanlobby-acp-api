package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	I18n "github.com/Happy2018new/nemc-tan-lobby-solver/bunker/i18n"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login"
	"github.com/Happy2018new/nemc-tan-lobby-solver/ws_api"
)

var auth_server_url = flag.String("A", "http://47.101.71.192:24990", "Auth Service URL")
var ws_server_port = flag.Int("port", 25010, "WebSocket API port")
var room_id = flag.String("R", "", "Room ID")
var room_passcode = flag.String("P", "", "Room Password")
var token = flag.String("T", "", "Your token")

func main() {
	flag.Parse()

	client, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: *auth_server_url,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(I18n.T(I18n.ACP_ContactingWithAuthServer))
	wrapper := auth.NewAccessWrapper(client, *room_id, *room_passcode, *token)
	netConn, err := login.Dial(wrapper)
	if err != nil {
		panic(err)
	}

	fmt.Println(I18n.T(I18n.ACP_ConnectingToGame))
	conn, err := minecraft.DialContext(context.Background(), netConn)
	if err != nil {
		panic(err)
	}

	fmt.Println(I18n.T(I18n.ConnectionEstablished))

	go ws_api.StartWSServer(*ws_server_port)
	go ws_api.HandleEndpointsMessages(conn.WritePacket)

	fmt.Println(I18n.T(I18n.ACP_ApiReady))
	go ws_api.SetReady()

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}
		go ws_api.HandleServerPacketsToEndpoints(conn, pk)
	}
}
