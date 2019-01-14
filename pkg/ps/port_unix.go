// +build linux solaris

package ps

import (
	"fmt"
	"os/exec"
	"strconv"
)

func associatedPorts(pid int) (ports []uint16, err error) {
	cmd := fmt.Sprintf("ss -l -p -n | grep \"pid=%d,\"", pid)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return ports, fmt.Errorf("Failed to execute command: %s", err)
	}
	outputString := string(out)
	var port string
eachCharacter:
	for i := 0; i < len(outputString); i++ {
		r := outputString[i : i+1]
		if r == " " && len(port) > 1 {
			p, err := strconv.Atoi(port)
			if err != nil {
				return ports, err
			}
			ports = append(ports, uint16(p))
			port = ""
			continue
		}
		// to accommodate ipv6, skip all entries between [ and ]
		if r == "[" {
			for i < len(outputString) {
				i++
				if string(outputString[i:i+1]) == "]" {
					continue eachCharacter
				}
			}
		}
		_, err := strconv.Atoi(r)
		if err != nil {
			continue
		}
		if string(r) == " " {
			continue
		}
		if i > 0 && string(out[i-1]) == ":" {
			port = string(r)
			continue
		}
		if port != "" {
			port += string(r)
		}
	}

	return
}
