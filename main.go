package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func parseLocationCount(locationCode string, body []byte) (int, error) {
	// yuck
	pattern := fmt.Sprintf(`'%s' : \{\s*'capacity' : \d+,\s*'count' : (\d+),`, locationCode)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return -1, err
	}

	matches := regex.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		count, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, err
		}
		return count, nil
	}

	return -1, fmt.Errorf("no count found for %s", locationCode)
}

func getOccupancy() (int, int, error) {
	resp, err := http.Get("https://portal.rockgympro.com/portal/public/8490bc5e774d0034d09420df23a224b9/occupancy?=&iframeid=occupancyCounter&fId=")
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, err
	}

	ptmCount, ptmErr := parseLocationCount("PTM", body)
	scbCount, scbErr := parseLocationCount("SCB", body)

	if ptmErr != nil && scbErr != nil {
		return -1, -1, errors.New("no count found in response")
	}

	return ptmCount, scbCount, nil
}

func main() {
	ptmCount, scbCount, _ := getOccupancy()
	fmt.Printf("Portsmouth occupancy: %d\n", ptmCount)
	fmt.Printf("Scarborough occupancy: %d\n", scbCount)
}
