package flags

import (
	"errors"
	goflag "flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/youta-t/flarc/flarcerror"
	"github.com/youta-t/flarc/utils"
)

// Flag represents commandline flag.
type Flag interface {

	// Name of this flag.
	Name() string

	// Aliases of this flag.
	Alias() []string

	// Match returns true when s in Name() or Alias()
	Match(s string) bool

	// parse and set value
	Set(string) error

	// Found this flag but no value is given
	Found() error

	// Help message
	Help() string

	// usage of this flag
	Usage() string
}

func setfn[T any](dest reflect.Value) func(T) {

	wrap := func(v reflect.Value) reflect.Value {
		return v
	}

	next := dest.Type()
	for {
		switch next.Kind() {
		case reflect.Pointer:
			_w := wrap
			wrap = func(v reflect.Value) reflect.Value {
				p := reflect.New(v.Type())
				p.Elem().Set(v)
				return _w(p)
			}
			next = next.Elem()
			continue
		case reflect.Slice:
			_w := wrap
			s := reflect.MakeSlice(next, 0, 0)
			wrap = func(v reflect.Value) reflect.Value {
				s = reflect.Append(s, v)
				return _w(s)
			}
			next = next.Elem()
			continue
		default:
		}
		break
	}

	return func(t T) {
		dest.Set(wrap(reflect.ValueOf(t)))
	}
}

type flag[T any] struct {
	name       string
	alias      []string
	set        func(T)
	translator func(string) (T, error)
	action     func() (T, error)

	help      string
	metaValue string
}

func (f flag[T]) Usage() string {
	flgname := f.name
	if len(flgname) == 1 {
		flgname = "-" + flgname
	} else {
		flgname = "--" + flgname
	}

	if f.metaValue == "" {
		return flgname
	}

	return fmt.Sprintf("%s=%s", flgname, f.metaValue)
}

func (f flag[T]) hypen(n string) string {
	if len(n) == 1 {
		return "-" + n
	}
	return "--" + n
}

func (f flag[T]) Name() string {
	return f.hypen(f.name)
}

func (f flag[T]) Alias() []string {
	a := make([]string, len(f.alias))
	for i := range f.alias {
		a[i] = f.hypen(f.alias[i])
	}
	return a
}

func (f flag[T]) Match(given string) bool {
	if given == f.name {
		return true
	}
	for _, pat := range f.alias {
		if pat == given {
			return true
		}
	}
	return false
}

func (f flag[T]) Set(s string) error {
	val, err := f.translator(s)
	if err != nil {
		return err
	}
	f.set(val)
	return nil
}

func (f flag[T]) Found() error {
	if f.action == nil {
		return fmt.Errorf("%w: %s", ErrValueRequired, f.name)
	}
	val, err := f.action()
	if err != nil {
		return err
	}
	f.set(val)
	return nil
}

func (f flag[T]) Help() string {
	return f.help
}

func elem(t reflect.Type) reflect.Type {
	next := t
	for {
		switch next.Kind() {
		case reflect.Slice, reflect.Pointer:
			next = next.Elem()
		default:
			return next
		}
	}
}

