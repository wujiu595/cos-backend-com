package util

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecCmd(name string, args ...string) (output string, err error) {
	cmdWArgs := strings.Split(name, " ")
	var cmd *exec.Cmd
	if len(cmdWArgs) > 1 {
		cmd = exec.Command(cmdWArgs[0], append(cmdWArgs[1:], args...)...)
	} else {
		cmd = exec.Command(name, args...)
	}

	data, err := cmd.CombinedOutput()
	if err != nil {
		cmdstr := strings.Join(append([]string{name}, args...), " ")
		err = fmt.Errorf("cmd: %s, failed: %v", cmdstr, err)
		return string(data), err
	}
	return string(data), nil
}
