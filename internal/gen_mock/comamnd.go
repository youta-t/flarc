// Code generated -- DO NOT EDIT
package gen_mock

import (
	its "github.com/youta-t/its"
	itskit "github.com/youta-t/its/itskit"
	mockkit "github.com/youta-t/its/mocker/mockkit"
	pkg1 "context"
	pkg2 "github.com/youta-t/flarc"
	
)

type _TaskCallSpec[T any] struct {
	arg0 its.Matcher[pkg1.Context]
	
	arg1 its.Matcher[pkg2.Commandline[T]]
	
	arg2 its.Matcher[[]any]
	
	
}

type _TaskCall[T any] struct {
	name itskit.Label
	spec _TaskCallSpec[T]
}

func Task_Expects[T any](
	arg0 its.Matcher[pkg1.Context],
	
	arg1 its.Matcher[pkg2.Commandline[T]],
	
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

type _TaskBehavior [T any] struct {
	name itskit.Label
	spec _TaskCallSpec[T]
	effect func(arg0 pkg1.Context, arg1 pkg2.Commandline[T], arg2 []any) error
}

func (b *_TaskBehavior[T]) Fn(t mockkit.TestLike) func(arg0 pkg1.Context, arg1 pkg2.Commandline[T], arg2 []any) error {
	return func (
		
		arg0 pkg1.Context,
		
		arg1 pkg2.Commandline[T],
		
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
				matcher = its.Never[pkg1.Context]()
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
				matcher = its.Never[pkg2.Commandline[T]]()
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

) mockkit.FuncBehavior[ func (arg0 pkg1.Context, arg1 pkg2.Commandline[T], arg2 []any) error  ] {
	return c.ThenEffect(func(
		
		pkg1.Context,
		
		pkg2.Commandline[T],
		
		[]any,
		
		
	)(
		error,
		
	){
		
		return ret0
		
	})
}

func (c _TaskCall[T]) ThenEffect(effect func(arg0 pkg1.Context, arg1 pkg2.Commandline[T], arg2 []any) error) mockkit.FuncBehavior[ func (arg0 pkg1.Context, arg1 pkg2.Commandline[T], arg2 []any) error ] {
	return &_TaskBehavior[T] {
		name: c.name,
		spec: c.spec,
		effect: effect,
	}
}




