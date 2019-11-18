package handler

import (
	"io"

	"github.com/zhiguochi/chess/pb/center"
	"github.com/zhiguochi/chess/util/buf_pool"
	"github.com/zhiguochi/chess/util/log"
	"github.com/zhiguochi/chess/util/rpc"
	"github.com/golang/protobuf/proto"
)

func sendNewConnInfoNotify(info *center.ConnInfo, excludeClient io.Writer) {
	notify := &center.NewConnInfoNotify{info}

	writer := buf_pool.Get()
	defer buf_pool.Put(writer)

	err := rpc.EncodePb(writer, notify)
	if err != nil {
		log.Warn("encode protobuf fail:%s", err.Error())
		return
	}

	sendClientNotify(writer.Bytes(), excludeClient)
}

func sendDelConnInfoNotify(info *center.ConnInfo, excludeClient io.Writer) {
	notify := &center.DelConnInfoNotify{info}

	writer := buf_pool.Get()
	defer buf_pool.Put(writer)

	err := rpc.EncodePb(writer, notify)
	if err != nil {
		log.Warn("encode protobuf fail:%s", err.Error())
		return
	}

	sendClientNotify(writer.Bytes(), excludeClient)
}

func sendDelConnInfoByGateidNotify(gateid uint32, excludeClient io.Writer) {
	notify := &center.DelConnInfoByGateidNotify{gateid}

	writer := buf_pool.Get()
	defer buf_pool.Put(writer)

	err := rpc.EncodePb(writer, notify)
	if err != nil {
		log.Warn("encode protobuf fail:%s", err.Error())
		return
	}

	sendClientNotify(writer.Bytes(), excludeClient)
}

func sendResp(client io.Writer, resp proto.Message) error {
	writer := buf_pool.Get()
	defer buf_pool.Put(writer)

	if err := rpc.EncodePb(writer, resp); err != nil {
		return err
	}

	_, err := client.Write(writer.Bytes())
	return err
}
