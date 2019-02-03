package jsonrpc

import (
	"encoding/json"
	"fmt"
	"github.com/biety/base"
	"github.com/biety/config"
	"io/ioutil"
	"net/http"
	"strconv"
)

var server RPCServer

type RPCServer struct {
	m     map[string]func([]interface{}) map[string]interface{}
}

func StartRPCServer() error {
	server.m = make(map[string]func([]interface{}) map[string]interface{})
	http.HandleFunc("/", Handle)

	HandleFunc("getversion", GetNodeVersion)

	err := http.ListenAndServe(":"+strconv.Itoa(config.Rpc_port), nil)
	if err != nil {
		return fmt.Errorf("ListenAndServe error : %s\n", err)
	}
	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", "application/json;charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	if r.Method != "POST" {
		fmt.Printf("HTTP JSON RCP handle - Method != POST\n")
		return
	}

	if r.Body == nil {
		fmt.Printf("HTTP JSON RCP handle - Request body is nil\n")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("HTTP JSON RCP handle - ioutil.ReadAll: %s\n", err)
		return
	}

	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Printf("HTTP JSON RCP handle - json.Unmarshal: %s\n", err)
		return
	}

	if request["method"] == nil {
		fmt.Printf("HTTP JSON RCP handle - method not found: %s\n", err)
		return
	}

	method, ok := request["method"].(string)
	if !ok {
		fmt.Printf("HTTP JSON RCP handle - method is not string: \n")
		return
	}

	function, ok := server.m[method]
	if ok {
		response := function(request["params"].([]interface{}))
		data, err := json.Marshal(map[string]interface{} {
			"jsonrpc": "2.0",
			"error": response["error"],
			"desc": response["desc"],
			"result": response["result"],
			"id": request["id"],
		})

		if err != nil {
			fmt.Printf("HTTP JSON RCP handle - json.marshal: %s\n", err)
			return
		}

		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", "application/json;charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	} else {
		fmt.Printf("HTTP JSON RPC Handle - No function to call for %s \n", method)
		data, err := json.Marshal(map[string]interface{}{
			"error": base.INVALID_METHOD,
			"result": map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
				"data":    "The called method was not found on the server",
			},
			"id": request["id"],
		})

		if err != nil {
			fmt.Printf("HTTP JSON RPC Handle - json.Marshal: %s\n", err)
			return
		}
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", "application/json;charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	}
}

func HandleFunc(pattern string, handler func([]interface {}) map[string]interface{}) {
	server.m[pattern] = handler
}

func GetNodeVersion(params []interface{}) map[string]interface{} {
	return responseSuccess("1.0")
}

func responseSuccess(result interface{}) map[string]interface{} {
	return responsePack(0, result)
}

func responsePack(errcode int64, result interface{}) map[string]interface{} {
	resp := map[string]interface{} {
		"error": errcode,
		"desc": "SUCCESS",
		"result": result,
	}

	return resp
}