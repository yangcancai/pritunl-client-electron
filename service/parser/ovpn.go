package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Remote struct {
	Host  string
	Port  int
	Proto string
}

type Ovpn struct {
	EnvId             string
	EnvName           string
	Dev               string
	DevType           string
	Remotes           []Remote
	RemoteRandom      bool
	NoBind            bool
	PersisitTun       bool
	Cipher            string
	Auth              string
	Verb              int
	Mute              int
	PushPeerInfo      bool
	Ping              int
	PingExit          int
	HandWindow        int
	ServerPollTimeout int
	RenegSec          int
	SndBuf            int
	RcvBuf            int
	RemoteCertTls     string
	Compress          string
	BlockOutsideDns   bool
	AuthUserPass      bool
	KeyDirection      int
	CaCert            string
	TlsAuth           string
	Cert              string
	Key               string
}

func (o *Ovpn) Export() string {
	output := ""

	if o.EnvId != "" {
		output += fmt.Sprintf("setenv UV_ID %s\n", o.EnvId)
	}
	if o.EnvName != "" {
		output += fmt.Sprintf("setenv UV_NAME %s\n", o.EnvName)
	}
	output += "client\n"
	output += fmt.Sprintf("dev %s\n", o.Dev)
	output += fmt.Sprintf("dev-type %s\n", o.DevType)
	for _, remote := range o.Remotes {
		output += fmt.Sprintf(
			"remote %s %d %s\n",
			remote.Host,
			remote.Port,
			remote.Proto,
		)
	}
	if o.RemoteRandom {
		output += "remote-random\n"
	}
	if o.NoBind {
		output += "nobind\n"
	}
	if o.PersisitTun {
		output += "persist-tun\n"
	}
	if o.Cipher != "" {
		output += fmt.Sprintf("cipher %s\n", o.Cipher)
	}
	if o.Auth != "" {
		output += fmt.Sprintf("auth %s\n", o.Auth)
	}
	if o.Verb > 0 {
		output += fmt.Sprintf("verb %d\n", o.Verb)
	}
	if o.Mute > 0 {
		output += fmt.Sprintf("mute %d\n", o.Mute)
	}
	if o.PushPeerInfo {
		output += "push-peer-info\n"
	}
	if o.Ping > 0 {
		output += fmt.Sprintf("ping %d\n", o.Ping)
	}
	if o.PingExit > 0 {
		output += fmt.Sprintf("ping-exit %d\n", o.PingExit)
	}
	if o.HandWindow > 0 {
		output += fmt.Sprintf("hand-window %d\n", o.HandWindow)
	}
	if o.ServerPollTimeout > 0 {
		output += fmt.Sprintf("server-poll-timeout %d\n", o.ServerPollTimeout)
	}
	if o.RenegSec > 0 {
		output += fmt.Sprintf("reneg-sec %d\n", o.RenegSec)
	}
	if o.SndBuf > 0 {
		output += fmt.Sprintf("sndbuf %d\n", o.SndBuf)
	}
	if o.RcvBuf > 0 {
		output += fmt.Sprintf("rcvbuf %d\n", o.RcvBuf)
	}
	if o.RemoteCertTls != "" {
		output += fmt.Sprintf("remote-cert-tls %s\n", o.RemoteCertTls)
	}
	if o.Compress != "" {
		output += fmt.Sprintf("compress %s\n", o.Compress)
	}
	if o.AuthUserPass {
		output += "auth-user-pass\n"
	}
	if o.KeyDirection > 0 {
		output += fmt.Sprintf("key-direction %d\n", o.KeyDirection)
	}

	if o.CaCert != "" {
		output += fmt.Sprintf("<ca>\n%s</ca>\n", o.CaCert)
	}
	if o.TlsAuth != "" {
		output += fmt.Sprintf("<tls-auth>\n%s</tls-auth>\n", o.TlsAuth)
	}
	if o.Cert != "" {
		output += fmt.Sprintf("<cert>\n%s</cert>\n", o.Cert)
	}
	if o.Key != "" {
		output += fmt.Sprintf("<key>\n%s</key>\n", o.Key)
	}

	return output
}

