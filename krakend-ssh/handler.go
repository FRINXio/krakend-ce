package ssh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/ssh"
)

var sshCli = []*sshClient{}

const (
	messageWait    = 10 * time.Second
	maxMessageSize = 2048
)

var terminalModes = ssh.TerminalModes{
	ssh.ECHO:          1,     // enable echoing (different from the example in docs)
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

type sshClient struct {
	host     string
	addr     string
	user     string
	secret   string
	client   *ssh.Client
	sess     *ssh.Session
	sessIn   io.WriteCloser
	sessOut  io.Reader
	closeSig chan struct{}
	data     []byte
}

func (c *sshClient) wsWrite() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		_, readErr := c.sessOut.Read(c.data)
		fmt.Printf(string(c.data))

		if readErr != nil {
			return fmt.Errorf("sessOut.Read: %w", readErr)
		}
	}
}

func (c *sshClient) wsRead() error {
	defer func() {
		c.closeSig <- struct{}{}
	}()

	src := strings.NewReader("")
	_, err := io.Copy(c.sessIn, src)
	for {

	}
	return fmt.Errorf("conn.NextReader: %w", err)

}

func (c *sshClient) wsRead2(cmd string) {
	src := strings.NewReader(cmd + "\n")
	_, err := io.Copy(c.sessIn, src)

	if err != nil {
		log.Println("bridgeWSAndSSH: ssh.Dial:", err)
	}
}

func (c *sshClient) sshConnection() {

	defer panic("Exited session")

	err := *new(error)
	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.secret),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	c.client, err = ssh.Dial("tcp", c.addr, config)
	if err != nil {
		log.Println("bridgeWSAndSSH: ssh.Dial:", err)
		return
	}
	defer c.client.Close()

	c.sess, err = c.client.NewSession()
	if err != nil {
		log.Println("bridgeWSAndSSH: client.NewSession:", err)
		return
	}
	defer c.sess.Close()

	c.sess.Stderr = os.Stderr // TODO: check proper Stderr output
	c.sessOut, err = c.sess.StdoutPipe()
	if err != nil {
		log.Println("bridgeWSAndSSH: session.StdoutPipe:", err)
		return
	}

	c.sessIn, err = c.sess.StdinPipe()
	if err != nil {
		log.Println("bridgeWSAndSSH: session.StdinPipe:", err)
		return
	}
	defer c.sessIn.Close()

	if err := c.sess.RequestPty("xterm", 20, 100, terminalModes); err != nil {
		log.Println("bridgeWSAndSSH: session.RequestPty:", err)
		return
	}

	if err := c.sess.Shell(); err != nil {
		log.Println("bridgeWSAndSSH: session.Shell:", err)
		return
	}

	log.Println("started a login shell on the remote host")
	defer log.Println("closed a login shell on the remote host")

	go func() {
		if err := c.wsWrite(); err != nil {
			log.Println("bridgeWSAndSSH: wsWrite:", err)
		}
	}()

	<-c.closeSig
}

type sshHandler struct {
	host   string
	addr   string
	port   string
	user   string
	secret string
}

func (c *sshClient) handleExitGoRoutine(_sshCli sshClient, f func()) {
	defer func() {
		if v := recover(); v != nil {
			// A panic is detected.
			log.Println("SSL client exited. Restart it now.")
			go _sshCli.handleExitGoRoutine(_sshCli, f) // restart
		}
	}()
	f()
}

func (h *sshHandler) sshClient() {

	_addr := h.addr + ":" + h.port
	_sshCli := &sshClient{
		host:     h.host,
		addr:     _addr,
		user:     h.user,
		secret:   h.secret,
		closeSig: make(chan struct{}, 1),
		data:     make([]byte, maxMessageSize, maxMessageSize),
	}
	sshCli = append(sshCli, _sshCli)
	go _sshCli.handleExitGoRoutine(*_sshCli, _sshCli.sshConnection)
}

type test_struct struct {
	Test string
}

func sendCommand(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["uniconfig"]

	for i, s := range sshCli {
		if s.host == name {
			fmt.Println(i, s.host)
			decoder := json.NewDecoder(r.Body)
			var t test_struct
			err := decoder.Decode(&t)
			if err != nil {
				panic(err)
			}
			log.Println(t.Test)
			s.wsRead2(t.Test)
			w.Write([]byte("OK"))
			return
		}
	}
	w.Write([]byte("Failed"))
}

func getOutput(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["uniconfig"]

	for _, s := range sshCli {
		if s.host == name {
			fmt.Sprintf("%s", string(s.data))
			log.Println(s)
			w.Write([]byte(string(s.data)))
			return
		}
	}
	w.Write([]byte("Failed"))
}
