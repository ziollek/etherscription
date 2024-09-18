package etherum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ziollek/etherscription/pkg/config"
	"github.com/ziollek/etherscription/pkg/logging"
	"github.com/ziollek/etherscription/pkg/model"
)

const (
	createFilter     = "eth_newFilter"
	getFilterChanges = "eth_getFilterChanges"
	getTransaction   = "eth_getTransactionByHash"
	rpcID            = 111
	rpcVersion       = "2.0"
	startBlock       = "latest"
)

type createFilterRequest struct {
	FromBlock string `json:"fromBlock,omitempty"`
	ToBlock   string `json:"toBlock,omitempty"`
	Address   string `json:"address,omitempty"`
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      int         `json:"id"`
	Params  interface{} `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (m rpcError) ToError() error {
	return fmt.Errorf("RPC error code: %d,  %s", m.Code, m.Message)
}

type jsonRPCResponse[T any] struct {
	JSONRPC string    `json:"jsonrpc"`
	Method  string    `json:"method"`
	Result  T         `json:"result"`
	Error   *rpcError `json:"error"`
}

func (response *jsonRPCResponse[T]) ToResponse() (T, error) {
	if response.Error != nil {
		return response.Result, response.Error.ToError()
	}
	return response.Result, nil
}

func newEncodedJSONRPCRequest(method string, params interface{}) ([]byte, error) {
	return json.Marshal(jsonRPCRequest{
		JSONRPC: rpcVersion,
		Method:  method,
		Params:  params,
		ID:      rpcID,
	})
}

type RPCClient struct {
	node          string
	httpClient    *http.Client
	slowDownDelay time.Duration
}

func NewRPCClient(node string, cfg *config.RPCConfig) *RPCClient {
	return &RPCClient{
		node: node,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		slowDownDelay: cfg.TooManyRequestsDelay,
	}
}

func (c *RPCClient) createFilter() (string, error) {
	result := jsonRPCResponse[string]{}
	if err := c.makeCall(createFilter, []createFilterRequest{{FromBlock: startBlock}}, &result); err != nil {
		return "", err
	}
	return result.ToResponse()
}

func (c *RPCClient) getChanges(filterID string) (logEntries, error) {
	result := jsonRPCResponse[logEntries]{}
	if err := c.makeCall(getFilterChanges, []string{filterID}, &result); err != nil {
		return nil, err
	}
	return result.ToResponse()
}

func (c *RPCClient) getTransaction(hash string) (*model.RawTransaction, error) {
	result := jsonRPCResponse[*model.RawTransaction]{}
	if err := c.makeCall(getTransaction, []string{hash}, &result); err != nil {
		return nil, fmt.Errorf("error while read body from HTTP response: %w", err)
	}
	return result.ToResponse()
}

func (c *RPCClient) makeCall(method string, input, output interface{}) error {
	payload, err := newEncodedJSONRPCRequest(method, input)
	logging.Logger().Debug().Str("payload", string(payload)).Msgf("payload for %s", method)
	if err != nil {
		return fmt.Errorf("error while encoding JSON RPC request: %w", err)
	}
	request, err := http.NewRequest("POST", c.node, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error while creating HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error while making HTTP request: %w", err)
	}
	if resp.StatusCode == 429 {
		logging.Logger().Warn().Str("response", resp.Status).Msgf("rate limit reached, trying slow down")
		time.Sleep(c.slowDownDelay)
		resp.Body.Close()
		request, err = http.NewRequest("POST", c.node, bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("error while creating HTTP request: %w", err)
		}
		resp, err = c.httpClient.Do(request)
		if err != nil {
			return fmt.Errorf("error while making HTTP request: %w", err)
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error while making HTTP request: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while read body from HTTP response: %w", err)
	}
	logging.Logger().Debug().Str("response", string(body)).Msgf("response fetched for %s", createFilter)
	err = json.Unmarshal(body, &output)
	if err != nil {
		return fmt.Errorf("error while unmarshall body from HTTP response: %w", err)
	}
	return nil
}
