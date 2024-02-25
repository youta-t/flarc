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

// NewCommandGroup creates command group.
//
// # Args
//
// - name: name of this command group.
//
// - shortDescription: short (one line or less) description for this group.
//
// - flagdef: struct defining flags. This flags affected global flag under this group.
//
// - option: options.
func NewCommandGroup[T any](
	shortDescription string,
	flagdef T,
	option ...CommandGroupOption,
) (Command, error) {
	opt := &commandGroupOption{
		subCommands: map[string]Command{},
	}
	for _, f := range option {
		var err error
		opt, err = f(opt)
		if err != nil {
			return nil, err
		}
	}

	parser, err := parser.New(&flagdef, []params.ArgDef{})
	if err != nil {
		return nil, err
	}

	cg := &commandGroup[T]{
		shortDescription: shortDescription,
		parser:           parser,

		description: opt.description,
		subcommands: opt.subCommands,
	}

	return cg, nil
}

type CommandGroupOption func(*commandGroupOption) (*commandGroupOption, error)

type commandGroupOption struct {
	description *template.Template
	subCommands map[string]Command
}

func WithGroupDescription(d string) CommandGroupOption {
	return func(p *commandGroupOption) (*commandGroupOption, error) {
		var err error
		p.description, err = template.New("").Parse(d)
		return p, err
	}
}

func WithSubcommand(name string, c Command) CommandGroupOption {
	return func(p *commandGroupOption) (*commandGroupOption, error) {
		if _, ok := p.subCommands[name]; ok {
			return nil, fmt.Errorf("subcommand name conflicts: %s", name)
		}
		p.subCommands[name] = c
		return p, nil
	}
}

type commandGroup[T any] struct {
	name             string
	shortDescription string
	description      *template.Template

	parser parser.Parser[T]

	subcommands map[string]Command
}

func (cg *commandGroup[T]) Name() string {
	return cg.name
}

func (cg *commandGroup[T]) ShortDescription() string {
	return cg.shortDescription
}

func (cg *commandGroup[T]) newHelp(fullname string) help.Help {
	cmds := map[string]help.CommandDescriptor{}
	for name := range cg.subcommands {
		cmds[name] = cg.subcommands[name]
	}

	return help.New(
		fullname, cg.ShortDescription(),
		help.WithFlags(cg.parser.Flags()),
		help.WithDescription(cg.description),
		help.WithSubcommands(cmds),
	)
}

func (cg *commandGroup[T]) prepare(
	fullname string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args []string,
	params ...any,
) runner {

	flags, _, rem, err := cg.parser.Parse(args)
	if err != nil {
		return runner{
			Run:  func(context.Context) error { return err },
			Help: func() help.Help { return cg.newHelp(fullname) },
		}
	}

	if len(rem) == 0 {
		return runner{
			Run: func(context.Context) error {
				return fmt.Errorf("%w: no subcommands", flarcerror.ErrUsage)
			},
			Help: func() help.Help { return cg.newHelp(fullname) },
		}
	}

	for name, sub := range cg.subcommands {
		if name == rem[0] {
			p := append([]any{}, params...)
			p = append(p, *flags)

			r := sub.prepare(fullname+" "+name, stdin, stdout, stderr, rem[1:], p...)

			return runner{
				Run: r.Run,
				Help: func() help.Help {
					h := r.Help()
					h.AppendFlags(cg.parser.Flags()...)
					return h
				},
			}
		}
	}

	return runner{
		Run: func(ctx context.Context) error {
			return fmt.Errorf("%w: unknown subcommand: %s", flarcerror.ErrUsage, rem[0])
		},
		Help: func() help.Help { return cg.newHelp(fullname) },
	}
}
