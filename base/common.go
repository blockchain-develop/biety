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

const (
	SYNC_MAX_HEADER_FORWARD_SIZE = 5000       //keep CurrentHeaderHeight - CurrentBlockHeight <= SYNC_MAX_HEADER_FORWARD_SIZE
	SYNC_MAX_FLIGHT_HEADER_SIZE  = 1          //Number of headers on flight
	SYNC_MAX_FLIGHT_BLOCK_SIZE   = 50         //Number of blocks on flight
	SYNC_MAX_BLOCK_CACHE_SIZE    = 500        //Cache size of block wait to commit to ledger
	SYNC_HEADER_REQUEST_TIMEOUT  = 2          //s, Request header timeout time. If header haven't receive after SYNC_HEADER_REQUEST_TIMEOUT second, retry
	SYNC_BLOCK_REQUEST_TIMEOUT   = 2          //s, Request block timeout time. If block haven't received after SYNC_BLOCK_REQUEST_TIMEOUT second, retry
	SYNC_NEXT_BLOCK_TIMES        = 3          //Request times of next height block
	SYNC_NEXT_BLOCKS_HEIGHT      = 2          //for current block height plus next
	SYNC_NODE_RECORD_SPEED_CNT   = 3          //Record speed count for accuracy
	SYNC_NODE_RECORD_TIME_CNT    = 3          //Record request time  for accuracy
	SYNC_NODE_SPEED_INIT         = 100 * 1024 //Init a big speed (100MB/s) for every node in first round
	SYNC_MAX_ERROR_RESP_TIMES    = 5          //Max error headers/blocks response times, if reaches, delete it
	SYNC_MAX_HEIGHT_OFFSET       = 5          //Offset of the max height and current height
)

// multi-sig constants
const MULTI_SIG_MAX_PUBKEY_SIZE = 16

// transaction constants
const TX_MAX_SIG_SIZE = 16

const (
	MAX_CAPACITY     = 100140                           // The tx pool's capacity that holds the verified txs
	MAX_PENDING_TXN  = 4096 * 10                        // The max length of pending txs
	MAX_WORKER_NUM   = 2                                // The max concurrent workers
	MAX_RCV_TXN_LEN  = MAX_WORKER_NUM * MAX_PENDING_TXN // The max length of the queue that server can hold
	MAX_RETRIES      = 0                                // The retry times to verify tx
	EXPIRE_INTERVAL  = 9                                // The timeout that verify tx
	STATELESS_MASK   = 0x1                              // The mask of stateless validator
	STATEFUL_MASK    = 0x2                              // The mask of stateful validator
	VERIFY_MASK      = STATELESS_MASK | STATEFUL_MASK   // The mask that indicates tx valid
	MAX_LIMITATION   = 10000                            // The length of pending tx from net and http
	UPDATE_FREQUENCY = 100                              // The frequency to update gas price from global params
	MAX_TX_SIZE      = 1024 * 1024                      // The max size of a transaction to prevent DOS attacks
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
