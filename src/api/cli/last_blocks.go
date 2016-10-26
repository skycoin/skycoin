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
		Name:        "lastBlocks",
		ArgsUsage:   "Displays the content of the most recently N generated blocks.",
		Usage:       "[numberOfBlocks]",
		Description: "All results returned in JSON format.",
		Action:      getLastBlocks,
	}
	Commands = append(Commands, cmd)
}

func getLastBlocks(c *gcli.Context) error {
	num := c.Args().First()
	if num == "" {
		num = "1"
	}

	n, err := strconv.Atoi(num)
	if err != nil {
		return errors.New("error block number")
	}
	url := fmt.Sprintf("%s/last_blocks?num=%d", nodeAddress, n)
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
