//Inspired by: https://github.com/42wim/cssh/blob/master/device/cisco.go
package brocadevdx

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"golang.org/x/crypto/ssh/agent"
	"fmt"
	"strings"
	"io"
	"bufio"
	"time"
	"log"
)

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func ConnectToSwitchWithPassword(address string, username string, password string) (*VDXSwitch, error) {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	//Connect to switch
	return SSHToSwitch(address, sshConfig)
}

func ConnectToSwitchWithSSHAgent(address string, username string) (*VDXSwitch, error) {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	//Connect to switch
	return SSHToSwitch(address, sshConfig)
}

func ConnectToSwitchWithUnencryptedCertificate(address string, username string, pathToKey string) (*VDXSwitch, error) {
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(pathToKey)
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	//Connect to switch
	return SSHToSwitch(address, sshConfig)
}

func SSHToSwitch(address string, sshConfig *ssh.ClientConfig) (*VDXSwitch, error) {
	var device VDXSwitch
	//Set fields
	device.Hostname = address
	device.Timeout = 5

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	connection, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err.Error())
	}

	//Get a session
	session, err := connection.NewSession()
	if err != nil {
		connection.Conn.Close()
		return nil, err
	}

	//Init connection management channels
	device.ReadChan = make(chan *string, 20)
	device.StopChan = make(chan struct{})

	//Make terminal config
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	//Get a PTY
	if err := session.RequestPty("xterm", 8000, 40, modes); err != nil {
		session.Close()
		return nil, err
	}

	//Redirect output and input
	device.stdin, err = session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	device.stdout, err = session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}

	//Get a shell
	session.Shell()

	//Save our values
	device.client = connection
	device.session = session

	//Wait for prompt
	bufstdout := bufio.NewReader(device.stdout)
	buf := make([]byte, 1000)
	loadStr := ""
	for {
		n, err := bufstdout.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		loadStr += string(buf[:n])
		//If the line has a # then we are at the prompt
		if strings.Contains(loadStr, "#") {
			break
		}
	}

	return &device, nil
}

func (d *VDXSwitch) readln(r io.Reader) {
	//Make a buffer for our output
	buf := make([]byte, 10000)
	loadStr := ""
	for {
		//Read in output byte by byte and build a string
		n, err := r.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("ERROR ", err)
			}
			d.StopChan <- struct{}{}
		}
		//Add to our output string
		loadStr += string(buf[:n])

		//If the output has '# ' then it is the prompt and we should break
		if strings.HasSuffix(loadStr, "# ") {
			break
		}
		// keepalive
		d.ReadChan <- nil
	}
	//Signal when we are done
	//loadStr = strings.Replace(loadStr, "\r", "", -1)
	d.ReadChan <- &loadStr
}

func (d *VDXSwitch) Cmd(cmd string) (string, error) {
	var result string
	bufstdout := bufio.NewReader(d.stdout)
	lines := strings.Split(cmd, "\n")
	for _, line := range lines {
		io.WriteString(d.stdin, line+"\n")
		time.Sleep(time.Millisecond * 100)
	}
	go d.readln(bufstdout)
	for {
		select {
		case output := <-d.ReadChan:
			{
				if output == nil {
					continue
				}
				result = strings.Replace(*output, lines[0], "", 1)
				return result, nil
			}
		case <-d.StopChan:
			{
				if d.session != nil {
					d.session.Close()
				}
				d.client.Conn.Close()
				return "", fmt.Errorf("EOF")
			}
		case <-time.After(time.Second * time.Duration(d.Timeout)):
			{
				fmt.Println("timeout on", d.Hostname)
				if d.session != nil {
					d.session.Close()
				}
				d.client.Conn.Close()
				return "", nil
			}
		}
	}
}