package goshnix

import (
	"fmt"
	"testing"
)

var (
	host     = "127.0.0.1"
	port     = "22"
	uname    = "swarvanu"
	pass     = "swarvanu"
	passcmd  = "echo \"Hello\""
	passresp = "Hello"
	failcmd  = "swarvanu"
	testfile = "test.txt"
)

func TestSshClientCreation(t *testing.T) {
	client, err := create_ssh_client(host, port, uname, pass)
	if err != nil || client == nil {
		t.Errorf("Failed to create ssh client: %s", err)
	}
	fmt.Println("SSH Client: ", client)
}

func TestExecuteCommand(t *testing.T) {
	client, err := create_ssh_client(host, port, uname, pass)
	if err != nil || client == nil {
		t.Errorf("Failed to create ssh client: %s", err)
	}
	op, excerr := client.execute_command(passcmd)
	if excerr != nil {
		t.Errorf("Failed to execute command: %s", excerr)
	}
	if op != passresp {
		t.Errorf("Failed to execute command, op is different than expected")
	}
	fmt.Printf("%s :%s\n", passcmd, op)

	op, excerr = client.execute_command(failcmd)
	if excerr == nil {
		t.Errorf("No command exist like swarvanu, op: %s", op)
	}
	fmt.Printf("The error string:%s\n", excerr)
}

func TestGetFileContent(t *testing.T) {
	client, err := create_ssh_client(host, port, uname, pass)
	if err != nil || client == nil {
		t.Errorf("Failed to create ssh client: %s", err)
	}
	op, excerr := client.get_file_content(testfile)
	if excerr != nil {
		t.Errorf("Failed to get file content: %s", excerr)
	}
	if op == "" {
		t.Errorf("Failed to get file content, op is blank")
	}
	fmt.Printf("%s content:%s\n", testfile, op)
}
