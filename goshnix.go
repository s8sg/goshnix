package goshnix

import "fmt"
import "strings"
import "os"
import "time"
import "strconv"

// The Goshnix lib client
type Goshnix struct {
	client *ssh_client
}

type fileinfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modtime time.Time
	isdir   bool
}

//Initialize the goshnix client
func Init(host, port, uname, pass string) (*Goshnix, error) {
	client, err := create_sshclient(host, port, uname, pass)
	if err != nil {
		return nil, err
	}
	op, excerr := client.execute_command("echo \"Hello\"")
	if excerr != nil {
		return nil, excerr
	}
	if op != "Hello" {
		return nil, throw_connerror("Test command execution invalid output")
	}
	goshnix := &Goshnix{}
	goshnix.client = client

	return goshnix, nil
}

//Isssherror check if error is due to any ssh connection error
func Isssherror(err error) bool {
	_, ok := err.(connerror)
	return ok
}

// Iscmderror check if error is due to command execution failed
func Iscmderror(err error) bool {
	_, ok := err.(Cmderror)
	return ok
}

//Chmod changes the mode of the named file to mode. If the file is a symbolic link, it changes the mode of the link's target
func (goshnix *Goshnix) Chmod(name string, mode os.FileMode) error {
	modei := int(mode)
	command := fmt.Sprintf("chmod %d %s", modei, name)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

//Chown changes the numeric uid and gid of the named file. If the file is a symbolic link, it changes the uid and gid of the link's target
func (goshnix *Goshnix) Chown(name string, uid, gid int) error {
	command := fmt.Sprintf("chown %d:%d %s", uid, gid, name)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

//Environ returns a copy of strings representing the environment, in the form "key=value"
func (goshnix *Goshnix) Environ() ([]string, error) {
	command := "env"
	op, err := goshnix.client.execute_command(command)
	if err != nil {
		return nil, err
	}
	envs := strings.Split(op, "\n")
	return envs, nil
}

//Getenv retrieves the value of the environment variable named by the key. It returns the value, which will be empty if the variable is not present
func (goshnix *Goshnix) Getenv(key string) (string, error) {
	envs, err := goshnix.Environ()
	if err != nil {
		return "", err
	}
	for _, env := range envs {
		keyval := strings.Split(env, "=")
		if keyval[0] == key {
			return keyval[1], nil
		}
	}
	return "", nil
}

//LookupEnv retrieves the value of the environment variable named by the key. If the variable is present in the environment the value (which may be empty) is returned and the boolean is true
func (goshnix *Goshnix) LookupEnv(key string) (string, bool, error) {
	envs, err := goshnix.Environ()
	if err != nil {
		return "", false, err
	}
	for _, env := range envs {
		keyval := strings.Split(env, "=")
		if keyval[0] == key {
			return keyval[1], true, nil
		}
	}
	return "", false, nil
}

//Hostname returns the host name reported by the kernel
func (goshnix *Goshnix) Hostname() (name string, err error) {
	command := "hostname"
	op, err := goshnix.client.execute_command(command)
	if err != nil {
		return "", err
	}
	return op, nil
}

// Link creates newname as a hard link to the oldname file
func (goshnix *Goshnix) Link(oldname, newname string) error {
	command := fmt.Sprintf("link %s %s", oldname, newname)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// Mkdir creates a new directory with the specified name and permission bits
func (goshnix *Goshnix) Mkdir(name string, perm os.FileMode) error {
	mode := int(perm)
	command := fmt.Sprintf("mkdir %s --mode=%d", name, mode)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// Readlink returns the destination of the named symbolic link
func (goshnix *Goshnix) Readlink(name string) (string, error) {
	command := fmt.Sprintf("readlink %s", name)
	op, err := goshnix.client.execute_command(command)
	if err != nil {
		return "", err
	}
	return op, nil
}

// Remove removes the named file or directory
func (goshnix *Goshnix) Remove(name string) error {
	command := fmt.Sprintf("rm -r %s", name)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters
func (goshnix *Goshnix) RemoveAll(path string) error {
	command := fmt.Sprintf("rm -rf %s", path)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// Rename renames (moves) oldpath to newpath. If newpath already exists, Rename replaces it. OS-specific restrictions may apply when oldpath and newpath are in different directories
func (goshnix *Goshnix) Rename(oldpath, newpath string) error {
	command := fmt.Sprintf("mv %s %s", oldpath, newpath)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// Symlink creates newname as a symbolic link to oldname
func (goshnix *Goshnix) Symlink(oldname, newname string) error {
	command := fmt.Sprintf("ln -s %s %s", oldname, newname)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// Setenv sets the value of the environment variable named by the key (TODO: As the ssh session is new each time, it could make the Setenv useless)
func (goshnix *Goshnix) Setenv(key, value string) error {
	command := fmt.Sprintf("export \"%s=%s\"", key, value)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

/* XXX: File utils */

// Stat returns a FileInfo describing the named file
func (goshnix *Goshnix) Stat(name string) (os.FileInfo, error) {
	command := fmt.Sprintf("stat %s", name)
	statop, staterr := goshnix.client.execute_command(command)
	if staterr != nil {
		return nil, staterr
	}
	lines := strings.Split(statop, "\n")
	var file string
	var isdir bool
	var mode os.FileMode
	var modtime time.Time
	var size int64
	for _, line := range lines {
		trimline := strings.Trim(line, " ")
		trimlines := strings.Fields(trimline)
		switch trimlines[0] {
		case "File:":
			file = strings.Trim(trimlines[1], "'")
		case "Size":
			sizei, _ := strconv.Atoi(trimlines[1])
			size = int64(sizei)
			if trimlines[7] == "directory" {
				isdir = true
			} else {
				isdir = false
			}
		case "Access":
			for _, key := range trimlines {
				if key == "Uid" {
					mdata := strings.Split(strings.Trim(trimlines[1], "()"), "/")[0]
					modei, _ := strconv.Atoi(mdata)
					mode = os.FileMode(modei)
				}
			}
		case "Modify":
			/*
				date := trimlines[1]
				tym := trimlines[2]
			*/
			// TODO: Claculate time and set modtime
		}
	}
	finfo := &fileinfo{}
	finfo.name = file
	finfo.isdir = isdir
	finfo.mode = mode
	finfo.modtime = modtime
	finfo.size = size
	return finfo, nil
}

// Get Base name of the file
func (finfo *fileinfo) Name() string {
	return finfo.name
}

// Get length in bytes for regular files
func (finfo *fileinfo) Size() int64 {
	return finfo.size
}

// Get file mode bits
func (finfo *fileinfo) Mode() os.FileMode {
	return finfo.mode
}

// Get the modification time
func (finfo *fileinfo) ModTime() time.Time {
	return finfo.modtime
}

// Check if directory
func (finfo *fileinfo) IsDir() bool {
	return finfo.isdir
}

// Get underying data source
func (finfo *fileinfo) Sys() interface{} {
	return finfo
}

// Kill a running process by its pid
func (goshnix *Goshnix) Kill(pid int) error {
	command := fmt.Sprintf("kill -9 %d", pid)
	_, err := goshnix.client.execute_command(command)
	if err != nil {
		return err
	}
	return nil
}

// ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename
func (goshnix *Goshnix) ReadDir(dirname string) ([]os.FileInfo, error) {
	command := fmt.Sprintf("ls -a %s", dirname) // It should sort by default
	op, err := goshnix.client.execute_command(command)
	if err != nil {
		return nil, err
	}
	files := strings.Split(op, "\n")
	var fileinfos []os.FileInfo
	for _, file := range files {
		fileinfo, err := goshnix.Stat(file)
		if err != nil {
			return fileinfos, err
		}
		fileinfos = append(fileinfos, fileinfo)
	}
	return fileinfos, nil
}

// ReadFile reads the file named by filename and returns the contents. A successful call returns err == nil, not err == EOF. Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func (goshnix *Goshnix) ReadFile(filename string) ([]byte, error) {
	op, err := goshnix.client.get_file_content(filename)
	if err != nil {
		return nil, err
	}
	return []byte(op), nil
}
