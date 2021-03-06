package module

import (
	. "core/libs"
	"core/libs/sessions"
	"core/protos"
	"core/protos/gameProto"
	"github.com/golang/protobuf/proto"
	"servives/chat/cache"
	"servives/public"
	"servives/public/errCodes"
	"servives/public/redisCaches"
)

//获取用户信息
func JoinChat(clientSession *sessions.BackSession, msgData proto.Message) {
	data := msgData.(*gameProto.UserJoinChatC2S)
	token := data.GetToken()
	userId := public.GetUserIdByToken(token)
	if userId == 0 {
		public.SendErrorMsgToClient(clientSession, errCodes.PARAM_ERROR)
		return
	}

	//获取redis缓存中用户数据
	dbUser := redisCaches.GetUser(userId)
	if dbUser == nil {
		public.SendErrorMsgToClient(clientSession, errCodes.PARAM_ERROR)
		return
	}

	//保存到内存中
	clientSession.SetUserId(dbUser.Id)
	cache.AddUser(dbUser.Id, dbUser.Account, clientSession)

	//用户下线处理
	clientSession.AddCloseCallback(nil, "user.joinChatSuccess", func() {
		cache.RemoveUser(dbUser.Id)
		DEBUG("用户下线：当前在线人数", cache.GetOnlineUsersNum(), sessions.BackSessionLen())
	})
	DEBUG("用户上线：当前在线人数", cache.GetOnlineUsersNum(), sessions.BackSessionLen())

	//返回客户端
	sendMsg := &gameProto.UserJoinChatS2C{}
	public.SendMsgToClient(clientSession, sendMsg)
}

func Chat(clientSession *sessions.BackSession, msgData proto.Message) {
	data := msgData.(*gameProto.UserChatC2S)
	chatUser := cache.GetUser(clientSession.UserID())
	if chatUser == nil {
		public.SendErrorMsgToClient(clientSession, errCodes.PARAM_ERROR)
		return
	}

	//发送给所有人
	sendMsg := &gameProto.UserChatNoticeS2C{
		UserId:   protos.Uint64(chatUser.UserID),
		UserName: protos.String(chatUser.UserName),
		Msg:      protos.String(data.GetMsg()),
	}
	public.SendMsgToAllClient(sendMsg)
}
