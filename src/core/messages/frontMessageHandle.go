package messages

import (
	"core"
	"core/consts/errCode"
	"core/consts/service"
	. "core/libs"
	"core/libs/grpc/ipc"
	"core/libs/sessions"
	"core/protos"
	"core/protos/gameProto"
	"errors"
)

func dealConnectorMsg(session *sessions.FrontSession, msgBody []byte) {
	protoMsg := protos.UnmarshalProtoMsg(msgBody)
	if protoMsg == protos.NullProtoMsg {
		return
	}

	//Ping消息特殊处理
	if protoMsg.ID == gameProto.ID_client_ping_c2s {
		session.UpdatePingTime()
		return
	}
}

func dealGameMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToIpcService(Service.Game, session, msgBody)
	if err != nil {
		ERR("dealGameMsg", err)
		sendErrorMsgToClient(session)
	}
}

func dealLoginMsg(session *sessions.FrontSession, msgBody []byte) {
	err := sendMsgToIpcService(Service.Login, session, msgBody)
	if err != nil {
		ERR("dealLoginMsg", err)
		sendErrorMsgToClient(session)
	}
}

func getGameService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	//1: 获取用户数据，根据Token分配
	msgId := protos.UnmarshalProtoId(msgBody)
	if msgId == gameProto.ID_user_getInfo_c2s {
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			return ""
		}
		protoMsgData := protoMsg.Body.(*gameProto.UserGetInfoC2S)
		return ipcClient.GetServiceByFlag(protoMsgData.GetToken())
	} else {
		return session.GetIpcService(Service.Game)
	}
}

func getLoginService(session *sessions.FrontSession, msgBody []byte, ipcClient *ipc.Client) string {
	//1: 登录，根据Account分配
	msgId := protos.UnmarshalProtoId(msgBody)
	if msgId == gameProto.ID_user_login_c2s {
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			return ""
		}
		protoMsgData := protoMsg.Body.(*gameProto.UserLoginC2S)
		return ipcClient.GetServiceByFlag(protoMsgData.GetAccount())
	} else {
		return session.GetIpcService(Service.Login)
	}
}

func sendErrorMsgToClient(session *sessions.FrontSession) {
	sendMsg := protos.MarshalProtoMsg(&gameProto.ErrorNoticeS2C{
		ErrorCode: protos.Int32(ErrCode.SYSTEM_ERR),
	})
	session.Send(sendMsg)
}

func sendMsgToIpcService(serviceName string, clientSession *sessions.FrontSession, msgBody []byte) error {
	ipcClient := core.Service.GetIpcClient(serviceName)
	if ipcClient == nil {
		return errors.New(serviceName + ": ipcClient not exists")
	}

	var service string
	if serviceName == Service.Login {
		service = getLoginService(clientSession, msgBody, ipcClient)
	} else if serviceName == Service.Game {
		service = getGameService(clientSession, msgBody, ipcClient)
	}

	if service == "" {
		return errors.New(serviceName + ": service not exists")
	}

	err := ipcClient.Send(core.Service.Identify(), clientSession.ID(), msgBody, service)
	if err == nil {
		clientSession.SetIpcService(serviceName, service)
	}
	return err
}
