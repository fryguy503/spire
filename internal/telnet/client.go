package telnet

import (
	"fmt"
	"github.com/Akkadius/spire/internal/env"
	"github.com/k0kubun/pp/v3"
	"github.com/sirupsen/logrus"
	"github.com/ziutek/telnet"
	"io"
	"net"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Client struct {
	debugging bool
	t         *telnet.Conn
	logger    *logrus.Logger
}

func NewClient(logger *logrus.Logger) *Client {
	return &Client{
		debugging: env.GetInt("DEBUG", "0") >= 3,
		logger:    logger,
	}
}

const (
	linebreak = "\n\r> "
)

func expect(t *telnet.Conn, d ...string) bool {
	err := t.SkipUntil(d...)
	if err != nil {
		return false
	}

	return true
}

func sendln(t *telnet.Conn, s string) error {
	defer func() {
		if r := recover(); r != nil {
			t.Close()
		}
	}()

	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	return err
}

func connCheck(conn net.Conn) error {
	var sysErr error = nil
	rc, err := conn.(syscall.Conn).SyscallConn()
	if err != nil {
		return err
	}
	err = rc.Read(func(fd uintptr) bool {
		var buf []byte = []byte{0}
		n, _, err := syscall.Recvfrom(int(fd), buf, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	})
	if err != nil {
		return err
	}

	return sysErr
}

func (c *Client) Connect() error {
	var err error

	if c.t != nil {
		err := connCheck(c.t.Conn)
		if err != nil {
			c.Close()
			c.t = nil
		}
	}

	if c.t != nil {
		return nil
	}

	d := 1000 * time.Second
	c.t, err = telnet.DialTimeout("tcp", "localhost:9000", d)
	if err != nil {
		return err
	}

	err = c.t.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}
	err = c.t.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}

	err = c.t.SetEcho(false)
	if err != nil {
		return err
	}

	// what the console expects when connecting locally
	if expect(c.t, "assuming admin") {
		c.debug("\n###################################\n# Logging into World\n###################################")

		expect(c.t, ">")
		_ = sendln(c.t, "echo off")
		expect(c.t, ">")
		_ = sendln(c.t, "acceptmessages on")
		expect(c.t, ">")
	}

	return nil
}

func (c *Client) Command(cmd string) (string, error) {
	var err error

	err = c.Connect()
	if err != nil {
		return "", err
	}

	err = c.t.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return "", err
	}
	err = c.t.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return "", err
	}

	sendln(c.t, cmd)

	defer func() {
		if r := recover(); r != nil {
			c.debug("Panic in read, close connection")
			c.Close()
		}
	}()

	var data []byte
	var output string
	for {
		start := time.Now()
		data, err = c.t.ReadUntil(linebreak)
		c.debug("Read operation took %v", time.Since(start))
		if err != nil {
			c.logger.Warnf("[telnet] read failed: %s", err)
			return "", err
		}

		output += string(data)

		if strings.Contains(output, linebreak) {
			output = strings.Replace(output, linebreak, "", 1)
			c.debug("[Output] %v", output)
			return output, nil
		}
	}
}

func (c *Client) Close() {
	if c.t != nil {
		err := c.t.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}
}

func (c *Client) debug(msg string, a ...interface{}) {
	if c.debugging {
		_, file, _, ok := runtime.Caller(1)
		if ok {
			file = filepath.Base(file)
			if len(a) > 0 {
				pp.Printf(fmt.Sprintf("[%v] ", file) + fmt.Sprintf(msg, a...) + "\n")
				return
			}
			pp.Printf(fmt.Sprintf("[%v] ", file) + fmt.Sprintf(msg) + "\n")
		}
	}
}