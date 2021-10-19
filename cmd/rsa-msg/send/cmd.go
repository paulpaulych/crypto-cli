package send

import (
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/infra/cli"
	"io"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "send"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"bob-pub": "path to file containing destination public key",
		"i":       "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr := flags.Args[0]

	bobPubFName := flags.Flags["bob-pub"].Get()
	input := flags.Flags["i"].GetOr("console")

	msgReader, e := cli.NewInputReader(input, flags.Args[1:])
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	if bobPubFName == nil {
		return nil, cli.NewCmdConfError("required flag: -bob-pub", nil)
	}
	return &Cmd{addr: addr, writer: protocols.RsaWriter(*bobPubFName), msg: msgReader}, nil
}

type Cmd struct {
	addr   string
	writer msg_core.ConnWriter
	msg    io.Reader
}

func (cmd *Cmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, cmd.writer)
}