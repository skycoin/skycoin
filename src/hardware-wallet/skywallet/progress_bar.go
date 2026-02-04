package skywallet

import "fmt"

const (
	maxbars int = 100
)

// Progbar progress bar for cli command in the style:
type Progbar struct {
	total int
}

// PrintProg print the progress var for the portion value
func (p *Progbar) PrintProg(portion int) {
	bars := p.calcBars(portion)
	spaces := maxbars - bars - 1
	percent := 100 * (float32(portion) / float32(p.total))
	fmt.Print("\r[")
	for i := 0; i < bars; i++ {
		fmt.Print("=")
	}
	fmt.Print(">")
	for i := 0; i <= spaces; i++ {
		fmt.Print(" ")
	}
	fmt.Printf(" ] %3.2f%% (%d/%d)", percent, portion, p.total)
}

// PrintComplete print the progress bar as completed
func (p *Progbar) PrintComplete() {
	p.PrintProg(p.total)
	fmt.Print("\n")
}

func (p *Progbar) calcBars(portion int) int {
	if portion == 0 {
		return portion
	}
	return int(float32(maxbars) / (float32(p.total) / float32(portion)))
}
