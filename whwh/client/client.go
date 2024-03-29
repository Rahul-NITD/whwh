package client

import (
	"encoding/json"
	"net/http"

	"github.com/aargeee/whwh/whwh"
)

func RequestChannel(requestUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, requestUrl, http.NoBody)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var requestChannelResponse whwh.CreateChannelResponse
	if err := json.NewDecoder(res.Body).Decode(&requestChannelResponse); err != nil {
		return "", err
	}
	return requestChannelResponse.Response.ChannelID, nil
}
