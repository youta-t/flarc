package parser

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/youta-t/flarc/flarcerror"
	"github.com/youta-t/flarc/params"
)

type Parser[T any] interface {
	Parse([]string) (
		flags *T,
		args map[string][]string,
		rem []string,
		err error,
	)

	String() string

	Flags() []params.Flag
	Args() []params.Arg
}

func New[T any](
	flagdef *T, pos []params.ArgDef,
) (Parser[T], error) {
	refd := reflect.ValueOf(flagdef)
	trefd := refd.Elem().Type()

	_pos := make([]params.Arg, len(pos))
	for i := range pos {
		_pos[i] = pos[i].Freeze()
	}

	psr := &parser[T]{
		dest: flagdef,
		args: _pos,
	}

	names := map[string]struct{}{}
	for i := 0; i < trefd.NumField(); i += 1 {
		ref := trefd.Field(i)

		flg, err := params.NewFlag(ref, refd.Elem().Field(i))
		if err != nil {
			return nil, err
		}

		if _, ok := names[flg.Name()]; ok {
			return nil, errors.New("")
		}

		psr.flags = append(psr.flags, flg)
	}

	return psr, nil
}

type parser[T any] struct {
	dest  *T
	flags []params.Flag
	args  []params.Arg
}

func (p *parser[T]) String() string {
	items := []string{}
	for _, f := range p.flags {
		items = append(items, f.Usage())
	}

	for _, pos := range p.args {
		items = append(items, pos.Usage())
	}

	return strings.Join(items, " ")
}

func (p *parser[T]) Flags() []params.Flag {
	return p.flags
}

func (p *parser[T]) Args() []params.Arg {
	return p.args
}

func (p *parser[T]) Parse(args []string) (*T, map[string][]string, []string, error) {
	argv := []string{}
	// parse flags.
ARGS:
	for i := 0; i < len(args); i += 1 {
		token := args[i]
		if token == "--" {
			if i+1 < len(args) {
				argv = append(argv, args[i+1:]...)
			}
			break
		}

		token, val, eqok := strings.Cut(token, "=")
		flagName, ok := seemsFlag(token)
		if !ok {
			argv = append(argv, args[i])
			continue
		}

		for _, f := range p.flags {
			lookAhead := 0
			if !f.Match(flagName) {
				continue
			}

			if !eqok {
				lookAhead += 1
				if len(args) <= i+lookAhead {
					if err := f.Found(); err != nil {
						return nil, nil, nil, fmt.Errorf("%w: %s", err, token)
					}
					break ARGS
				}
				val = args[i+lookAhead]
			}

			if err := f.Set(val); err == nil {
				i += lookAhead
			} else {
				if eqok || !errors.Is(err, params.ErrPushBack) {
					return nil, nil, nil, fmt.Errorf("%w: %s", err, token)
				} else if _err := f.Found(); _err != nil {
					return nil, nil, nil, fmt.Errorf("%w: %s", err, token)
				}
			}

			continue ARGS
		}

		argv = append(argv, args[i])
	}

	if len(p.args) == 0 {
		return p.dest, map[string][]string{}, argv, nil
	}

	// assign posargs
	requiredPosArgs := 0
	foundPosArgs := map[string][]string{}
	for _, pos := range p.args {
		foundPosArgs[pos.Name()] = []string{}
		if pos.Required() {
			requiredPosArgs += 1
		}
	}
	set := false

	posargs := p.args[:]

	for 0 < len(posargs) {
		p := posargs[0]
		restv := len(argv)
		if restv <= 0 {
			break
		}
		if restv <= requiredPosArgs {
			if set {
				posargs = posargs[1:]
				p = posargs[0]
				set = false
			}

			if p.Required() && 0 < restv {
				foundPosArgs[p.Name()] = append(foundPosArgs[p.Name()], argv[0])
				argv = argv[1:]
				requiredPosArgs -= 1
			}

			if 0 < len(posargs) {
				posargs = posargs[1:]
			}
			continue
		}

		foundPosArgs[p.Name()] = append(foundPosArgs[p.Name()], argv[0])
		argv = argv[1:]
		if !set && p.Required() {
			requiredPosArgs -= 1
		}
		set = true
		if !p.Repeatable() {
			set = false
			posargs = posargs[1:]
		}
	}

	if 0 < requiredPosArgs {
		return nil, nil, nil, ErrNotEnoughArgs
	}

	return p.dest, foundPosArgs, argv, nil
}

func seemsFlag(arg string) (name string, ok bool) {
	if arg == "--" || arg == "" || arg[0] != '-' {
		return "", false
	}

	l := len(arg)

	if arg[:2] == "--" && arg[:3] != "---" && 4 <= l {
		return arg[2:], true
	}
	if arg[0] == '-' && l == 2 {
		return arg[1:], true
	}

	return "", false
}

var ErrNotEnoughArgs = fmt.Errorf("%w: not enough args", flarcerror.ErrUsage)
