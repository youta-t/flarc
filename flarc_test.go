package flarc_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/youta-t/flarc"
	"github.com/youta-t/flarc/internal/gen_mock"
	"github.com/youta-t/its"
	"github.com/youta-t/its/itskit"
	"github.com/youta-t/its/mocker/scenario"
)

type CommandlineSpec[T any] struct {
	Fullname its.Matcher[string]
	Flags    its.Matcher[T]
	Args     its.Matcher[map[string][]string]
}

func nilsafe[T any](m its.Matcher[T]) its.Matcher[T] {
	if m != nil {
		return m
	}
	return its.Never[T]()
}

func ItsCommandline[T any](spec CommandlineSpec[T]) its.Matcher[flarc.Commandline[T]] {
	cancel := itskit.SkipStack()
	defer cancel()

	return its.All[flarc.Commandline[T]](
		itskit.Property(
			".Fullname()", flarc.Commandline[T].Fullname,
			nilsafe(spec.Fullname),
		),
		itskit.Property(
			".Flags()", flarc.Commandline[T].Flags,
			nilsafe(spec.Flags),
		),
		itskit.Property(
			".Args()", flarc.Commandline[T].Args,
			nilsafe(spec.Args),
		),
	)
}

func TestCommand(t *testing.T) {
	type Flag struct {
		F bool `alias:"flag" help:"help message" metavar:"FLAG"`
	}

	type When struct {
		commandName      string
		shortDescription string

		args   []string
		params []any
		stdin  string
	}

	type Then struct {
		status its.Matcher[int]
		stdout its.Matcher[string]
		stderr its.Matcher[string]
	}

	theory := func(
		when When,
		task *gen_mock.TaskBehaviour[Flag],
		then Then,
	) func(*testing.T) {
		return func(t *testing.T) {
			sc := scenario.Begin(t)
			defer sc.End()

			var taskfn flarc.Task[Flag]
			if task != nil {
				_, taskfn = scenario.Next(sc, task.Mock(t))
			}

			cmd, err := flarc.NewCommand(
				when.shortDescription,
				Flag{},
				flarc.Args{
					{Name: "source", Repeatable: true, Help: "something input"},
					{Name: "dest", Required: true, Help: "where output is written"},
				},
				taskfn,
				flarc.WithDescription("description..."),
			)
			if err != nil {
				t.Fatal(err)
			}

			its.EqEq(when.shortDescription).Match(cmd.ShortDescription()).OrError(t)

			stdin := new(bytes.Buffer)
			stdin.WriteString(when.stdin)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)
			ctx := context.Background()
			status := flarc.Run(
				ctx, cmd,
				flarc.WithName(when.commandName),
				flarc.WithArgs(when.args),
				flarc.WithInput(stdin),
				flarc.WithOutput(stdout, stderr),
				flarc.WithParams(when.params),
			)
			then.status.Match(status).OrError(t)

			then.stdout.Match(stdout.String()).OrError(t)
			then.stderr.Match(stderr.String()).OrError(t)
		}
	}

	t.Run("command return nil", theory(
		When{
			commandName:      "test",
			shortDescription: "short description...",
			stdin:            "stdin!",
			args:             []string{"-f", "source1", "source2", "dest"},
		},
		gen_mock.NewTaskCall[Flag](
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[Flag]{
				Fullname: its.EqEq("test"),
				Flags:    its.EqEq(Flag{F: true}),
				Args: its.Map(its.MapSpec[string, []string]{
					"source": its.Slice(its.EqEq("source1"), its.EqEq("source2")),
					"dest":   its.Slice(its.EqEq("dest")),
				}),
			}),
			its.Slice[any](),
		).
			ThenEffect(func(_ context.Context, f flarc.Commandline[Flag], params []any) error {
				read := make([]byte, 6)
				io.ReadAtLeast(f.Stdin(), read, len(read))

				its.EqEq("stdin!").Match(string(read)).OrError(t)

				fmt.Fprint(f.Stdout(), "stdout!")
				fmt.Fprint(f.Stderr(), "stderr!")
				return nil
			}),
		Then{
			status: its.EqEq(0),
			stdout: its.EqEq("stdout!"),
			stderr: its.EqEq("stderr!"),
		},
	))

	t.Run("command return nil, with extra params", theory(
		When{
			commandName:      "test",
			shortDescription: "short description...",
			stdin:            "stdin!",
			args:             []string{"-f", "source1", "source2", "dest"},
			params:           []any{42, "foo"},
		},
		gen_mock.NewTaskCall[Flag](
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[Flag]{
				Fullname: its.EqEq("test"),
				Flags:    its.EqEq(Flag{F: true}),
				Args: its.Map(its.MapSpec[string, []string]{
					"source": its.Slice(its.EqEq("source1"), its.EqEq("source2")),
					"dest":   its.Slice(its.EqEq("dest")),
				}),
			}),
			its.Slice[any](
				its.EqEq[any](42),
				its.EqEq[any]("foo"),
			),
		).
			ThenReturn(nil),
		Then{
			status: its.EqEq(0),
			stdout: its.EqEq(""),
			stderr: its.EqEq(""),
		},
	))

	t.Run("return ErrUsage", theory(
		When{
			commandName: "test", shortDescription: "testing flarg",
			stdin: "stdin!",
			args:  []string{"-f", "source1", "source2", "dest"},
		},
		gen_mock.NewTaskCall[Flag](
			its.Always[context.Context](),
			its.Always[flarc.Commandline[Flag]](),
			its.Slice[any](),
		).
			ThenReturn(flarc.ErrUsage),
		Then{
			status: its.EqEq(2),
			stdout: its.EqEq(""),
			stderr: its.Text(`usage error

test -- testing flarg

Usage:

    test -f=FLAG --help=false [source[, ...]] dest

Description:

    description...

Flags:

    -f, --flag  help message
    --help, -h  show help message

Args:

    source      something input
    dest        where output is written
`),
		},
	))

	t.Run("on parse error", theory(
		When{
			commandName: "test", shortDescription: "testing flarg",
			stdin: "stdin!",
			args:  []string{"-f=not-bool", "source1", "source2", "dest"},
		},
		nil,
		Then{
			status: its.EqEq(2),
			stdout: its.EqEq(""),
			stderr: its.Text(`usage error: parse error: not-bool is not bool: -f

test -- testing flarg

Usage:

    test -f=FLAG --help=false [source[, ...]] dest

Description:

    description...

Flags:

    -f, --flag  help message
    --help, -h  show help message

Args:

    source      something input
    dest        where output is written
`),
		},
	))

	t.Run("return other error", theory(
		When{
			commandName: "test",
			stdin:       "stdin!",
			args:        []string{"-f", "source1", "source2", "dest"},
		},
		gen_mock.NewTaskCall[Flag](
			its.Always[context.Context](),
			its.Always[flarc.Commandline[Flag]](),
			its.Slice[any](),
		).
			ThenReturn(errors.New("fake error")),
		Then{
			status: its.EqEq(1),
			stdout: its.EqEq(""),
			stderr: its.EqEq("fake error\n"),
		},
	))
}