func Import(data string) (o *Ovpn) {
	o = &Ovpn{
		Remotes: []Remote{},
	}

	inCa := false
	inTls := false
	inCert := false
	inKey := false

	data = strings.ReplaceAll(data, "\r", "")

	for _, line := range strings.Split(data, "\n") {
		line = FilterStr(line, 256)

		if inCa {
			if line == "</ca>" {
				inCa = false
				continue
			}
			o.CaCert += line + "\n"
		} else if inTls {
			if line == "</tls-auth>" {
				inTls = false
				continue
			}
			o.TlsAuth += line + "\n"
		} else if inCert {
			if line == "</cert>" {
				inCert = false
				continue
			}
			o.Cert += line + "\n"
		} else if inKey {
			if line == "</key>" {
				inKey = false
				continue
			}
			o.Key += line + "\n"
		}

		lines := strings.Split(line, " ")

		key := strings.ToLower(lines[0])

		switch key {
		case "<ca>":
			inCa = true
			break
		case "<tls-auth>":
			inTls = true
			break
		case "<cert>":
			inCert = true
			break
		case "<key>":
			inKey = true
			break
		case "setenv":
			if len(lines) != 3 {
				continue
			}
			switch strings.ToLower(lines[1]) {
			case "uv_id":
				o.EnvId = lines[2]
				break
			case "uv_name":
				o.EnvName = lines[2]
				break
			}
			break
		case "dev":
			switch strings.ToLower(lines[1]) {
			case "tun":
				o.Dev = "tun"
				break
			case "tap":
				o.Dev = "tap"
				break
			}
			break
		case "dev-type":
			switch strings.ToLower(lines[1]) {
			case "tun":
				o.Dev = "tun"
				break
			case "tap":
				o.Dev = "tap"
				break
			}
			break
		case "remote":
			if len(lines) != 4 {
				continue
			}

			port, e := strconv.Atoi(lines[2])
			if e != nil {
				continue
			}

			remote := Remote{
				Host: lines[1],
				Port: port,
			}

			switch strings.ToLower(lines[3]) {
			case "udp":
				remote.Proto = "udp"
				break
			case "udp6":
				remote.Proto = "udp6"
				break
			case "tcp":
				remote.Proto = "tcp"
				break
			case "tcp6":
				remote.Proto = "tcp6"
				break
			default:
				continue
			}

			o.Remotes = append(o.Remotes, remote)

			break
		case "remote-random":
			o.RemoteRandom = true
			break
		case "nobind":
			o.NoBind = true
			break
		case "persist-tun":
			o.PersisitTun = true
			break
		case "cipher":
			if len(lines) != 2 {
				continue
			}

			o.Cipher = lines[1]
			break
		case "auth":
			if len(lines) != 2 {
				continue
			}

			o.Auth = lines[1]
			break
		case "verb":
			if len(lines) != 2 {
				continue
			}

			verb, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.Verb = verb
			break
		case "mute":
			if len(lines) != 2 {
				continue
			}

			mute, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.Mute = mute
			break
		case "push-peer-info":
			o.PushPeerInfo = true
			break
		case "ping":
			if len(lines) != 2 {
				continue
			}

			ping, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.Ping = ping
			break
		case "ping-restart":
			if len(lines) != 2 {
				continue
			}

			pingRestart, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.PingExit = pingRestart
			break
		case "ping-exit":
			if len(lines) != 2 {
				continue
			}

			pingExit, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.PingExit = pingExit
			break
		case "hand-window":
			if len(lines) != 2 {
				continue
			}

			handWindow, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.HandWindow = handWindow
			break
		case "server-poll-timeout":
			if len(lines) != 2 {
				continue
			}

			serverPollTimeout, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.ServerPollTimeout = serverPollTimeout
			break
		case "reneg-sec":
			if len(lines) != 2 {
				continue
			}

			renegSec, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.RenegSec = renegSec
			break
		case "sndbuf":
			if len(lines) != 2 {
				continue
			}

			sndbuf, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.SndBuf = sndbuf
			break
		case "rcvbuf":
			if len(lines) != 2 {
				continue
			}

			rcvbuf, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.RcvBuf = rcvbuf
			break
		case "remote-cert-tls":
			if len(lines) != 2 {
				continue
			}

			switch strings.ToLower(lines[1]) {
			case "server":
				o.RemoteCertTls = "server"
				break
			default:
				continue
			}

			break
		case "comp-lzo":
			if len(lines) != 2 {
				continue
			}

			switch strings.ToLower(lines[1]) {
			case "on":
				o.Compress = "lzo"
				break
			default:
				continue
			}

			break
		case "compress":
			if len(lines) != 2 {
				continue
			}

			switch strings.ToLower(lines[1]) {
			case "lzo":
				o.Compress = "lzo"
				break
			case "lz4":
				o.Compress = "lz4"
				break
			default:
				continue
			}

			break
		case "auth-user-pass":
			o.AuthUserPass = true
			break
		case "key-direction":
			if len(lines) != 2 {
				continue
			}

			keyDirection, e := strconv.Atoi(lines[1])
			if e != nil {
				continue
			}

			o.KeyDirection = keyDirection
			break
		}
	}

	if o.Dev == "" {
		o.Dev = "tun"
	}
	if o.DevType == "" {
		o.DevType = "tun"
	}

	return
}