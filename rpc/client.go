package rpc

type Client interface {
	Call(method string, params map[string]interface{}, result interface{}) (interface{}, error)
}
