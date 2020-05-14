package dbconn

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strings"
)

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...driver.NamedValue) (context.Context, error) {
	for i := range args {
		arg := args[i].Value
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.String {
			args[i].Value = strings.TrimSpace(val.String())
		}
	}
	return ctx, nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...driver.NamedValue) (context.Context, error) {
	return ctx, nil
}
