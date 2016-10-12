package cli

import gcli "gopkg.in/urfave/cli.v1"

func init() {
	cmd := gcli.Command{
		Name:      "generateWallet",
		Usage:     "Generate a new wallet from seed.",
		ArgsUsage: "[options]",
		Description: `
		Use caution when using the “-p” command. If you have command history enabled your 
		wallet encryption password can be recovered from the history log. If you do not 
		include the “-p” option you will be prompted to enter your password after you enter 
		your command. 
		
		All results are returned in JSON format. 
                      `,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "s",
				Usage: "Your seed.",
			},
			gcli.StringFlag{
				Name:  "r",
				Usage: "A random alpha numeric seed will be generated for you.",
			},
			gcli.StringFlag{
				Name:  "rd",
				Usage: "A random seed consisting of 12 dictionary words will be generated for you.",
			},
			gcli.IntFlag{
				Name:  "m",
				Usage: "[numberOfAddresses] Number of addresses to generate. By default 1 address is generated.",
			},
			gcli.StringFlag{
				Name:  "p",
				Usage: "Password used to encrypt the wallet locally.",
			},
			gcli.StringFlag{
				Name:  "n",
				Usage: `[walletName] Name of wallet. The final format will be "yourName.wlt". If no wallet name is specified a generic name will be selected.`,
			},
			gcli.StringFlag{
				Name:  "l",
				Usage: "[label] Label used to idetify your wallet.",
			},
		},
		Action: func(c *gcli.Context) error {
			return nil
		},
	}
	Commands = append(Commands, cmd)
}
