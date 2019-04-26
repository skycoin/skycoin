package bip32

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// Path represents a parsed HD wallet path
type Path struct {
	Elements []PathNode
}

// PathNode is an element of an HD wallet path
type PathNode struct {
	Master      bool
	ChildNumber uint32
}

// Hardened returns true if this path node is hardened
func (p PathNode) Hardened() bool {
	return p.ChildNumber >= FirstHardenedChild
}

var (
	// ErrPathNoMaster HD wallet path does not start with m
	ErrPathNoMaster = errors.New("Path must start with m")
	// ErrPathChildMaster HD wallet path contains m in a child node
	ErrPathChildMaster = errors.New("Path contains m as a child node")
	// ErrPathNodeNotNumber HD wallet path node is not a valid uint32 number
	ErrPathNodeNotNumber = errors.New("Path node is not a valid uint32 number")
	// ErrPathNodeNumberTooLarge HD wallet path node is >= 2^31
	ErrPathNodeNumberTooLarge = errors.New("Path node must be less than 2^31")
)

// ParsePath parses a bip32 HD wallet path. The path must start with m/.
func ParsePath(p string) (*Path, error) {
	pts := strings.Split(p, "/")

	path := &Path{
		Elements: []PathNode{
			{
				Master:      true,
				ChildNumber: 0,
			},
		},
	}

	for i, x := range pts {
		if i == 0 {
			if x != "m" {
				return nil, ErrPathNoMaster
			}
			// Path.Elements already initialized with master node
			continue
		} else if x == "m" {
			return nil, ErrPathChildMaster
		}

		n, err := parseNode(x)
		if err != nil {
			return nil, err
		}

		path.Elements = append(path.Elements, n)
	}

	return path, nil
}

func parseNode(x string) (PathNode, error) {
	// Hardened nodes have an apostrophe ' appended
	hardened := false
	if strings.HasSuffix(x, "'") {
		hardened = true
		x = x[:len(x)-1]
	}

	// Node element (minus a single trailing apostrophe) must be a valid uint32 number
	n, err := strconv.ParseUint(x, 10, 32)
	if err != nil {
		return PathNode{}, ErrPathNodeNotNumber
	}

	// Node number must be <2^31. If hardened, 2^31 will be added to it
	if n >= uint64(FirstHardenedChild) {
		return PathNode{}, ErrPathNodeNumberTooLarge
	}

	nn := uint32(n)

	// Add 2^31 to the base number for hardened nodes
	if hardened {
		nnn := nn + FirstHardenedChild

		// Sanity check
		if nnn < nn {
			log.Panic("Unexpected overflow of path number when adjusting hardened child number")
		}
		nn = nnn
	}

	return PathNode{
		Master:      false,
		ChildNumber: nn,
	}, nil
}
