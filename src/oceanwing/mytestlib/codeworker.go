package mytestlib

import (
	"fmt"
	"os/exec"
)

func execShell() {
	// "ping", "-c4", "127.0.0.1"
	cmd := exec.Command("ping", "-c4", "127.0.0.1")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error occur: %s\n", err)
		return
	}
	fmt.Printf("out put value is %s\n", out)
}
