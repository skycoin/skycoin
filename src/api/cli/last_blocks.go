package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	gcli "gopkg.in/urfave/cli.v1"
)

func init() {
	cmd := gcli.Command{
		Name:        "lastBlocks",
		Usage:       "Displays the content of the most recently N generated blocks.",
		Description: "All results returned in JSON format.",
		ArgsUsage:   "[numberOfBlocks]",
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
		return err
	}
	url := fmt.Sprintf("%s/last_blocks?num=%d", nodeAddress, n)
	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	d, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(d))
	return nil
}
