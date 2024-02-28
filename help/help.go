package help

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/youta-t/flarc/params"
)

type Option func(*help) *help

func WithFlags(flg []params.Flag) Option {
	return func(h *help) *help {
		h.flags.Append(flg...)
		return h
	}
}

func WitArgs(arg []params.Arg) Option {
	return func(h *help) *help {
		h.args.Append(arg...)
		return h
	}
}

type CommandDescriptor interface {
	ShortDescription() string
}

func WithSubcommands(cmds map[string]CommandDescriptor) Option {
	return func(h *help) *help {
		h.subcommands.Append(cmds)
		return h
	}
}

func WithDescription(tpl *template.Template) Option {
	return func(h *help) *help {
		h.description = tpl
		return h
	}
}

func New(
	fullname string, shortDescription string,
	options ...Option,
) Help {
	h := &help{
		fullname:         fullname,
		shortDescription: shortDescription,
		flags:            new(paramSection[params.Flag]),
		args:             new(paramSection[params.Arg]),

		subcommands: &subcommandSection{
			cmds: map[string]CommandDescriptor{},
		},
	}

	for _, o := range options {
		h = o(h)
	}

	return h
}

type Help interface {
	Write(w io.Writer) error

	AppendFlags(...params.Flag)
}

type help struct {
	fullname         string
	shortDescription string
	description      *template.Template
	flags            *paramSection[params.Flag]
	args             *paramSection[params.Arg]
	subcommands      *subcommandSection
}

func (h *help) AppendFlags(flgs ...params.Flag) {
	h.flags.Append(flgs...)
}

func (h *help) Write(w io.Writer) error {
	fmt.Fprint(w, h.fullname)
	if h.shortDescription != "" {
		fmt.Fprint(w, " -- ", h.shortDescription)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "    %s", h.fullname)
	h.flags.WriteUsage(w)
	h.args.WriteUsage(w)
	fmt.Fprintln(w)

	if h.description != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Description:")
		fmt.Fprintln(w)

		sb := new(strings.Builder)
		err := h.description.Execute(sb, struct{ Command string }{Command: h.fullname})
		if err != nil {
			return err
		}

		for _, line := range strings.Split(sb.String(), "\n") {
			fmt.Fprint(w, "    ")
			fmt.Fprintln(w, line)
		}
	}

	if 0 < h.flags.Len() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Flags:")
		fmt.Fprintln(w)
		h.flags.WriteHelp(w)
	}

	if 0 < h.args.Len() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Args:")
		fmt.Fprintln(w)
		h.args.WriteHelp(w)
	}

	if 0 < h.subcommands.Len() {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Subcommands:")
		fmt.Fprintln(w)
		h.subcommands.WriteHelp(w)
	}

	return nil
}

type paramSection[T interface {
	Name() string
	Help() string
	Usage() string
}] struct {
	content []T
}

func (s *paramSection[T]) Len() int {
	return len(s.content)
}

func (s *paramSection[T]) Append(item ...T) {
	s.content = append(s.content, item...)
}

func (s *paramSection[T]) Merge(o *paramSection[T]) {
	s.Append(o.content...)
}

func (s *paramSection[T]) WriteUsage(w io.Writer) {
	for _, c := range s.content {
		fmt.Fprint(w, " ", c.Usage())
	}
}

func (s *paramSection[T]) WriteHelp(w io.Writer) {
	for _, c := range s.content {
		var names []string
		if a, ok := any(c).(interface{ Alias() []string }); !ok {
			names = []string{c.Name()}
		} else {
			alias := a.Alias()
			names = make([]string, 1+len(alias))
			names[0] = c.Name()
			copy(names[1:], alias)
		}

		helpText := c.Help()
		n := strings.Join(names, ", ")

		if helpText == "" {
			fmt.Fprintf(w, "    %s", n)
		} else if len(n) <= 12 {
			s := fmt.Sprintf("    %-12s", n)
			fmt.Fprint(w, s)
		} else {
			fmt.Fprintf(w, "    %s\n            ", n)
		}

		help := strings.Split(helpText, "\n")
		if 0 < len(help) {
			fmt.Fprint(w, help[0])
			help = help[1:]
		}
		for _, l := range help {
			fmt.Fprint(w, "           ", l)
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w)
	}
}

type subcommandSection struct {
	cmds map[string]CommandDescriptor
}

func (scs *subcommandSection) Len() int {
	return len(scs.cmds)
}

func (scs *subcommandSection) Append(cmds map[string]CommandDescriptor) {
	for name, cmd := range cmds {
		scs.cmds[name] = cmd
	}
}

func (scs *subcommandSection) WriteHelp(w io.Writer) {

	type subcommand struct {
		Name    string
		Command CommandDescriptor
	}

	subcommands := make([]subcommand, len(scs.cmds))
	{
		i := 0
		for name, cmd := range scs.cmds {
			subcommands[i] = subcommand{Name: name, Command: cmd}
			i += 1
		}

		slices.SortFunc(subcommands, func(a, b subcommand) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}

	for _, c := range subcommands {
		sd := c.Command.ShortDescription()

		if sd == "" {
			fmt.Fprintf(w, "    %s", c.Name)
		} else if name := c.Name; len(name) <= 12 {
			fmt.Fprintf(w, "    %-12s", c.Name)
		} else {
			fmt.Fprintf(w, "    %s\n                ", c.Name)
		}
		help := strings.Split(sd, "\n")
		if 0 < len(help) {
			fmt.Fprintln(w, help[0])
			help = help[1:]
		}
		for _, l := range help {
			fmt.Fprintln(w, "                ", l)
		}
	}

	fmt.Fprintln(w)
}
