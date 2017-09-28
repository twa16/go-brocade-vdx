package brocadevdx

import (
	"strings"
	"strconv"
)

func (vdx *VDXSwitch) GetARPTable() ([]ARPTableEntry, error) {
	data, err := vdx.Cmd("show arp")
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processARPTableOutput(dataLines)

	//Return
	return entries, nil
}

func (vdx *VDXSwitch) GetARPTableFiltered(filter string) ([]ARPTableEntry, error) {
	data, err := vdx.Cmd("show arp | in "+filter)
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processARPTableOutput(dataLines)

	//Return
	return entries, nil
}


func processARPTableOutput(dataLines []string) ([]ARPTableEntry) {
	//Storage for our entries
	var entries []ARPTableEntry

	//Process Lines
	for _, line := range dataLines {
		fields := strings.Fields(line)                     // Split line into array delimited by whitespace
		//Ignore bad lines
		if len(fields) != 7 {
			continue
		}
		if _, err := strconv.Atoi(strings.Split(fields[0],".")[0]); err == nil { //Check if the line starts with a number, if it does, its a VLANPortConfig
			entry := ARPTableEntry{
				Address: fields[0],
				MACAddress: fields[1],
				Interface: fields[2]+" "+fields[3],
				MacResolved: fields[4],
				Age: fields[5],
				Type: fields[6],
			}

			//Add to list
			entries = append(entries, entry)
		}
	}
	return entries
}
