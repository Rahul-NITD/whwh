package whwh

type ChannelIDResponse struct {
	ChannelID string `json:"channelid"`
}

type CreateChannelResponse struct {
	Event    string            `json:"event"`
	Response ChannelIDResponse `json:"response"`
	Message  string            `json:"message"`
	Status   string            `json:"success"`
}
