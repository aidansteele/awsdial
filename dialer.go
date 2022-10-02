package awsdial

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mmmorris1975/ssm-session-client/datachannel"
	"io"
	"net"
	"strconv"
)

type Dialer struct {
	Client *ssm.Client
}

func (d *Dialer) Dial(ctx context.Context, target string, port int) (net.Conn, error) {
	in := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartSSHSession"),
		Target:       aws.String(target),
		Parameters: map[string][]string{
			"portNumber": {strconv.Itoa(port)},
		},
	}

	start, err := d.Client.StartSession(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("calling StartSession API: %w", err)
	}

	c := &datachannel.SsmDataChannel{}
	err = c.StartSessionFromDataChannelURL(*start.StreamUrl, *start.TokenValue)
	if err != nil {
		return nil, fmt.Errorf("opening ssm datachannel: %w", err)
	}

	err = c.WaitForHandshakeComplete()
	if err != nil {
		return nil, fmt.Errorf("waiting for ssm handshake: %w", err)
	}

	pr, pw := io.Pipe()
	go c.WriteTo(pw)

	conn := ssmconn{SsmDataChannel: c, pr: pr, target: target}
	return conn, nil
}
