package awsdial

import (
	"github.com/mmmorris1975/ssm-session-client/datachannel"
	"io"
	"net"
	"time"
)

type ssmconn struct {
	*datachannel.SsmDataChannel
	pr     *io.PipeReader
	target string
}

func (s ssmconn) Read(data []byte) (int, error) {
	return s.pr.Read(data)
}

func (s ssmconn) Close() error {
	s.SsmDataChannel.TerminateSession()
	s.SsmDataChannel.Close()
	s.pr.Close()
	return nil
}

func (s ssmconn) LocalAddr() net.Addr {
	return addrString("")
}

func (s ssmconn) RemoteAddr() net.Addr {
	return addrString(s.target)
}

func (s ssmconn) SetDeadline(t time.Time) error {
	return nil
}

func (s ssmconn) SetReadDeadline(t time.Time) error {
	return nil
}

func (s ssmconn) SetWriteDeadline(t time.Time) error {
	return nil
}

type addrString string

func (a addrString) Network() string {
	return "cmd"
}

func (a addrString) String() string {
	return string(a)
}
