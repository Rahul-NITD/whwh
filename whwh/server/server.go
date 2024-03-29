package server

import (
	"github.com/aargeee/whwh/whwh"
	"github.com/google/uuid"
)

func CreateNewChannel() whwh.CreateChannelResponse {
	sid := uuid.NewString()
	return whwh.CreateChannelResponse{
		Event:    "CreateChannel",
		Message:  "Channel Created Successfully",
		Status:   "SUCCESS",
		Response: whwh.ChannelIDResponse{ChannelID: sid},
	}
}
