package brocadevdx

import (
	"strings"
	"strconv"
	"regexp"
)

func (vdx *VDXSwitch) GetLLDPNeighbors() ([]LLDPEntry, error) {
	data, err := vdx.Cmd("show lldp ne")
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processLLDPNeighborOutput(dataLines)

	//Return
	return entries, nil
}

func (vdx *VDXSwitch) GetLLDPNeighborsFiltered(filter string) ([]LLDPEntry, error) {
	data, err := vdx.Cmd("show lldp ne | in "+filter)
	if err != nil {
		return nil, err
	}
	//Convert bytes to string
	dataString := string(data[:])
	//Convert to string array
	dataLines := strings.Split(dataString, "\r\n")

	//Process
	entries := processLLDPNeighborOutput(dataLines)

	//Return
	return entries, nil
}

func processLLDPNeighborOutput(dataLines []string) []LLDPEntry {
	//Storage for our entries
	var entries []LLDPEntry

	//Process Lines
	for _, line := range dataLines {
		fields := regSplit(line, "(\\s\\s+)")                   // Split line into array delimited by whitespace
		//Ignore bad lines
		if len(fields) < 8{
			continue
		}
		if _, err := strconv.Atoi(fields[2]); err == nil { //Check if the line starts with a number, if it does, its a VLANPortConfig
			entry := LLDPEntry{}
			entry.LocalInterface = fields[0]
			entry.DeadInterval, _ = strconv.Atoi(fields[1])
			entry.RemainingLife, _ = strconv.Atoi(fields[2])
			entry.RemoteInterface = fields[3]
			entry.ChassisID = fields[4]
			entry.Tx, _ = strconv.Atoi(fields[5])
			entry.Rx, _ = strconv.Atoi(fields[6])
			entry.SystemName = fields[7]

			//Add to list
			entries = append(entries, entry)
		}
	}
	return entries
}

func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes) + 1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}