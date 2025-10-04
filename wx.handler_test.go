package wx

import (
	"errors"
	"reflect"
	"testing"
)

type Context[TUser any] struct {
	User TUser
}

var cache = map[any]any{}

func AddFn[TController any, TUser any](c *TController, fn func(ctx *Context[TUser]) error) {
	cache[*c] = fn
}
func Call[TController any, TUser any](c *TController, ctx *Context[TUser]) (*TUser, error) {
	fn, ok := cache[*c]
	if !ok {
		return nil, errors.New("loi")
	}

	if fx, ok := fn.(func(ctx *Context[TUser]) error); ok {
		err := fx(ctx)
		if err != nil {
			return nil, err
		}
		return &ctx.User, nil
	}
	return nil, errors.New("loi")
}

type Controller struct {
}
type UserTest struct {
}

var cachFn = map[any]any{}

func NewInstance[TController any, TUser any](uri string) func(user TUser) {
	c := reflect.New(reflect.TypeFor[TController]()).Interface().(*TController)
	AddFn(c, func(ctx *Context[TUser]) error {

		return nil
	})

	cachFn[*c] = func(user TUser) {
		ctx := &Context[TUser]{}
		Call(c, ctx)
	}
	return func(user TUser) {
		fn := cachFn[*c]
		fn.(func(user TUser))(user)
	}

}
func TestCall(t *testing.T) {
	NewInstance[Controller, UserTest]("test-001")
	fx := NewInstance[Controller, UserTest]("test-002")
	//fx := cachFn["test-001"]
	fx(UserTest{})

}
