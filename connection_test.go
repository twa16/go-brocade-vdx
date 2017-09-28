package brocadevdx

import (
	"testing"
	"os"
	"github.com/k0kubun/pp"
)

func TestConnection(t *testing.T) {
	vdx, err := ConnectToSwitchWithPassword( os.Getenv("TEST_ADDRESS"),  os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"))
	if err != nil {
		t.Error(err)
	}

	_, err = vdx.Cmd("show system")
	if err != nil {
		t.Error(err)
	}

}


func TestMACTableGet(t *testing.T) {
	vdx, err := ConnectToSwitchWithPassword( os.Getenv("TEST_ADDRESS"),  os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"))
	if err != nil {
		t.Error(err)
	}

	macTable, err := vdx.GetMacTable()
	if err != nil {
		t.Error(err)
	}
	pp.Println(macTable[:1])
}

func TestARPTableGet(t *testing.T) {
	vdx, err := ConnectToSwitchWithPassword( os.Getenv("TEST_ADDRESS"),  os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"))
	if err != nil {
		t.Error(err)
	}

	arpTable, err := vdx.GetARPTable()
	if err != nil {
		t.Error(err)
	}
	pp.Println(arpTable[:1])
}

func TestLLDPNeighborGet(t *testing.T) {
	vdx, err := ConnectToSwitchWithPassword( os.Getenv("TEST_ADDRESS"),  os.Getenv("TEST_USERNAME"), os.Getenv("TEST_PASSWORD"))
	if err != nil {
		t.Error(err)
	}

	lldpNe, err := vdx.GetLLDPNeighbors()
	if err != nil {
		t.Error(err)
	}
	pp.Println(lldpNe)
}
