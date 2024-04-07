package main

import (
  "context"
  "errors"
  "io"
  "regexp"
  "strconv"
  "github.com/aws/aws-lambda-go/lambda"
  "net/http"
)

type MyEvent struct {
  Name string `json:"name"`
}


type Response struct {
    Count string `json:"count"`
}

func getOccupancy() ( int, error ){
  resp, err := http.Get("https://portal.rockgympro.com/portal/public/8490bc5e774d0034d09420df23a224b9/occupancy?=&iframeid=occupancyCounter&fId=")
  if err != nil {
    return -1, err
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return -1, err
  }

  // Dumb but works
  ptmRegex, err := regexp.Compile(`'PTM' : \{\s*'capacity' : \d+,\s*'count' : (\d+),`)
  if err != nil {
    return -1, err
  }

  matches := ptmRegex.FindStringSubmatch(string(body))
  if len(matches) > 1 {
    count := matches[1]

    i, err := strconv.Atoi(count)
    if err != nil  {
      return -1, err
    }

    return i, nil
  }

  return -1, errors.New("no count found in response")
}

func HandleRequest(ctx context.Context) (Response, error) {
  count, err := getOccupancy()
  if err != nil {
    return Response{}, err
  }

  return Response{Count: strconv.Itoa(count) }, nil
}

func main() {
  lambda.Start(HandleRequest)
}

