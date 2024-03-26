package whwh

type Doc struct {
	Name           string
	Endpoint       string
	Description    string
	Examples       string
	AllowedMethods []string
	Response       interface{}
}

var CreateChannelDoc = Doc{
	Name:           "Create Channel",
	Endpoint:       "/create",
	AllowedMethods: []string{"GET", "POST"},
	Description:    "Create a new channel by making a post request on /create",
	Examples:       "",
	Response:       ChannelIDResponse{},
}
