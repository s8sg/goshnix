package goshnix

import "golang.org/x/crypto/ssh"
import "fmt"
import "strconv"
import "strings"
import "bytes"

type ssh_client struct {
	config *ssh.ClientConfig
	addr   string
}

// Create a ssh client for a specific host
func create_ssh_client(host, port, uname, pass string) (*ssh_client, error) {
	sshConfig := &ssh.ClientConfig{
		User: uname,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	client := &ssh_client{}
	client.config = sshConfig
	client.addr = addr
	return client, nil
}

// Parse the response
func parse_response(output string) (string, int) {
	data := strings.Split(output, "\n")
	ept := len(data) - 2
	rt := data[ept]
	if len(data) == 2 {
		output = ""
	} else {
		output = strings.Join(data[0:ept], "\n")
	}
	returncode, _ := strconv.Atoi(rt)
	return output, returncode
}

// Execute a command over ssh
func (client *ssh_client) execute_command(command string) (string, error) {
	connection, err := ssh.Dial("tcp", client.addr, client.config)
	if err != nil {
		//fmt.Printf("Failed to dial: %s", err)
		return "", fmt.Errorf("Failed to dial: %s", err)
	}
	session, err := connection.NewSession()
	if err != nil {
		//fmt.Printf("Failed to initiate session: %s", err)
		return "", fmt.Errorf("Failed to initiate session: %s", err)
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	ncommand := fmt.Sprintf("%s;%s", command, "echo $?")
	err = session.Run(ncommand)
	if err != nil {
		//fmt.Printf("Unable to run command: %v", err)
		return "", fmt.Errorf("Unable to run command: %v", err)
	}
	var reterr error = nil
	errstr := stderrBuf.String()
	outstr := stdoutBuf.String()
	op, returncode := parse_response(outstr)
	if returncode == 0 && op == "" {
		op = errstr
	}
	if returncode != 0 {
		reterr = fmt.Errorf("%s", errstr)
	}
	return op, reterr
}

// Get file content over ssh
func (client *ssh_client) get_file_content(filepath string) (string, error) {
	command := "cat " + filepath
	data, err := client.execute_command(command)
	if err != nil {
		return "", err
	}
	return data, nil
}
