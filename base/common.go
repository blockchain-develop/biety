package base

import (
	"errors"
	"strconv"
	"strings"
)

const (
	PROTOCOL_VERSION      = 0     //protocol version
	UPDATE_RATE_PER_BLOCK = 2     //info update rate in one generate block period
	KEEPALIVE_TIMEOUT     = 15    //contact timeout in sec
	DIAL_TIMEOUT          = 6     //connect timeout in sec
	CONN_MONITOR          = 6     //time to retry connect in sec
	CONN_MAX_BACK         = 4000  //max backoff time in micro sec
	MAX_RETRY_COUNT       = 3     //max reconnect time of remote peer
	CHAN_CAPABILITY       = 10000 //channel capability of recv link
	SYNC_BLK_WAIT         = 2     //timespan for blk sync check
)




func ParseIPAddr(s string) (string, error) {
	i := strings.Index(s, ":")
	if i < 0 {
		return "",errors.New("split ip address error")
	}
	return s[:i], nil
}

func ParseIPPort(s string)(string, error) {
	i := strings.Index(s, ":")
	if i < 0 {
		return "",errors.New("split ip port error")
	}

	port, err := strconv.Atoi(s[i+1:])
	if err != nil {
		return "",errors.New("parse port error")
	}

	if port <= 0 || port >= 65536 {
		return "",errors.New("port out of bound")
	}

	return s[i:], nil
}
