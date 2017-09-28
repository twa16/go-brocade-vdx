package brocadevdx

import (
	"strings"
	"strconv"
)

func (vdx *VDXSwitch) GetMacTable() ([]MACTableEntry, error) {
	data, err := vdx.Cmd("show mac-address-table")
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processMACTableOutput(dataLines)

	//Return
	return entries, nil
}

func (vdx *VDXSwitch) GetMacTableFiltered(filter string) ([]MACTableEntry, error) {
	data, err := vdx.Cmd("show mac-address-table | in "+filter)
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processMACTableOutput(dataLines)

	//Return
	return entries, nil
}

func processMACTableOutput(dataLines []string) ([]MACTableEntry) {
	//Storage for our entries
	var entries []MACTableEntry

	//Process Lines
	for _, line := range dataLines {
		fields := strings.Fields(line)                     // Split line into array delimited by whitespace
		//Ignore bad lines
		if len(fields) != 6 {
			continue
		}
		if _, err := strconv.Atoi(fields[0]); err == nil { //Check if the line starts with a number, if it does, its a VLANPortConfig
			entry := MACTableEntry{}
			entry.VlanID, _ = strconv.Atoi(fields[0])
			entry.MACAddress = fields[1]
			entry.Type = fields[2]
			entry.State = fields[3]
			entry.Port = fields[4]+" "+fields[5]

			//Add to list
			entries = append(entries, entry)
		}
	}

	return entries
}
