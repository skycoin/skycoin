package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "blocks",
		Usage:     "Lists the content of a single block of a range of blocks. Block results are always in JSON format.",
		ArgsUsage: "[starting block or single block] [ending block]]",
		Action:    getBlocks,
	}
	Commands = append(Commands, cmd)
}

func getBlocks(c *gcli.Context) error {
	// get start
	start := c.Args().Get(0)
	end := c.Args().Get(1)
	if end == "" {
		end = start
	}

	s, err := strconv.Atoi(start)
	if err != nil {
		return err
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:6420/blocks?start=%d&end=%d", s, e)
	// blocks := visor.ReadableBlocks{}

	rsp, err := http.Get(url)
	if err != nil {
		return errConnectNodeFailed
	}
	defer rsp.Body.Close()
	d, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(d))
	return nil
}
