package main

import "testing"


func TestConnection(t *testing.T) {
	hostName := "127.0.0.1"
	port := 3306 //default for mysql
	if !makeTcpConnection(hostName, port) {
		t.Error()
	}
}

func TestBanner(t *testing.T) {
	//the test will succeed on a system with mysql running on 3306 port
	hostName := "127.0.0.1"
	port := 3306 //default for mysql
	resWithBanner := grabBanner(hostName, port)
	if !isMysqlBanneer(resWithBanner) {
		t.Error()
	}
}

func TestNmap(t *testing.T) {
	//the test will succeed on a system with mysql running on 3306 port and with nmap installed
	hostName := "127.0.0.1"
	port := 3306 //default for mysql
	nmapPath, err := findNmap()
	if err != nil {
		t.Error("nmap not installed")
	}
	resWithNmapShort, err := executeNmapForShortOutput(nmapPath, hostName, port)
	if err != nil {
		t.Error("nmap error")
	}
	nmapOutputFields := extractInfoFromShortNmapOutput(resWithNmapShort)
	if len(nmapOutputFields) == 0 {
		t.Error("nmap output error")
	}
	if nmapOutputFields[2] != "mysql"{
		t.Error()
	}
}

