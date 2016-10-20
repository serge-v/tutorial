package main

import (
	"bytes"
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"flag"
	"os"
	"time"
	"strings"
)

var (
	deploy = flag.String("deploy", "", "file name to deploy to ocean server")
	status = flag.Bool("status", false, "show ocean server status")
	tty = flag.Bool("tty", false, "open terminal")
)

type sshParams struct {
	PrivateKey string
	Host string // "127.0.0.1:22"
	User string
}

//var params = sshParams{}
var params = oceanParams // declared in separate file not included to git

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	signer, err := ssh.ParsePrivateKey([]byte(params.PrivateKey))
	if err != nil {
		log.Fatal("parse key failed: ", err)
	}

	config := &ssh.ClientConfig{
		User: params.User,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: time.Second * 3,
	}

	client, err := ssh.Dial("tcp", params.Host, config)
	if err != nil {
		log.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if *deploy != "" {
		fname := *deploy
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal("cannot open file: ", err)
		}
		session.Stdin = f
		
		fmt.Printf("copying %s to the ocean\n", fname)
		cmd := "cat > " + fname
		if err := session.Run(cmd); err != nil {
			log.Fatal("failed to run: ", err)
		}

		session.Close()

		session, err = client.NewSession()
		if err != nil {
			log.Fatal("failed to reopen session: ", err)
		}

		cmd = strings.Replace("sudo cp {fname} /usr/local/www/wet/{fname};" +
			"sudo chmod +x /usr/local/www/wet/{fname};" +
			"ls -l /usr/local/www/wet/", "{fname}", fname, -1)

		fmt.Println("deploying " + fname)
		if err := session.Run(cmd); err != nil {
			log.Fatal("failed to run:", err)
		}
		fmt.Println("done")

	} else if *status {
		if err := session.Run("ls -l /usr/local/www/wet/"); err != nil {
			log.Fatal("failed to run:", err)
		}
		fmt.Println(b.String())
	} else if *tty {
	
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		in, _ := session.StdinPipe()
	
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		
		// Request pseudo terminal
		if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
			log.Fatal("request for pseudo terminal failed: ", err)
		}

		// Start remote shell
		if err := session.Shell(); err != nil {
			log.Fatal("failed to start shell: ", err)
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			str, err := reader.ReadString('\n')
			if err != nil {
				println(err)
				break
			}
			fmt.Fprint(in, str)
		}
	}
}
