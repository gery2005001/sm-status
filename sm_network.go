package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sm-status/utility"
	"sort"
	"time"
)

const (
	LayerDuration     = 300   //每个Layer完成的时间
	T12HDuration      = 43200 //12小时Duration
	T24HDuration      = 86400 //24小时Duration
	EpochLayers       = 4032  //每个Epoch的Layer数量
	Gap12HLayers      = 2736  //Epoch开始后的第几个Layer开始Gap12H
	T12HLayers        = 144   //12小时完成的Layer
	T24HLayers        = 288   //24小时完成的Layer
	GetNetworkInfoUrl = "https://mainnet-explorer-api.spacemesh.network/network-info"
)

var BlockClass = []string{"block-now", "block-start-end", "block-running", "block-elgs"}

var SmNetworkInfo NetworkInfo = NetworkInfo{}

type Statistics struct {
	Capacity      int64 `json:"capacity" bson:"capacity"`         // Average tx/s rate over capacity considering all layers in the current epoch.
	Decentral     int64 `json:"decentral" bson:"decentral"`       // Distribution of storage between all active smeshers.
	Smeshers      int64 `json:"smeshers" bson:"smeshers"`         // Number of active smeshers in the current epoch.
	Transactions  int64 `json:"transactions" bson:"transactions"` // Total number of transactions processed by the state transition function.
	Accounts      int64 `json:"accounts" bson:"accounts"`         // Total number of on-mesh accounts with a non-zero coin balance as of the current epoch.
	Circulation   int64 `json:"circulation" bson:"circulation"`   // Total number of Smesh coins in circulation. This is the total balances of all on-mesh accounts.
	Rewards       int64 `json:"rewards" bson:"rewards"`           // Total amount of Smesh minted as mining rewards as of the last known reward distribution event.
	RewardsNumber int64 `json:"rewardsnumber" bson:"rewardsnumber"`
	Security      int64 `json:"security" bson:"security"`   // Total amount of storage committed to the network based on the ATXs in the previous epoch.
	TxsAmount     int64 `json:"txsamount" bson:"txsamount"` // Total amount of coin transferred between accounts in the epoch. Incl coin transactions and smart wallet transactions.
}

type Stats struct {
	Current    Statistics `json:"current"`
	Cumulative Statistics `json:"cumulative"`
}

type Gap12H struct {
	Start      uint32 `json:"start" bson:"start"`
	End        uint32 `json:"end" bson:"end"`
	LayerStart uint32 `json:"layerstart" bson:"layerstart"`
	LayerEnd   uint32 `json:"layerend" bson:"layerend"`
	Layers     uint32 `json:"layers" bson:"layers"`
}

type Gap24L struct {
	Start      uint32 `json:"start" bson:"start"`
	End        uint32 `json:"end" bson:"end"`
	LayerStart uint32 `json:"layerstart" bson:"layerstart"`
	LayerEnd   uint32 `json:"layerend" bson:"layerend"`
	Layers     uint32 `json:"layers" bson:"layers"`
}

type Epoch struct {
	Number     int32  `json:"number" bson:"number"`
	Start      uint32 `json:"start" bson:"start"`
	End        uint32 `json:"end" bson:"end"`
	LayerStart uint32 `json:"layerstart" bson:"layerstart"`
	LayerEnd   uint32 `json:"layerend" bson:"layerend"`
	Layers     uint32 `json:"layers" bson:"layers"`
	Stats      Stats  `json:"stats"`
	Gap12H     Gap12H `json:"gap12h"`
	Gap24L     Gap24L `json:"gap24l"`
}

type Layer struct {
	Number       uint32 `json:"number" bson:"number"`
	Status       int    `json:"status" bson:"status"`
	Txs          uint32 `json:"txs" bson:"txs"`
	Start        uint32 `json:"start" bson:"start"`
	End          uint32 `json:"end" bson:"end"`
	TxsAmount    uint64 `json:"txsamount" bson:"txsamount"`
	Rewards      uint64 `json:"rewards" bson:"rewards"`
	Epoch        uint32 `json:"epoch" bson:"epoch"`
	Hash         string `json:"hash" bson:"hash"`
	BlocksNumber uint32 `json:"blocksnumber" bson:"blocksnumber"`
}

type NetworkInfo struct {
	Epoch Epoch `json:"epoch"`
	Layer Layer `json:"layer"`
}

