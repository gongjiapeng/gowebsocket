package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"gowebsocket/impl"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func wsHander(w http.ResponseWriter,r *http.Request) {
	var (
		wsConn *websocket.Conn
		err error
		conn *impl.Connection
		data []byte
	)
	if wsConn,err = upgrader.Upgrade(w,r,nil);err!=nil{
		return
	}
	if conn,err = impl.InitConnection(wsConn);err!=nil{
		goto ERR
	}
	//心跳检查
	go func() {
		for  {
			var err error
			if err = conn.WriteMessage([]byte("heartcheck"));err!=nil{
				return
			}
		}
	}()
	for  {
		if data ,err = conn.ReadMessage();err!=nil{
			goto ERR
		}
		if err = conn.WriteMessage(data);err!=nil {
			goto ERR
		}
	}
	ERR:
		//todo guanbi conn
		conn.Close()
}
func main() {
	http.HandleFunc("/ws",wsHander)
	http.ListenAndServe("0.0.0.0:8000",nil)
}