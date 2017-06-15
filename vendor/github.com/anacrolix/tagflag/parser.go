package tagflag

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/huandu/xstrings"
)

type parser struct {
	// The value from which the parser is built, and values are assigned.
	cmd interface{}
	// Disables the default handling of -h and -help.
	noDefaultHelp bool
	program       string
	description   string

	posArgs []arg
	// Maps -K=V to map[K]arg(V)
	flags map[string]arg

	// Count of positional arguments parsed so far. Used to locate the next
	// positional argument where it's non-trivial (non-unity arity).
	numPos int
}

func (p *parser) hasOptions() bool {
	return len(p.flags) != 0
}

func (p *parser) parse(args []string) (err error) {
	args, err = p.parseAny(args)
	if err != nil {
		return
	}
	err = p.parsePosArgs(args)
	if err != nil {
		return
	}
	if p.numPos < p.minPos() {
		return userError{fmt.Sprintf("missing argument: %q", p.indexPosArg(p.numPos).name)}
	}
	return
}

func (p *parser) minPos() (min int) {
	for _, arg := range p.posArgs {
		min += arg.arity.min
	}
	return
}

func newParser(cmd interface{}, opts ...parseOpt) (p *parser, err error) {
	p = &parser{
		cmd: cmd,
	}
	for _, opt := range opts {
		opt(p)
	}
	err = p.parseCmd()
	return
}

func (p *parser) parseCmd() error {
	if p.cmd == nil {
		return nil
	}
	s := reflect.ValueOf(p.cmd).Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct got %s", s.Type())
	}
	return p.parseStruct(reflect.ValueOf(p.cmd).Elem())
}

// Positional arguments are marked per struct.
func (p *parser) parseStruct(st reflect.Value) (err error) {
	posStarted := false
	foreachStructField(st, func(f reflect.Value, sf reflect.StructField) (stop bool) {
		if !posStarted && f.Type() == reflect.TypeOf(StartPos{}) {
			posStarted = true
			return false
		}
		if sf.PkgPath != "" {
			return false
		}
		if canMarshal(f) {
			if posStarted {
				err = p.addPos(f, sf)
			} else {
				err = p.addFlag(f, sf)
				if err != nil {
					err = fmt.Errorf("error adding flag in %s: %s", st.Type(), err)
				}
			}
			return err != nil
		}
		if f.Kind() == reflect.Struct {
			if canMarshal(f.Addr()) {
				err = fmt.Errorf("field %q has type %s, did you mean to use %s?", sf.Name, f.Type(), f.Addr().Type())
				return true
			}
			err = p.parseStruct(f)
			return err != nil
		}
		err = fmt.Errorf("field has bad type: %v", f.Type())
		return true
	})
	return
}

func newArg(v reflect.Value, sf reflect.StructField, name string) arg {
	return arg{
		arity: fieldArity(v, sf),
		value: v,
		name:  name,
		help:  sf.Tag.Get("help"),
	}
}

func (p *parser) addPos(f reflect.Value, sf reflect.StructField) error {
	p.posArgs = append(p.posArgs, newArg(f, sf, strings.ToUpper(xstrings.ToSnakeCase(sf.Name))))
	return nil
}

func (p *parser) addFlag(f reflect.Value, sf reflect.StructField) error {
	name := structFieldFlag(sf)
	if _, ok := p.flags[name]; ok {
		return fmt.Errorf("flag %q defined more than once", name)
	}
	if p.flags == nil {
		p.flags = make(map[string]arg)
	}
	p.flags[name] = newArg(f, sf, name)
	return nil
}

func (p *parser) parseAny(args []string) (left []string, err error) {
	for len(args) != 0 {
		a := args[0]
		args = args[1:]
		if a == "--" {
			left = args[1:]
			return
		}
		if strings.HasPrefix(a, "-") && len(a) > 1 {
			err = p.parseFlag(a[1:])
		} else {
			err = p.parsePos(a)
		}
		if err != nil {
			break
		}
	}
	return
}

func (p *parser) parsePosArgs(args []string) (err error) {
	for _, a := range args {
		err = p.parsePos(a)
		if err != nil {
			break
		}
	}
	return
}

func (p *parser) parseFlag(s string) error {
	i := strings.IndexByte(s, '=')
	k := s
	v := ""
	if i != -1 {
		k = s[:i]
		v = s[i+1:]
	}
	flag, ok := p.flags[k]
	if !ok {
		if (k == "help" || k == "h") && !p.noDefaultHelp {
			return ErrDefaultHelp
		}
		return userError{fmt.Sprintf("unknown flag: %q", k)}
	}
	err := flag.marshal(v, i != -1)
	if err != nil {
		return fmt.Errorf("error setting flag %q: %s", k, err)
	}
	return nil
}

func (p *parser) indexPosArg(i int) *arg {
	for _, arg := range p.posArgs {
		if i < arg.arity.max {
			return &arg
		}
		i -= arg.arity.max
	}
	return nil
}

func (p *parser) parsePos(s string) (err error) {
	arg := p.indexPosArg(p.numPos)
	if arg == nil {
		return userError{fmt.Sprintf("excess argument: %q", s)}
	}
	err = arg.marshal(s, true)
	if err != nil {
		return
	}
	p.numPos++
	return
}

func structFieldFlag(sf reflect.StructField) string {
	name := sf.Tag.Get("name")
	if name != "" {
		return name
	}
	return fieldFlagName(sf.Name)
}

// Gets the reflect.Value for the nth positional argument.
func posIndexValue(v reflect.Value, _i int) (ret reflect.Value, i int) {
	i = _i
	log.Println("posIndexValue", v.Type(), i)
	switch v.Kind() {
	case reflect.Ptr:
		return posIndexValue(v.Elem(), i)
	case reflect.Struct:
		posStarted := false
		foreachStructField(v, func(fv reflect.Value, sf reflect.StructField) bool {
			log.Println("posIndexValue struct field", fv, sf)
			if !posStarted {
				if fv.Type() == reflect.TypeOf(StartPos{}) {
					// log.Println("posStarted")
					posStarted = true
				}
				return true
			}
			ret, i = posIndexValue(fv, i)
			if ret.IsValid() {
				return false
			}
			return true
		})
		return
	case reflect.Slice:
		ret = v
		return
	default:
		if i == 0 {
			ret = v
			return
		}
		i--
		return
	}
}

func (p *parser) posWithHelp() (ret []arg) {
	for _, a := range p.posArgs {
		if a.help != "" {
			ret = append(ret, a)
		}
	}
	return
}
