package args

// ArgDef holds configuration of positional arguments
type ArgDef struct {
	// Name of this arg
	Name string

	// true if this arg is required
	Required bool

	// true if this args is repeatable.
	Repeatable bool

	// Help message for this arg.
	Help string
}

func (p ArgDef) Freeze() Arg {
	return &arg{
		name:       p.Name,
		required:   p.Required,
		repeatable: p.Repeatable,
		help:       p.Help,
	}
}

type Arg interface {
	Name() string
	Required() bool
	Repeatable() bool
	Help() string
	Usage() string
}

type arg struct {
	name       string
	required   bool
	repeatable bool
	help       string
}

func (p arg) Name() string {
	return p.name
}

func (p arg) Required() bool {
	return p.required
}

func (p arg) Repeatable() bool {
	return p.repeatable
}

func (p arg) Help() string {
	return p.help
}

func (pos arg) Usage() string {
	s := pos.name
	if pos.repeatable {
		s = s + "[, ...]"
	}
	if !pos.required {
		s = "[" + s + "]"
	}
	return s
}
