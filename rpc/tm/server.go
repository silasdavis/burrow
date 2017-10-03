// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tm

import (
	"net"
	"net/http"

	"github.com/hyperledger/burrow/logging/structure"
	"github.com/hyperledger/burrow/rpc/core"
	"github.com/tendermint/tendermint/rpc/lib/server"
	"github.com/tendermint/tmlibs/events"
	"github.com/tendermint/tmlibs/log"
)

func StartServer(service core.Service, pattern, listenAddress string,
	evsw events.EventSwitch, logger log.Logger) (net.Listener, error) {

	logger = logger.With(structure.ComponentKey, "rpc/tm")
	routes := GetRoutes(service)
	mux := http.NewServeMux()
	wm := rpcserver.NewWebsocketManager(routes, evsw)
	mux.HandleFunc(pattern, wm.WebsocketHandler)
	rpcserver.RegisterRPCFuncs(mux, routes, logger)
	listener, err := rpcserver.StartHTTPServer(listenAddress, mux, logger)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
