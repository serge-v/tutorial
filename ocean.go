package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"flag"
	"os"
	"time"
	"strings"
	"compress/gzip"
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

func buildProgram(name string) {
	fmt.Println("building", name, "for the ocean")
	cmd := exec.Command("go", "build", name+".go")
	
	env := os.Environ()
	env = append(env, "GOOS=freebsd")
	env = append(env, "GOARCH=amd64")
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("errors:[\n", string(out), "]")
		panic(err)
	}
	fmt.Println("building done")
}

func deployProgram(fname string) {
	session := createSession()

	// this function does the same as the following commands
	// cat file | gzip | ssh host 'zcat > file'

	// create gzip stream
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	pr, pw := io.Pipe()
	gzwriter := gzip.NewWriter(pw)

	// copy to gzip
	go func() {
		written, err := io.Copy(gzwriter, f)
		if err != nil {
			panic(err)
		}
		println(written, "bytes transmitted")
		gzwriter.Close()
		f.Close()
		pw.Close()
	}()

	// connect ssh session input to the gzip output
	session.Stdin = pr

	fmt.Printf("copying %s to the ocean\n", fname)
	cmd := "zcat > " + fname
	if err := session.Run(cmd); err != nil {
		panic(err)
	}

	session.Close()
	session = createSession()
	defer session.Close()

	// run tuning commands
	commands := `
		sudo cp {fname} /usr/local/www/wet/{fname};
		sudo chmod +x /usr/local/www/wet/{fname};
		ls -l /usr/local/www/wet/`

	cmd = strings.Replace(commands, "{fname}", fname, -1)

	fmt.Println("deploying " + fname)
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

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
 
	if *deploy != "" {
		buildProgram(*deploy)
	}

	client = createSshClient()

	if *deploy != "" {
		deployProgram(*deploy)
	} else if *status {
		getStatus()
	} else if *tty {
		runTerminal()
	}
}
