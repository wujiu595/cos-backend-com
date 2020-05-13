package dbconn

import (
	"context"
	"reflect"
)

// Hooks satisfies the sqlhook.Hooks interface
//type Hooks interface {
//				Before(ctx context.Context, query string, args ...interface{}) (context.Context, error)
//				After(ctx context.Context, query string, args ...interface{}) (context.Context, error)
//}

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	for i := range args {
		arg := reflect.ValueOf(args[i]).FieldByName("Value")
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.String {
			arg.Elem().Set(val)
		}
	}
	return ctx, nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return ctx, nil
}
