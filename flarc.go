package flarc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/youta-t/flarc/flarcerror"
	"github.com/youta-t/flarc/help"
	"github.com/youta-t/flarc/params"
	"github.com/youta-t/flarc/parser"
)

var ErrUsage = flarcerror.ErrUsage

type helper struct {
	Help bool `alias:"h" help:"show help message"`
}

func helpParser() (parser.Parser[helper], error) {
	return parser.New(&helper{}, []params.ArgDef{})
}

type Command interface {
	ShortDescription() string

	prepare(
		invokedAs string,
		stdin io.Reader,
		stdout io.Writer,
		stderr io.Writer,
		args []string,
		params ...any,
	) runner
}

type runOption struct {
	name    string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	useHelp bool
	argv    []string
	params  []any
}

type RunOption func(*runOption) *runOption

// overwrite args.
//
// If do not pass this, args is defaulted as os.Args[1:]
func WithArgs(argv []string) RunOption {
	return func(ro *runOption) *runOption {
		ro.argv = argv
		return ro
	}
}

// WithHelp set global flag --help, -h
func WithHelp(need bool) RunOption {
	return func(ro *runOption) *runOption {
		ro.useHelp = need
		return ro
	}
}

// Replace Stdin
func WithInput(in io.Reader) RunOption {
	return func(ro *runOption) *runOption {
		ro.stdin = in
		return ro
	}
}

func WithName(name string) RunOption {
	return func(ro *runOption) *runOption {
		ro.name = name
		return ro
	}
}

// Replace Stdout and Stderr
//
// If passing nil as out or errout, it is handled as io.Discard.
func WithOutput(out io.Writer, errout io.Writer) RunOption {
	if out == nil {
		out = io.Discard
	}
	if errout == nil {
		errout = io.Discard
	}

	return func(ro *runOption) *runOption {
		ro.stdout = out
		ro.stderr = errout
		return ro
	}
}

// Pass extra parameters
//
// If pass this multiple times, parameters are appended with previous ones.
func WithParams(param []any) RunOption {
	return func(ro *runOption) *runOption {
		ro.params = append(ro.params, param...)
		return ro
	}
}

// Run command.
//
// # Args
//
// - ctx context.Context
//
// - cmd: command to be executed
//
// - options: options.
func Run(ctx context.Context, cmd Command, options ...RunOption) int {
	runOpt := &runOption{
		name:    filepath.Base(os.Args[0]),
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		argv:    os.Args[1:],
		useHelp: true,
	}
	for _, o := range options {
		runOpt = o(runOpt)
	}

	argv := runOpt.argv
	showHelp := false
	var helpPsr parser.Parser[helper]
	if runOpt.useHelp {
		hp, err := helpParser()
		if err != nil {
			fmt.Fprintln(runOpt.stderr, err)
			return 1
		}
		helpPsr = hp

		hf, _, argv_, err := helpPsr.Parse(argv)
		if err != nil {
			fmt.Fprintln(runOpt.stderr, err)
			return 1
		}
		showHelp = hf.Help
		argv = argv_
	}

	r := cmd.prepare(
		runOpt.name,
		runOpt.stdin, runOpt.stdout, runOpt.stderr,
		argv, runOpt.params...,
	)

	if showHelp {
		hlp := r.Help()
		if helpPsr != nil {
			hlp.AppendFlags(helpPsr.Flags()...)
		}
		hlp.Write(runOpt.stderr)
		return 0
	}

	if err := r.Run(ctx); err == nil {
		return 0
	} else if errors.Is(err, ErrUsage) {
		fmt.Fprintln(runOpt.stderr, err)
		fmt.Fprintln(runOpt.stderr)

		hlp := r.Help()
		if helpPsr != nil {
			hlp.AppendFlags(helpPsr.Flags()...)
		}
		hlp.Write(runOpt.stderr)
		return 2
	} else {
		fmt.Fprintln(runOpt.stderr, err)
		return 1
	}
}

type runner struct {
	Run  func(context.Context) error
	Help func() help.Help
}

// FindParam finds T-typed value from params.
func FindParam[T any](params []any) (T, bool) {
	for _, p := range params {
		t, ok := p.(T)
		if ok {
			return t, true
		}
	}

	return *new(T), false
}
