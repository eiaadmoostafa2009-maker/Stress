package script

import (
	"context"
	"io"
	"net/http"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type HttpApi struct {
	client  *http.Client
	baseUrl string
	context context.Context
}

func (api HttpApi) Register(L *lua.LState) {
	user := L.NewTable()

	L.SetField(user, "get", L.NewFunction(
		func(L *lua.LState) int {
			return api.request(L, http.MethodGet)
		},
	))

	L.SetField(user, "post", L.NewFunction(
		func(L *lua.LState) int {
			return api.request(L, http.MethodPost)
		},
	))

	L.SetGlobal("user", user)
}

func (api *HttpApi) request(L *lua.LState, method string) int {
	route := L.CheckString(2)

	if !strings.HasPrefix("/", route) {
		L.RaiseError("the route should starts with /")
		return 0
	}

	req, err := http.NewRequestWithContext(api.context, method, api.baseUrl+route, nil)
	if err != nil {
		L.RaiseError("Error: %v", err)
		return 0
	}

	resp, err := api.client.Do(req)
	if err != nil {
		L.RaiseError("Error: %v", err)
		return 0
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		L.RaiseError("Response body error: %v", err)
		return 0
	}

	result := L.NewTable()
	result.RawSetString(
		"status",
		lua.LNumber(resp.StatusCode),
	)

	result.RawSetString(
		"body",
		lua.LString(string(body)),
	)

	L.Push(result)
	return 1
}
