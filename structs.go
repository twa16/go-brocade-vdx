package brocadevdx

import (
	"golang.org/x/crypto/ssh"
	"io"
)

type VDXSwitch struct {
	Hostname string
	Username string
	Password string
	stdin    io.WriteCloser
	stdout   io.Reader
	session  *ssh.Session
	ReadChan chan *string
	StopChan chan struct{}
	client   *ssh.Client
	Timeout  int
}

type MACTableEntry struct {
	VlanID     int
	MACAddress string
	Type       string
	State      string
	Port       string
}

type ARPTableEntry struct {
	Address     string
	MACAddress  string
	Interface   string
	MacResolved string
	Age         string
	Type        string
}

type LLDPEntry struct {
	LocalInterface  string
	DeadInterval    int
	RemainingLife   int
	RemoteInterface string
	ChassisID       string
	Tx              int
	Rx              int
	SystemName      string
}
