package system

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)


func ExecShellScript(shell_script string) (int, error){
	command := exec.Command("/bin/sh", shell_script)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	status := command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if err != nil {
		return status,err
	}
	return status,nil
}



func Exec(cmd string) (string, error) {
	command := exec.Command("sh", "-c", cmd)
	bytes, err := command.Output()
	return string(bytes), err
}

func ExecOutput(cmd string) string {
	output, err := Exec(cmd)
	if err != nil {
		return ""
	}
	Trim(&output)
	return output
}


func Trim(str *string) {
	*str = strings.TrimSpace(*str)
	*str = strings.Replace(*str, "\n", "", -1)
}


