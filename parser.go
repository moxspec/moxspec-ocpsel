package main

import "strings"

func parseBody(bod []byte) [][]string {
	if len(bod) == 0 {
		return [][]string{}
	}

	// chunk example
	// -------------
	// SEL Record ID          : 0001
	//  Record Type           : 02
	//  Timestamp             : 04/15/2019 12:01:29
	//  Generator ID          : 0020
	//  EvM Revision          : 04
	//  Sensor Type           : Event Logging Disabled
	//  Sensor Number         : d1
	//  Event Type            : Sensor-specific Discrete
	//  Event Direction       : Assertion Event
	//  Event Data            : 02ffff
	//  Description           : Log area reset/cleared
	var chunk []string
	var chunkList [][]string
	for _, line := range strings.Split(string(bod), "\n") {
		if strings.HasPrefix(line, "SEL Record ID") {
			if len(chunk) > 0 {
				chunkList = append(chunkList, chunk)
			}
			chunk = []string{}
		}
		if line == "" {
			continue
		}
		chunk = append(chunk, line)
	}
	chunkList = append(chunkList, chunk)

	return chunkList
}

func parseChunkList(chunkList [][]string) []string {
	var result []string
	for _, chunk := range chunkList {
		result = append(result, parseChunk(chunk)...)
	}
	return result
}

func parseChunk(chunk []string) []string {
	gens, sens, eds := extractChunkInfo(chunk)

	res, err := decode(sens, gens, eds)
	if err != nil {
		return []string{}
	}

	return []string{
		strings.Join(chunk, "\n"),
		res.printIndent(" "),
		"",
	}
}

func extractChunkInfo(chunk []string) (string, string, string) {
	var gens, sens, eds string
	for _, line := range chunk {
		e := strings.Split(line, ":")
		if len(e) < 2 {
			continue
		}

		k := strings.TrimSpace(e[0])
		v := strings.TrimSpace(strings.Join(e[1:], ":"))

		switch k {
		case "Generator ID":
			gens = v
		case "Sensor Number":
			sens = v
		case "Event Data":
			eds = v
		}
	}
	return gens, sens, eds
}
