# BCSV

Marsgaller / Unmarshaller for the BCSV format

## Installation

```ssh
go get -u github.com/gouniverse/bcsv@v1.10.0
```

## Example

- Using array

```go
import "github.com/gouniverse/bcsv"

rows := [][]string{}

rows = append(rows, []string{"city", "county"})
rows = append(rows, []string{"Sofia", "BG"})
rows = append(rows, []string{"London", "UK"})

bcsvString, _ := bcsv.MarshalToString(rows)

```

- Using struct

```go

import "github.com/gouniverse/bcsv"

type Address struct {
	City    string `bcsv:"city"`
	Country string `bcsv:country"`
}

addressRows := []Address{
  {
    City:    "Sofia",
    Country: "BG",
  },
  {
    City:    "London",
    Country: "UK",
  },
}

bcsvString, _ := bcsv.MarshalToString(addressRows)
```