// 从网络获取Epoch信息
func GetNetworkInfo() error {
	log.Println("start get network info...")
	resp, err := http.Get(GetNetworkInfoUrl)
	if err != nil {
		log.Println("get epoch infomation failed: ", err)
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if err := json.NewDecoder(resp.Body).Decode(&SmNetworkInfo); err != nil {
		log.Println("decode epoch Json information failed: ", err)
		return err
	}
	log.Println("successfully get current epoch")

	//Gap 12H开始和结束的Layer
	SmNetworkInfo.Epoch.Gap12H.LayerStart = SmNetworkInfo.Epoch.LayerStart + Gap12HLayers
	SmNetworkInfo.Epoch.Gap12H.LayerEnd = SmNetworkInfo.Epoch.Gap12H.LayerStart + T12HLayers
	SmNetworkInfo.Epoch.Gap12H.Layers = T12HLayers
	//Gap 12H开始和结束的时间
	SmNetworkInfo.Epoch.Gap12H.Start = LayerDuration*Gap12HLayers + SmNetworkInfo.Epoch.Start
	SmNetworkInfo.Epoch.Gap12H.End = SmNetworkInfo.Epoch.Gap12H.Start + T12HDuration

	//Gap 24L开始和结束的Layer
	SmNetworkInfo.Epoch.Gap24L.LayerStart = SmNetworkInfo.Epoch.Gap12H.LayerEnd + T24HLayers
	SmNetworkInfo.Epoch.Gap24L.LayerEnd = SmNetworkInfo.Epoch.Gap24L.LayerStart + T24HLayers
	SmNetworkInfo.Epoch.Gap24L.Layers = T24HLayers
	//Gap 24L开始和结束的时间
	SmNetworkInfo.Epoch.Gap24L.Start = SmNetworkInfo.Epoch.Gap12H.End + T24HDuration
	SmNetworkInfo.Epoch.Gap24L.End = SmNetworkInfo.Epoch.Gap24L.Start + T12HDuration

	return nil
}

func (x *NetworkInfo) GetHtmlString() string {
	if x.Epoch.Number == 0 {
		return ""
	}
	type Layers struct {
		Name  string
		Desc  string
		Time  time.Duration
		Layer uint32
		Type  int // 0: now 1: start-end 2: 12H-24L
	}

	layerTimes := []Layers{}

	layerTimes = append(layerTimes, Layers{
		Name:  fmt.Sprintf("Epoch %d", x.Epoch.Number),
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.LayerStart,
		Type:  1,
	})
	layerTimes = append(layerTimes, Layers{
		Name:  fmt.Sprintf("Epoch %d", x.Epoch.Number+1),
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.LayerEnd,
		Type:  1,
	})
	layerTimes = append(layerTimes, Layers{
		Name:  "12H begin",
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.Gap12H.LayerStart,
		Type:  2,
	})
	layerTimes = append(layerTimes, Layers{
		Name:  "12H end",
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.Gap12H.LayerEnd,
		Type:  2,
	})
	layerTimes = append(layerTimes, Layers{
		Name:  "24L begin",
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.Gap24L.LayerStart,
		Type:  2,
	})
	layerTimes = append(layerTimes, Layers{
		Name:  "24L end",
		Desc:  "",
		Time:  0,
		Layer: x.Epoch.Gap24L.LayerEnd,
		Type:  2,
	})

	over := false
	for n, layer := range layerTimes {
		if layer.Layer == x.Layer.Number {
			over = true
			layerTimes[n].Desc = time.Now().Format("2006-01-02 15:04:05")
			layerTimes[n].Type = 0
		} else {
			incLayer := int64(layer.Layer) - int64(x.Layer.Number)
			layerTimes[n].Time = time.Duration(incLayer*LayerDuration) * time.Second
			layerTimes[n].Desc = utility.DurationToTimeFormat(layerTimes[n].Time)
		}
	}
	if !over {
		layerTimes = append(layerTimes, Layers{
			Name:  "Now",
			Desc:  time.Now().Format("2006-01-02 15:04:05"),
			Time:  0,
			Layer: x.Layer.Number,
			Type:  0,
		})
	}

	sort.Slice(layerTimes, func(i, j int) bool {
		return layerTimes[i].Layer < layerTimes[j].Layer
	})

	var HtmlData string

	HtmlData += "<table class=\"block-table\"><tr>"
	for _, lt := range layerTimes {
		class := GetBlockColorClass(lt.Type)
		HtmlData += fmt.Sprintf("<td class=\"%s\">", class)
		HtmlData += fmt.Sprintf("<b>%d</b><br />%s<br />%s", lt.Layer, lt.Name, lt.Desc)
		HtmlData += "</td>"
	}
	HtmlData += "</tr></table>"

	return HtmlData
}

func GetBlockColorClass(n int) string {
	if n < len(BlockClass) {
		return BlockClass[n]
	} else {
		return ""
	}
}
