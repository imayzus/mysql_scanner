package main

import (
	//"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// fmt.Println("len(args)=" + strconv.Itoa(len(os.Args)))
	hostName := "127.0.0.1"
	port := 3306 //default for mysql

	if len(os.Args) == 3 {
		hostName = os.Args[1]
		portStr := os.Args[2]
		portNum, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println(("non-numeric port: " + portStr))
			return
		} else {
			port = portNum
		}
	} else {
		progName := "go run scan_for_mysql.go"
		fmt.Println("Usage: " + progName + " <hostName> <port>, running with defaults hostName=" + hostName +
			", port=" + strconv.Itoa(port))
	}
	fmt.Println("checking if the port is open...")
	if makeTcpConnection(hostName, port) {
		fmt.Println("port is open, proceeding with other tests")
	} else {
		fmt.Println(" a tcp connection could not be established")
		return
	}

	if !checkWithBanner(hostName, port) {
		fmt.Println("based on the banner, mysql is not running on " + hostName + ":" + strconv.Itoa(port))
	}

	fmt.Println("checking with nmap...")
	nmapPath, err := findNmap()
	if err != nil {
		fmt.Println("nmap not found on the system")
		return
	}
	resWithNmapShort, err := executeNmapForShortOutput(nmapPath, hostName, port)
	if err != nil {
		fmt.Println("error running nmap: ", err)
		return
	}
	//fmt.Println(resWithNmapShort)
	nmapOutputFields := extractInfoFromShortNmapOutput(resWithNmapShort)
	if len(nmapOutputFields) == 0 {
		fmt.Println("unable to parse nmap output: " + resWithNmapShort)
		return
	}
	if nmapOutputFields[2] == "mysql" {
		fmt.Println("based on nmap, mysql is running on port " + nmapOutputFields[0] +
			", the state is " + nmapOutputFields[1] + ", mysql version: " +
			nmapOutputFields[3] + " " + nmapOutputFields[4])
	} else {
		fmt.Println("based on nmap, mysql is not running on " + hostName + ":" + strconv.Itoa(port))
		fmt.Println("nmap output: port: " + nmapOutputFields[0] +
			", state: " + nmapOutputFields[1] + ", service: " + nmapOutputFields[2] +
		 	", version: " + nmapOutputFields[3] + " " + nmapOutputFields[4])
	}
}


func checkWithBanner(hostName string, port int) bool {
	fmt.Println("retrieving the service banner...")
	resWithBanner := grabBanner(hostName, port)
	fmt.Println("retrieved banner: ", resWithBanner)
	if isMysqlBanneer(resWithBanner) {
		fmt.Println("based on the banner, mysql is running on specified host and port")
		return true
	} else {
		return false
	}
}

func makeTcpConnection(hostName string, port int) bool {
	_, err := net.DialTimeout("tcp", hostName + ":" + strconv.Itoa(port),
		time.Second*10)
	if err == nil {
		return true
	} else {
		return false
	}
}



func findNmap() (string, error) {
	nmapPath, err := exec.LookPath("nmap")
	if err != nil {
		fmt.Println("Error: ", err)
		return "Nmap not found", err
	} else {
		return nmapPath, nil
	}
}


func executeNmapForShortOutput(nmapPath string, hostName string, port int) (string, error) {
	//for Xml output, use Args: []string{ nmapPath, "-p", strconv.Itoa(port),
	//			"-sV", "-T4", "-oX", "-", hostName,
	//		}
	cmd := &exec.Cmd{
		Path: nmapPath,
		Args: []string{ nmapPath, "-sV", "-p", strconv.Itoa(port),
			hostName,
		},
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	} else {
		return string(output), nil
	}
}


func extractInfoFromShortNmapOutput(resWithNmapShort string) [] string {
	lines := strings.Split(resWithNmapShort, "\n")
	re := regexp.MustCompile("PORT\\s+STATE\\s+SERVICE\\s+VERSION")
	ind := 0
	for ; ind < len(lines); ind++ {
		line := lines[ind]
		matched := re.MatchString(line)
		if matched {
			break
		}
	}
	//if match was not found, ind would iterate all the way to len(lines) and ind+1 would be outside len(lines)
	if (ind+1) < len(lines) {
		infoLine := lines[ind+1]
		nmapOutputFields := regexp.MustCompile("\\s+").Split(infoLine, -1)
		return nmapOutputFields
	} else {
		return []string{}
	}
}


func grabBanner(ip string, port int) string {
	connection, err := net.DialTimeout(
		"tcp",
		ip + ":"+strconv.Itoa(port),
		time.Second*10)
	if err != nil {
		return "no connection"
	}
	buffer := make([]byte, 4096)
	connection.SetReadDeadline(time.Now().Add(time.Second*5))
	// Set timeout
	numBytesRead, err := connection.Read(buffer)
	if err != nil {
		return err.Error()
	}
	res := buffer[0:numBytesRead]
	//log.Printf("Banner from port %d\n%s\n", port,
	//	buffer[0:numBytesRead])
	return string(res)

}

func isMysqlBanneer(banner string) bool {
	return strings.Contains(banner, "mysql")
}