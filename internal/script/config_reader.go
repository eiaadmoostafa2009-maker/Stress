package script

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"os"
	"stress/internal/config"
	"time"
)

func LoadConfig(path string) (config.ScriptConfig, error) {
	// Read the file before starting any workload.
	if _, err := os.Stat(path); err != nil {
		return config.ScriptConfig{}, fmt.Errorf("scenario file: %w", err)
	}

	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile(path); err != nil {
		return config.ScriptConfig{}, fmt.Errorf("load Lua scenario: %w", err)
	}

	value := L.GetGlobal("config")
	table, ok := value.(*lua.LTable)
	if !ok {
		return config.ScriptConfig{}, fmt.Errorf("script must define a config table")
	}

	cfg := config.ScriptConfig{
		Name:      lua.LVAsString(table.RawGetString("name")),
		BaseURL:   lua.LVAsString(table.RawGetString("base_url")),
		Users:     int(lua.LVAsNumber(table.RawGetString("users"))),
		Duration:  time.Duration(lua.LVAsNumber(table.RawGetString("duration_seconds"))) * time.Second,
		Timeout:   time.Duration(lua.LVAsNumber(table.RawGetString("timeout_seconds"))) * time.Second,
		ThinkTime: time.Duration(lua.LVAsNumber(table.RawGetString("think_time_seconds"))) * time.Second,
	}

	if err := cfg.ScriptValidate(); err != nil {
		return config.ScriptConfig{}, err
	}

	if L.GetGlobal("script").Type() != lua.LTFunction {
		return config.ScriptConfig{}, fmt.Errorf("script must define function script(user)")
	}

	return cfg, nil
}
