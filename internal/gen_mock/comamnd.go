// Code generated -- DO NOT EDIT
package gen_mock

import (
	its "github.com/youta-t/its"
	itskit "github.com/youta-t/its/itskit"
	testee "github.com/youta-t/flarc"
	u_context "context"
	
)

type _TaskReturnFixture[T any] struct {
	ret0 error
	
}

type _TaskReturn[T any] struct {
	fixture _TaskReturnFixture[T]
}

func (rfx _TaskReturn[T]) Get() (
	error,
	
) {
	return rfx.fixture.ret0
}

type _TaskCallSpec[T any] struct {
	arg0 its.Matcher[u_context.Context]
	
	arg1 its.Matcher[testee.Commandline[T]]
	
	arg2 its.Matcher[[]any]
	
	
}

type _TaskCall[T any] struct {
	name itskit.Label
	spec _TaskCallSpec[T]
}

func NewTaskCall[T any](
	arg0 its.Matcher[u_context.Context],
	
	arg1 its.Matcher[testee.Commandline[T]],
	
	arg2 its.Matcher[[]any],
	
) _TaskCall[T] {
	cancel := itskit.SkipStack()
	defer cancel()

	spec := _TaskCallSpec[T] {}
	spec.arg0 = itskit.Named(
		"arg0",
		arg0,
	)
	
	spec.arg1 = itskit.Named(
		"arg1",
		arg1,
	)
	
	spec.arg2 = itskit.Named(
		"arg2",
		arg2,
	)
	
	
	return _TaskCall[T]{
		name: itskit.NewLabelWithLocation("func Task"),
		spec: spec,
	}
}

type TaskBehaviour [T any] struct {
	name itskit.Label
	spec _TaskCallSpec[T]
	effect func(arg0 u_context.Context, arg1 testee.Commandline[T], arg2 []any) error
}

func (b *TaskBehaviour[T]) Mock(t interface { Error(...any) }) func(arg0 u_context.Context, arg1 testee.Commandline[T], arg2 []any) error {
	return func (
		
		arg0 u_context.Context,
		
		arg1 testee.Commandline[T],
		
		arg2 []any,
		
		
	) (
		error,
		
	) {
		if h, ok := t.(interface { Helper() }); ok {
			h.Helper()
		}
		ok := 0
		matches := []itskit.Match{}
		
		{
			matcher := b.spec.arg0
			if matcher == nil {
				matcher = its.Never[u_context.Context]()
			}
			m := matcher.Match(arg0)
			if m.Ok() {
				ok += 1
			}
			matches = append(matches, m)
		}
		
		{
			matcher := b.spec.arg1
			if matcher == nil {
				matcher = its.Never[testee.Commandline[T]]()
			}
			m := matcher.Match(arg1)
			if m.Ok() {
				ok += 1
			}
			matches = append(matches, m)
		}
		
		{
			matcher := b.spec.arg2
			if matcher == nil {
				matcher = its.Never[[]any]()
			}
			m := matcher.Match(arg2)
			if m.Ok() {
				ok += 1
			}
			matches = append(matches, m)
		}
		
		itskit.NewMatch(
			ok == len(matches),
			b.name.Fill(itskit.Missing),
			matches...,
		).OrError(t)
		return b.effect(
			
			arg0,
			
			arg1,
			
			arg2,
			
			
		)
	}
}

func (c _TaskCall[T]) ThenReturn(

	ret0 error,

) *TaskBehaviour[T] {
	return c.ThenEffect(func(
		
		u_context.Context,
		
		testee.Commandline[T],
		
		[]any,
		
		
	)(
		error,
		
	){
		
		return ret0
		
	})
}

func (c _TaskCall[T]) ThenEffect(effect func(arg0 u_context.Context, arg1 testee.Commandline[T], arg2 []any) error) *TaskBehaviour[T] {
	return &TaskBehaviour[T] {
		name: c.name,
		spec: c.spec,
		effect: effect,
	}
}




