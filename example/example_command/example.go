package main

// ...
import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/youta-t/flarc"
)

// ...

// declare struct for flags.
// Flag options can be set with tag.
type Flag struct {
	Foo  string        `alias:"f" metavar:"FOO" help:"flag foo"`
	Bar  int           // when no tag, assumed `flag:"${field-name-in-kebab-case}"` is given.
	Fizz bool          `alias:"F"` // flag name is case sensitive
	Bazz time.Duration `metavar:"DURATION"`
}

func main() {

	cmd, err := flarc.NewCommand(
		"short description...",

		// declare default values for flags
		Flag{
			Foo:  "default foo",
			Bar:  42,
			Fizz: false,
			Bazz: 3 * time.Second,
		},

		// declare positional args
		flarc.Args{
			{
				Name:       "SOURCE", // positional arg name
				Repeatable: true,     // set true if this arg takes many items
				Required:   true,     // require at least 1 item
				Help:       "help message of SOURCE",
			},
			{
				Name: "DEST", Required: true,
				Help: "help message of DEST",
			},
		},

		func(
			ctx context.Context,
			commandline flarc.Commandline[Flag],
			param []any,
		) error {
			// parsed flag
			var flag Flag = commandline.Flags()

			// key is arg's name. values are assigned commandline
			var args map[string][]string = commandline.Args()

			buf := make([]byte, 1024)
			io.ReadAtLeast(commandline.Stdin(), buf, 0)

			fmt.Fprintf(commandline.Stdout(), "passed flag: %+v\n\n", flag)
			fmt.Fprintf(commandline.Stdout(), "passed arg: %+v\n\n", args)
			fmt.Fprintf(commandline.Stdout(), "passed params: %+v\n\n", param)

			return nil
		},

		flarc.WithDescription(`this is example command.

To show how to declare flags and args, also genereatad help message.

This command is called "{{ .Command }}"
`),
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "unexpected error:", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt,
	)
	defer cancel()

	type param struct {
		ParamValue string
	}

	os.Exit(flarc.Run(
		ctx, cmd,

		// optional extra params.
		//
		// For example, loggers can be injected with this.
		flarc.WithParams([]any{param{ParamValue: "this is param value"}}),
	))
}
