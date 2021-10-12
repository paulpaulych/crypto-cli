package recv

import (
	"flag"
	"fmt"
	. "github.com/paulpaulych/crypto/shamir/cli/cmd/common"
	"github.com/paulpaulych/crypto/shamir/cli/internal/messaging"
	"github.com/paulpaulych/crypto/shamir/cli/internal/messaging/protocols"
	"net"
)

type RecvConf struct{}

func (conf *RecvConf) CmdName() string {
	return "recv"
}

func (conf *RecvConf) InitCmd(args []string) (Cmd, CmdConfError) {
	flags := flag.NewFlagSet(conf.CmdName(), flag.ContinueOnError)

	protocol := flags.String("protocol", "", "protocol")
	bindHostPtr := flags.String("host", "localhost", "host to bind")
	bindPortPtr := flags.String("port", "80", "port to bind")

	err := Parse(flags, args)
	if err != nil {
		return nil, err
	}

	if flags.NArg() < 1 {

	}

	addr := net.JoinHostPort(*bindHostPtr, *bindPortPtr)
	reader, err := readerForProtocol(*protocol)
	if err != nil {
		return nil, err
	}
	return &RecvCmd{bindAddr: addr, reader: reader}, nil
}

type RecvCmd struct {
	bindAddr string
	reader   messaging.ReadMsg
}

func (cmd *RecvCmd) Run() error {
	return messaging.ListenForMsg(cmd.bindAddr, cmd.reader)
}

func readerForProtocol(name string) (messaging.ReadMsg, CmdConfError) {
	switch name {
	case "shamir":
		return protocols.ShamirReader(), nil
	default:
		msg := fmt.Sprintf("unknown protocol '%s'", name)
		return nil, NewCmdConfError(msg, nil)
	}
}
