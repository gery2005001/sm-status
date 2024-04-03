package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ResActivation struct {
	Data       []Activation       `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

type PaginationMetadata struct {
	TotalCount  int64 `json:"totalCount"`
	PageCount   int64 `json:"pageCount"`
	PerPage     int64 `json:"perPage"`
	Next        int64 `json:"next"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
	Current     int64 `json:"current"`
	Previous    int64 `json:"previous"`
}

type Activation struct {
	Id                string `json:"id" bson:"id"`             //nolint will fix it later.
	SmesherId         string `json:"smesher" bson:"smesher"`   //nolint will fix it later // id of smesher who created the ATX
	Coinbase          string `json:"coinbase" bson:"coinbase"` // coinbase account id
	PrevAtx           string `json:"prevAtx" bson:"prevAtx"`   // previous ATX pointed to
	NumUnits          uint32 `json:"numunits" bson:"numunits"` // number of PoST data commitment units
	CommitmentSize    uint64 `json:"commitmentSize" bson:"commitmentSize"`
	PublishEpoch      uint32 `json:"publishEpoch" bson:"publishEpoch"`
	TargetEpoch       uint32 `json:"targetEpoch" bson:"targetEpoch"`
	TickCount         uint64 `json:"tickCount" bson:"tickCount"`
	Weight            uint64 `json:"weight" bson:"weight"`
	EffectiveNumUnits uint32 `json:"effectiveNumUnits" bson:"effectiveNumUnits"`
	Received          int64  `json:"received" bson:"received"`
}

const apiAeerss = "https://mainnet-explorer-api.spacemesh.network"

func GetActivations(id string) (ResActivation, error) {
	atxs := ResActivation{}
	api := fmt.Sprintf("%s/smeshers/%s/atxs", apiAeerss, id)

	resp, err := http.Get(api)
	if err != nil {
		log.Println(err)
		return atxs, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return atxs, err
	}

	err = json.Unmarshal(body, &atxs)
	if err != nil {
		log.Println(err)
		return atxs, err
	}

	return atxs, nil
}
