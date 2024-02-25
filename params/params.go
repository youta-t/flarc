package params

import (
	"reflect"

	"github.com/youta-t/flarc/params/internal/args"
	"github.com/youta-t/flarc/params/internal/flags"
)

type Flag flags.Flag
type Arg args.Arg
type ArgDef args.ArgDef

func (a ArgDef) Freeze() Arg {
	return (args.ArgDef)(a).Freeze()
}

var ErrParse = flags.ErrParse
var ErrPushBack = flags.ErrPushBack
var ErrValueRequired = flags.ErrValueRequired

func NewFlag(tfld reflect.StructField, dest reflect.Value) (Flag, error) {
	return flags.New(tfld, dest)
}
