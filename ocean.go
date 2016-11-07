// ocean is a tool to deploy golang application to digitalocean cloud server.
// It builds the app for freebsd-amd64 and copies to /usr/local/www/wet virtual directory.
// Uses ocean-params.go file which contains credentials to connect to the ocean.
package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"flag"
	"os"
	"time"
	"strings"
	"os/exec"
)

var (
	deploy = flag.String("deploy", "", "name of the program to deploy to ocean server. Should have corresponding .go source file.")
	status = flag.Bool("status", false, "show ocean server status")
	tty = flag.Bool("tty", false, "open terminal")

	client *ssh.Client
)

type sshParams struct {
	PrivateKey string
	Host string // "127.0.0.1:22"
	User string
}

//var params = sshParams{}
var params = oceanParams // declared in separate file not included to git

func createSshClient() *ssh.Client {
	signer, err := ssh.ParsePrivateKey([]byte(params.PrivateKey))
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: params.User,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: time.Second * 3,
	}

	client, err := ssh.Dial("tcp", params.Host, config)
	if err != nil {
		panic(err)
	}
	return client
}

func createSession() *ssh.Session {
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	session.Stdout = os.Stdout
	session.Stdout = os.Stderr

	return session
}

func deployFile(fname string) {
	session := createSession()

	// this function does the same as the following commands
	// cat file | ssh host 'cat > file'

	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	// connect ssh session input to the opened file
	session.Stdin = f

	fmt.Printf("copying %s to the ocean\n", fname)
	cmd := "cat > " + fname
	if err := session.Run(cmd); err != nil {
		panic(err)
	}

	session.Close()
	session = createSession()
	defer session.Close()

	// run remote deploy commands
	fmt.Println("deploying " + fname)
	cmd = "sudo tar -C / -xzf " + fname + "; sudo /usr/local/lib/" + fname + "-configure.sh"
	if err := session.Run(cmd); err != nil {
		panic(err)
	}
	fmt.Println("done")
}

func runTerminal() {
	session := createSession()
	defer session.Close()

	in, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		panic(err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		panic(err)
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

func getStatus() {
	session := createSession()
	defer session.Close()
	if err := session.Run("ls -l /usr/local/www/wet/"); err != nil {
		panic(err)
	}
}

func execInstallRules(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "// ocean-build ") {
			pars := strings.SplitN(s, " ", 3)
			if len(pars) != 3 {
				panic("invalid ocean-build command " + s)
			}

			cmd := exec.Command(pars[2])
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("errors:[\n", string(out), "]")
				panic(err)
			}
			fmt.Println(string(out))
		} else if strings.HasPrefix(s, "// ocean-deploy ") {
			pars := strings.SplitN(s, " ", 3)
			if len(pars) != 3 {
				panic("invalid ocean-build command " + s)
			}
			deployFile(pars[2])
			fmt.Println("deployment done")
			break
		}
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
 
	client = createSshClient()
	defer client.Close()

	if *deploy != "" {
		execInstallRules(*deploy)
	} else if *status {
		getStatus()
	} else if *tty {
		runTerminal()
	}
}
