package main

import (
  "errors"
  "io"
  "regexp"
  "strconv"
  "net/http"
  "fmt"
)

type MyEvent struct {
  Name string `json:"name"`
}


type Response struct {
  Count string `json:"count"`
}

func getOccupancy() ( int, int, error ){
  scbCount, ptmCount := -1, -1
  resp, err := http.Get("https://portal.rockgympro.com/portal/public/8490bc5e774d0034d09420df23a224b9/occupancy?=&iframeid=occupancyCounter&fId=")
  if err != nil {
    return ptmCount, scbCount, err
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return ptmCount, scbCount, err
  }

  // Dumb but works
  ptmRegex, err := regexp.Compile(`'PTM' : \{\s*'capacity' : \d+,\s*'count' : (\d+),`)
  scbRegex, err := regexp.Compile(`'SCB' : \{\s*'capacity' : \d+,\s*'count' : (\d+),`)
  if err != nil {
    return ptmCount, scbCount, err
  }


  matches := ptmRegex.FindStringSubmatch(string(body))
  if len(matches) > 1 {
    count := matches[1]

    i, err := strconv.Atoi(count)
    if err != nil  {
      return ptmCount, scbCount, err
    }

    ptmCount = i;
  }

  matches = scbRegex.FindStringSubmatch(string(body))
  if len(matches) > 1 {
    count := matches[1]

    i, err := strconv.Atoi(count)
    if err != nil  {
      return ptmCount, scbCount, err
    }

    scbCount = i;
  }

  return ptmCount, scbCount, errors.New("no count found in response")
}

func main() {
  ptmCount, scbCount, _:= getOccupancy()
  fmt.Printf("Portsmouth occupancy: %d\n", ptmCount)
  fmt.Printf("Scarborough occupancy: %d\n", scbCount)
}

