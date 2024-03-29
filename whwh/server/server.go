package server

import (
	"github.com/aargeee/whwh/whwh"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

var GlobalSSEServer = sse.New()

func CreateNewChannel() whwh.CreateChannelResponse {
	sid := uuid.NewString()
	for GlobalSSEServer.StreamExists(sid) {
		sid = uuid.NewString()
	}

	GlobalSSEServer.CreateStream(sid)

	return whwh.CreateChannelResponse{
		Event:    "CreateChannel",
		Message:  "Channel Created Successfully",
		Status:   "SUCCESS",
		Response: whwh.ChannelIDResponse{ChannelID: sid},
	}
}
