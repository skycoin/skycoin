package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "blocks",
		ArgsUsage: "Lists the content of a single block of a range of blocks. Block results are always in JSON format.",
		Usage:     "[starting block or single block seq] [ending block seq]",
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
		return errors.New("error block seq")
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		return errors.New("error block seq")
	}

	url := fmt.Sprintf("http://localhost:6420/blocks?start=%d&end=%d", s, e)
	// blocks := visor.ReadableBlocks{}

	rsp, err := http.Get(url)
	if err != nil {
		return errConnectNodeFailed
	}
	defer rsp.Body.Close()
	d, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errReadResponse
	}

	fmt.Println(string(d))
	return nil
}