func TestSubcommand(t *testing.T) {

	type FlagSuper struct {
		I int `alias:"int" help:"help for command group flag" metavar:"N"`
	}

	type FlagSub struct {
		F bool `alias:"flag" help:"help message" metavar:"FLAG"`
	}

	type When struct {
		args   []string
		params []any
		stdin  string
	}

	type Then struct {
		status its.Matcher[int]
		stdout its.Matcher[string]
		stderr its.Matcher[string]
	}

	theory := func(when When, taskMock *gen_mock.TaskBehaviour[FlagSub], then Then) func(*testing.T) {
		return func(t *testing.T) {
			sc := scenario.Begin(t)
			defer sc.End()

			var task flarc.Task[FlagSub]
			if taskMock != nil {
				_t := taskMock.Mock(t)
				_, task = scenario.Next(sc, _t)
			}

			subcom, err := flarc.NewCommand(
				"this is subcommand, called {{ .Command }}",
				FlagSub{}, flarc.Args{
					{Name: "arg1", Required: true, Help: "subcommand's arg"},
				},
				task,
				flarc.WithDescription("this is subcommand description"),
			)
			if err != nil {
				t.Fatal(err)
			}

			cg, err := flarc.NewCommandGroup(
				"this is command group, called {{ .Command }}",
				FlagSuper{},
				flarc.WithGroupDescription("this is group description."),
				flarc.WithSubcommand("sub", subcom),
			)
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()
			stdin := new(bytes.Buffer)
			stdin.WriteString(when.stdin)
			stdout := new(strings.Builder)
			stderr := new(strings.Builder)
			status := flarc.Run(
				ctx, cg,
				flarc.WithName("test"),
				flarc.WithInput(stdin),
				flarc.WithOutput(stdout, stderr),
				flarc.WithArgs(when.args),
				flarc.WithParams(when.params),
			)

			then.status.Match(status).OrError(t)
			then.stdout.Match(stdout.String()).OrError(t)
			then.stderr.Match(stderr.String()).OrError(t)
		}
	}

	t.Run("`-h` shows commandgroup help", theory(
		When{
			args: []string{"-h"},
		},
		nil,
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(`test -- this is command group

Usage:

    test -i=N --help=false

Description:

    this is group description.

Flags:

    -i, --int   help for command group flag
    --help, -h  show help message

Subcommands:

    sub         this is subcommand

`),
		},
	))

	t.Run("`-i 99` shows commandgroup help with `no subcommand` error", theory(
		When{
			args: []string{"-i", "99"},
		},
		nil,
		Then{
			status: its.EqEq(2),
			stdout: its.Text(""),
			stderr: its.Text(`usage error: no subcommands

test -- this is command group

Usage:

    test -i=N --help=false

Description:

    this is group description.

Flags:

    -i, --int   help for command group flag
    --help, -h  show help message

Subcommands:

    sub         this is subcommand

`),
		},
	))

	t.Run("`-i 99 nosub` shows commandgroup help with `unknwon subcommand` error", theory(
		When{
			args: []string{"-i", "99", "nosub"},
		},
		nil,
		Then{
			status: its.EqEq(2),
			stdout: its.Text(""),
			stderr: its.Text(`usage error: unknown subcommand: nosub

test -- this is command group

Usage:

    test -i=N --help=false

Description:

    this is group description.

Flags:

    -i, --int   help for command group flag
    --help, -h  show help message

Subcommands:

    sub         this is subcommand

`),
		},
	))

	t.Run("`-h sub aaa` shows subcommand help", theory(
		When{
			args: []string{"-h", "sub", "aaa"},
		},
		nil,
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(`test sub -- this is subcommand

Usage:

    test sub -f=FLAG -i=N --help=false arg1

Description:

    this is subcommand description

Flags:

    -f, --flag  help message
    -i, --int   help for command group flag
    --help, -h  show help message

Args:

    arg1        subcommand's arg
`),
		},
	))

	t.Run("`sub -h aaa` shows subcommand help", theory(
		When{
			args: []string{"sub", "-h", "aaa"},
		},
		nil,
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(`test sub -- this is subcommand

Usage:

    test sub -f=FLAG -i=N --help=false arg1

Description:

    this is subcommand description

Flags:

    -f, --flag  help message
    -i, --int   help for command group flag
    --help, -h  show help message

Args:

    arg1        subcommand's arg
`),
		},
	))

	t.Run("`sub aaa` invokes subcommand", theory(
		When{
			args:  []string{"sub", "aaa"},
			stdin: "stdin!!",
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: false}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.Slice(
				its.EqEq[any](FlagSuper{I: 0}),
			),
		).ThenEffect(func(arg0 context.Context, arg1 flarc.Commandline[FlagSub], arg2 []any) error {
			buf := make([]byte, 7)
			io.ReadAtLeast(arg1.Stdin(), buf, len(buf))
			its.EqEq("stdin!!").Match(string(buf)).OrError(t)

			fmt.Fprintf(arg1.Stdout(), "stdout!!!")
			fmt.Fprintf(arg1.Stderr(), "stderr!!!")
			return nil
		}),
		Then{
			status: its.EqEq(0),
			stdout: its.Text("stdout!!!"),
			stderr: its.Text("stderr!!!"),
		},
	))

	t.Run("`sub -i 99 -f true aaa` invokes subcommand with flags", theory(
		When{
			args:  []string{"sub", "-i", "99", "-f", "true", "aaa"},
			stdin: "stdin!!",
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: true}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.Slice(
				its.EqEq[any](FlagSuper{I: 99}),
			),
		).ThenReturn(nil),
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(""),
		},
	))

	t.Run("`sub -i 99 --flag true aaa` invokes subcommand with flags", theory(
		When{
			args:  []string{"sub", "-i", "99", "--flag", "true", "aaa"},
			stdin: "stdin!!",
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: true}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.Slice(
				its.EqEq[any](FlagSuper{I: 99}),
			),
		).ThenReturn(nil),
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(""),
		},
	))

	t.Run("`sub -i 99 --flag true aaa` with params invokes subcommand with flags", theory(
		When{
			args:   []string{"sub", "-i", "99", "--flag", "true", "aaa"},
			params: []any{42, "param"},
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: true}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.SliceUnordered(
				its.EqEq[any](42),
				its.EqEq[any]("param"),
				its.EqEq[any](FlagSuper{I: 99}),
			),
		).ThenReturn(nil),
		Then{
			status: its.EqEq(0),
			stdout: its.Text(""),
			stderr: its.Text(""),
		},
	))

	t.Run("when subcommand returns ErrUsage, it prints help", theory(
		When{
			args: []string{"sub", "aaa"},
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: false}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.Slice(
				its.EqEq[any](FlagSuper{I: 0}),
			),
		).ThenReturn(fmt.Errorf("%w: fake error", flarc.ErrUsage)),
		Then{
			status: its.EqEq(2),
			stdout: its.Text(""),
			stderr: its.Text(`usage error: fake error

test sub -- this is subcommand

Usage:

    test sub -f=FLAG -i=N --help=false arg1

Description:

    this is subcommand description

Flags:

    -f, --flag  help message
    -i, --int   help for command group flag
    --help, -h  show help message

Args:

    arg1        subcommand's arg
`),
		},
	))

	t.Run("when subcommand returns error, it exit with 1", theory(
		When{
			args: []string{"sub", "aaa"},
		},
		gen_mock.NewTaskCall(
			its.Always[context.Context](),
			ItsCommandline(CommandlineSpec[FlagSub]{
				Fullname: its.EqEq("test sub"),
				Flags:    its.EqEq(FlagSub{F: false}),
				Args: its.Map(its.MapSpec[string, []string]{
					"arg1": its.Slice(its.EqEq("aaa")),
				}),
			}),
			its.Slice(
				its.EqEq[any](FlagSuper{I: 0}),
			),
		).ThenReturn(errors.New("fake error")),
		Then{
			status: its.EqEq(1),
			stdout: its.Text(""),
			stderr: its.Text(`fake error
`),
		},
	))
}
