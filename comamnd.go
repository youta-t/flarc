//go:generate go run github.com/youta-t/its/mocker -t Task --dest internal/gen_mock
package flarc

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/youta-t/flarc/flarcerror"
	"github.com/youta-t/flarc/help"
	"github.com/youta-t/flarc/params"
	"github.com/youta-t/flarc/parser"
)

type Task[T any] func(context.Context, Commandline[T], []any) error

type Args []params.ArgDef

// NewCommand creates new command.
//
// # Args
//
// - name: the name of this command
//
// - shortDescription: short (one line or less) description for this command.
//
// - flagdef: a struct defining flags for this command
//
// - args: declaration of positional arguments
//
// - task: task of this command.
// Commandline.Flag() reflected flagdef. Commandline.Args() is reflected pos (keys are names of pos).
//
// - option: options
func NewCommand[T any](
	shortDescription string,
	flagdef T, args Args, task Task[T],
	option ...CommandOption,
) (Command, error) {
	opt := &commandOption{}
	for _, f := range option {
		var err error
		opt, err = f(opt)
		if err != nil {
			return nil, err
		}
	}

	aa := make([]params.ArgDef, len(args))
	copy(aa, args)

	parser, err := parser.New(&flagdef, args)
	if err != nil {
		return nil, err
	}

	return command[T]{
		shortDescription: shortDescription,
		task:             task,
		parser:           parser,
		description:      opt.description,
	}, nil
}

type CommandOption func(*commandOption) (*commandOption, error)

type commandOption struct {
	description *template.Template
}

func WithDescription(d string) CommandOption {
	return func(p *commandOption) (*commandOption, error) {
		var err error
		p.description, err = template.New("").Parse(d)
		return p, err
	}
}

type command[T any] struct {
	shortDescription string
	description      *template.Template

	parser parser.Parser[T]

	task Task[T]
}

func (cmd command[T]) ShortDescription() string {
	return cmd.shortDescription
}

func (cmd command[T]) newHelp(fullname string) help.Help {
	h := help.New(
		fullname, cmd.ShortDescription(),
		help.WithDescription(cmd.description),
		help.WitArgs(cmd.parser.Args()),
	)
	h.AppendFlags(cmd.parser.Flags()...)
	return h
}

func (cmd command[T]) prepare(
	fullname string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args []string,
	params ...any,
) runner {
	flags, argv, rem, err := cmd.parser.Parse(args)
	if err != nil {
		return runner{
			Run:  func(context.Context) error { return err },
			Help: func() help.Help { return cmd.newHelp(fullname) },
		}
	}
	if 0 < len(rem) {
		return runner{
			Run: func(context.Context) error {
				return fmt.Errorf("%w: too much args", flarcerror.ErrUsage)
			},
			Help: func() help.Help { return cmd.newHelp(fullname) },
		}
	}

	cl := commandline[T]{
		fullname: fullname,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
		flags:    *flags,
		args:     argv,
		params:   params,
	}

	return runner{
		Run: func(ctx context.Context) error {
			return cmd.task(ctx, cl, params)
		},
		Help: func() help.Help { return cmd.newHelp(fullname) },
	}
}

type commandline[T any] struct {
	fullname string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	flags  T
	args   map[string][]string
	params []any
}

// Commandline represents commandline interface.
type Commandline[T any] interface {
	// Fullname of this command.
	Fullname() string

	// Stdin returns standard input as io.Reader.
	Stdin() io.Reader

	// Stdout returns standard output as io.Writer.
	Stdout() io.Writer

	// Stderr returns standard error as io.Writer.
	Stderr() io.Writer

	// Flags returns flags passed on commandline.
	Flags() T

	// Args returns positional argument values for each positinal arguments' name.
	Args() map[string][]string
}

func (t commandline[T]) Fullname() string {
	return t.fullname
}

func (t commandline[T]) Stdin() io.Reader {
	return t.stdin
}

func (t commandline[T]) Stdout() io.Writer {
	return t.stdout
}

func (t commandline[T]) Stderr() io.Writer {
	return t.stderr
}

func (t commandline[T]) Flags() T {
	return t.flags
}

func (t commandline[T]) Args() map[string][]string {
	return t.args
}

func (t commandline[T]) Params() []any {
	return t.params
}
