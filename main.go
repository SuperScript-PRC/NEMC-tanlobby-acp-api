package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/Happy2018new/nemc-tan-lobby-solver/bunker/auth"
	"github.com/Happy2018new/nemc-tan-lobby-solver/minecraft"
	"github.com/Happy2018new/nemc-tan-lobby-solver/protocol/login"
	"github.com/Happy2018new/nemc-tan-lobby-solver/ws_api"
)

var auth_server_url = flag.String("A", "http://47.101.71.192:24990", "验证服务 URL")
var ws_server_port = flag.Int("port", 25010, "ToolDelta WebSocket API 服务端口")
var room_id = flag.String("R", "", "房间号")
var room_passcode = flag.String("P", "", "房间密码")
var token = flag.String("T", "", "验证服务器账号的 token")

func main() {
	flag.Parse()
	fmt.Println("正在从验证服务器取得信息")
	client, err := auth.CreateClient(&auth.ClientOptions{
		AuthServer: *auth_server_url,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("正在登录到房间")
	wrapper := auth.NewAccessWrapper(client, *room_id, *room_passcode, *token)
	netConn, err := login.Dial(wrapper)
	if err != nil {
		panic(err)
	}

	fmt.Println("正在登录到我的世界游戏网络")
	conn, err := minecraft.DialContext(context.Background(), netConn)
	if err != nil {
		panic(err)
	}

	fmt.Println("本地联机接入点已就绪")

	go ws_api.StartWSServer(*ws_server_port)
	go ws_api.HandleEndpointsMessages(conn.WritePacket)

	fmt.Println("WebSocket API 已就绪")
	go ws_api.SetReady()

	println("Starting to read packets...")
	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}
		go ws_api.HandleServerPacketsToEndpoints(conn, pk)
	}
}
