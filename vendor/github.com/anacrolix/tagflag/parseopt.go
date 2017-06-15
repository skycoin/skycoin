package tagflag

type parseOpt func(p *parser)

// Don't perform default behaviour if -h or -help are passed.
func NoDefaultHelp() parseOpt {
	return func(p *parser) {
		p.noDefaultHelp = true
	}
}

// Provides a description for the program to be shown in the usage message.
func Description(desc string) parseOpt {
	return func(p *parser) {
		p.description = desc
	}
}

func Program(name string) parseOpt {
	return func(p *parser) {
		p.program = name
	}
}
