flarc -- commandline parser, supports flag, arg, subcommands
============================================================

flarc is a commandline parser. With flarc, you can

- declare *fl*ags
- declare positional *ar*guments, and
- declare sub*c*ommands
- get generated help text help texts

Install
--------

```
go get githuc.com/youta-t/flarc
```

flarc needs go1.21+.

Usage
------

### Define command

```go
// ...
import "github.com/youta-t/flarc"
// ...

// declare struct for flags.
// Flag options can be set with tag.
type Flag struct {
    Foo string  `alias:"f" help:"flag foo" metavar:"FOO"`
    Bar int     // when no tag, assumed `flag:"${field-name-in-kebab-case}"` is given.
    Fizz bool   `flag:"F"`  // flag name is case sensitive
    Bazz time.Duration  `metavar:"DURATION"`
}

func main() {

// ...

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


	// ...
}
```

In description, you can use placeholder `{{ .Command }}`.
This will be filled with the command name at invoke time.

### Run command

```go
func main() {
	// ...

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
```

Return status is selected by flarc.

If the task of command returns,

- `nil`: success! exits with `0`.
- `flarc.ErrUsage`: prints help message and exits with `2`.
- other error: exits with `1`.


```
$ go run ./example/example_command -f foo -F false source1 source2 dest1
passed flag: {Foo:foo Bar:42 Fizz:false Bazz:3s}

passed arg: map[DEST:[dest1] SOURCE:[source1 source2]]

passed params: [{ParamValue:this is param value}]

```

By default, flarc provides `--help, -h` flag to show help message.

```
$ go run ./example/example_command --help
example_command -- short description...

Usage:

    example_command --foo=FOO --bar=0 --fizz=false --bazz=DURATION --help=false SOURCE[, ...] DEST

Description:

    this is example command.
    
    To show how to declare flags and args, also genereatad help message.
    
    This command is called "example_command"
    

Flags:

    --foo, -f   flag foo
    --bar
    --fizz, -F
    --bazz
    --help, -h  show help message

Args:

    SOURCE      help message of SOURCE
    DEST        help message of DEST
```

### Define Command Group and Subcommand

```go
// ...

type GroupFlag struct {
	Qux  string `alias:"q" help:"help for command group flag"`
	Quux string `alias:"Q"`
}

//...

func main() {

	// ...

	grp, err := flarc.NewCommandGroup(
		"description of command group",
		GroupFlag{
			Qux:  "qux",
			Quux: "QUUX",
		},
		flarc.WithGroupDescription(`command group description.

This command is called as "{{ .Command }}".
`),
		flarc.WithSubcommand("sub", cmd),
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "unexpected error:", err)
		os.Exit(1)
	}

	// ...
}
```

#### run it

Help for command group:

```
$ go run ./example/example_subcommand --help
example_subcommand -- description of command group

Usage:

    example_subcommand --qux --quux --help=false

Description:

    command group description.
    
    This command is called as "example_subcommand".
    

Flags:

    --qux, -q   help for command group flag
    --quux, -Q
    --help, -h  show help message

Subcommands:

    sub         short description...
```

Help for subcommand:

```
$ go run ./example/example_subcommand sub --help
example_subcommand sub -- short description...

Usage:

    example_subcommand sub --foo=FOO --bar=0 --fizz=false --bazz=DURATION --qux --quux --help=false SOURCE[, ...] DEST

Description:

    this is example command.
    
    To show how to declare flags and args, also genereatad help message.
    

Flags:

    --foo, -f   flag foo
    --bar
    --fizz, -F
    --bazz
    --qux, -q   help for command group flag
    --quux, -Q
    --help, -h  show help message

Args:

    SOURCE      help message of SOURCE
    DEST        help message of DEST
```

```
$ go run ./example/example_subcommand sub -f foo -F false -q "queue" -Q "Queue" source1 source2 dest1
passed flag: {Foo:foo Bar:42 Fizz:false Bazz:3s}

passed arg: map[DEST:[dest1] SOURCE:[source1 source2]]

passed params: [{ParamValue:this is param value} {Qux:queue Quux:Queue}]
```