func New(tfld reflect.StructField, dest reflect.Value) (Flag, error) {

	name := tfld.Tag.Get("flag")
	if name == "" {
		name = utils.ToKebab(tfld.Name)
	}
	alias := []string{}
	if s, ok := tfld.Tag.Lookup("alias"); ok {
		alias = append(alias, strings.Split(s, ",")...)
	}
	help := tfld.Tag.Get("help")

	metavar := ""
	if mv, ok := tfld.Tag.Lookup("metavar"); ok {
		metavar = mv
	}

	switch d := dest.Interface().(type) {
	case func(string) error:
		return flag[string]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set: func(s string) {
				// do nothing.
			},
			translator: func(s string) (string, error) {
				rets := dest.Call([]reflect.Value{reflect.ValueOf(s)})
				e := rets[0].Interface().(error)
				if e != nil {
					return "", e
				}
				return "", nil
			},
		}, nil
	case goflag.Value:
		if metavar == "" {
			metavar = d.String()
		}
		return flag[goflag.Value]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set: func(v goflag.Value) {},
			translator: func(s string) (goflag.Value, error) {
				err := d.Set(s)
				return d, err
			},
		}, nil
	}

	switch def := reflect.New(elem(tfld.Type)).Elem().Interface().(type) {
	case string:
		if metavar == "" {
			metavar = def
		}
		return flag[string]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[string](dest),
			action:     func() (string, error) { return "", ErrValueRequired },
			translator: readString,
		}, nil
	case bool:
		if metavar == "" {
			metavar = fmt.Sprintf("%#v", def)
		}
		return flag[bool]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[bool](dest),
			action:     func() (bool, error) { return true, nil },
			translator: readBool,
		}, nil
	case int:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[int]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[int](dest),
			action:     func() (int, error) { return 0, ErrValueRequired },
			translator: readInt[int],
		}, nil
	case int8:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[int8]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[int8](dest),
			action:     func() (int8, error) { return 0, ErrValueRequired },
			translator: readInt[int8],
		}, nil
	case int16:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[int16]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[int16](dest),
			action:     func() (int16, error) { return 0, ErrValueRequired },
			translator: readInt[int16],
		}, nil
	case int32:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[int32]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[int32](dest),
			action:     func() (int32, error) { return 0, ErrValueRequired },
			translator: readInt[int32],
		}, nil
	case int64:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[int64]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[int64](dest),
			action:     func() (int64, error) { return 0, ErrValueRequired },
			translator: readInt[int64],
		}, nil
	case uint:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[uint]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[uint](dest),
			action:     func() (uint, error) { return 0, ErrValueRequired },
			translator: readUint[uint],
		}, nil
	case uint8:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[uint8]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[uint8](dest),
			action:     func() (uint8, error) { return 0, ErrValueRequired },
			translator: readUint[uint8],
		}, nil
	case uint16:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[uint16]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[uint16](dest),
			action:     func() (uint16, error) { return 0, ErrValueRequired },
			translator: readUint[uint16],
		}, nil
	case uint32:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[uint32]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[uint32](dest),
			action:     func() (uint32, error) { return 0, ErrValueRequired },
			translator: readUint[uint32],
		}, nil
	case uint64:
		if metavar == "" {
			metavar = fmt.Sprintf("%d", def)
		}
		return flag[uint64]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[uint64](dest),
			action:     func() (uint64, error) { return 0, ErrValueRequired },
			translator: readUint[uint64],
		}, nil
	case float32:
		if metavar == "" {
			metavar = fmt.Sprintf("%f", def)
		}
		return flag[float32]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[float32](dest),
			action:     func() (float32, error) { return 0, ErrValueRequired },
			translator: readFloat[float32],
		}, nil
	case float64:
		if metavar == "" {
			metavar = fmt.Sprintf("%f", def)
		}
		return flag[float64]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[float64](dest),
			action:     func() (float64, error) { return 0, ErrValueRequired },
			translator: readFloat[float64],
		}, nil
	case time.Duration:
		if metavar == "" {
			metavar = def.String()
		}
		return flag[time.Duration]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[time.Duration](dest),
			action:     func() (time.Duration, error) { return 0, ErrValueRequired },
			translator: readDuration,
		}, nil
	case time.Time:
		if metavar == "" {
			metavar = def.Format(time.RFC3339Nano)
		}
		return flag[time.Time]{
			name: name, alias: alias, help: help, metaValue: metavar,
			set:        setfn[time.Time](dest),
			action:     func() (time.Time, error) { return time.Time{}, ErrValueRequired },
			translator: readTime(time.RFC3339Nano),
		}, nil
	}

	return nil, errors.New("unsupported type")
}

var ErrValueRequired = fmt.Errorf("%w: value required", flarcerror.ErrUsage)
