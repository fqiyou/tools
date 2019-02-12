package system

import (
	"os/exec"
	"strings"
)

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
