package script

import (
	"context"
	"fmt"
	"net/http"
	"stress/internal/config"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func RunVirtualUser(ctx context.Context, cfg config.ScriptConfig, path string) error {
	L := lua.NewState()
	defer L.Close()

	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	api := &HttpApi{
		client:  client,
		baseUrl: cfg.BaseURL,
		context: ctx,
	}
	api.Register(L)

	if err := L.DoFile(path); err != nil {
		return fmt.Errorf("lua script: %w", err)
	}

	fn := L.GetGlobal("script")
	if fn.Type() != lua.LTFunction {
		return fmt.Errorf("missing script(user) function")
	}

	user := L.GetGlobal("user")

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		err := L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, user)

		if err != nil {
			return fmt.Errorf("execute script: %w", err)
		}

		if cfg.ThinkTime > 0 {
			timer := time.NewTimer(cfg.ThinkTime)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
					return nil
				}
			case <-timer.C:
			}

		}

	}
}
