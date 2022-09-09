package ocpsel

import "strings"

// ParseText parses given ipmitool output
func ParseText(bod []byte) [][]string {
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

func ParseChunkList(chunkList [][]string) []string {
	var result []string
	for _, chunk := range chunkList {
		result = append(result, ParseChunk(chunk)...)
	}
	return result
}

func ParseChunk(chunk []string) []string {
	gens, sens, eds := ExtractChunkInfo(chunk)

	res, err := Decode(sens, gens, eds)
	if err != nil {
		return []string{}
	}

	return []string{
		strings.Join(chunk, "\n"),
		res.SummaryIndent(" "),
		"",
	}
}

func ExtractChunkInfo(chunk []string) (string, string, string) {
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
