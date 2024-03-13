package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aargeee/whwh/whwh"
)

func NewServer() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		res := whwh.CreateChannelResponse{
			Event:    "CreateChannel",
			Message:  "Channel Created Successfully",
			Status:   "SUCCESS",
			Response: whwh.ChannelIDResponse{ChannelID: "KJFNGAB7DFGDSGF7GFS7GF8S7"},
		}
		json.NewEncoder(w).Encode(res)
		w.WriteHeader(http.StatusOK)
	})
	return r
}
