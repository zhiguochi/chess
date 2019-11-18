package handler

import (
	"github.com/zhiguochi/chess/common"
	"github.com/zhiguochi/chess/game/server"
)

func HandleEcho(userid uint32, connid uint32, msgBody []byte) {
	server.SendResp(userid, MsgidEchoResp, common.ResultSuccess, msgBody)
}
