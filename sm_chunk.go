package main

import (
	"fmt"
	"sm-status/utility"
	"sort"
	"time"
)

// 根据NetworkInfo和SmConfig来生成Chunks表
type Chunk struct {
	Name  string        `json:"name"`
	Desc  string        `json:"desc"`
	When  time.Duration `json:"when"` // 0: now  >0: wait <0: end 到达和结束块的时长
	Layer uint32        `json:"layer"`
	Type  int           // 0: now 1: start-end 2: 12H-24L 3: Elgs
}

const MAXCOLUMN = 8

// var ElgChunks = []Chunk{}
// var TsChunks = []Chunk{}

func GetTimeStreamChunks() []Chunk {
	TsChunks := []Chunk{}
	if SmNetworkInfo.Epoch.Number <= 0 {
		return TsChunks
	}

	//添加Epoch开始Layer块
	whenLayer := int(SmNetworkInfo.Epoch.LayerStart) - int(SmNetworkInfo.Layer.Number)
	whenDuration := time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  fmt.Sprintf("Epoch %d", SmNetworkInfo.Epoch.Number),
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.LayerStart,
		Type:  1,
	})
	//下个epoch开始的layer
	nextEpochLayerStart := SmNetworkInfo.Epoch.LayerEnd + 1
	whenLayer = int(nextEpochLayerStart) - int(SmNetworkInfo.Layer.Number)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  fmt.Sprintf("Epoch %d", SmNetworkInfo.Epoch.Number+1),
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.LayerEnd + 1,
		Type:  1,
	})

	//开始Gap 12H的Layer
	whenLayer = int(SmNetworkInfo.Epoch.Gap12H.LayerStart) - int(SmNetworkInfo.Layer.Number)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "12H Begin",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap12H.LayerStart,
		Type:  2,
	})
	//结束Gap 12H的Layer
	whenLayer = int(SmNetworkInfo.Epoch.Gap12H.LayerEnd) - int(SmNetworkInfo.Layer.Number)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "12H End",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap12H.LayerEnd,
		Type:  2,
	})
	//开始Gap 24L的Layer
	whenLayer = int(SmNetworkInfo.Epoch.Gap24L.LayerStart) - int(SmNetworkInfo.Layer.Number)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "24L Begin",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap24L.LayerStart,
		Type:  2,
	})
	//结束Gap 12H的Layer
	whenLayer = int(SmNetworkInfo.Epoch.Gap24L.LayerEnd) - int(SmNetworkInfo.Layer.Number)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "24L End",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap24L.LayerEnd,
		Type:  2,
	})

	return TsChunks
}

func GetElgChunks() []Chunk {
	config := GetConfig()
	ElgChunks := []Chunk{}
	for i, node := range config.Node {
		if len(node.PostInfo) > 0 {
			for j, post := range config.Node[i].PostInfo {
				if len(post.Eligs) > 0 {
					for k, elg := range config.Node[i].PostInfo[j].Eligs {
						if elg.Epoch == uint32(SmNetworkInfo.Epoch.Number) {
							whenLayer := int(elg.Layer) - int(SmNetworkInfo.Layer.Number)
							whenDuration := time.Duration(whenLayer * LayerDuration)
							name := config.Node[i].Name + "<br />" + config.Node[i].PostInfo[j].Title
							nameTag := "【✓】"
							if whenDuration <= 0 {
								if whenDuration == 0 {
									nameTag = "【Now】"
								} else {
									if elg.Total > 0 {
										nameTag = fmt.Sprintf("<br />%.4f", float64(elg.Total)/1000000000)
									}
								}
							} else {
								nameTag = fmt.Sprintf("【%d】", config.Node[i].PostInfo[j].Eligs[k].Count)
							}
							ElgChunks = append(ElgChunks, Chunk{
								Name:  name + nameTag,
								Layer: elg.Layer,
								When:  whenDuration,
								Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
								Type:  3,
							})
						}
					}
				}
			}
		}
	}

	return ElgChunks
}

func GetChunksTableHTML() string {
	tc := GetTimeStreamChunks()
	ec := GetElgChunks()

	allChunks := []Chunk{}

	if len(tc) > 0 {
		allChunks = append(allChunks, tc...)
	}
	if len(ec) > 0 {
		allChunks = append(allChunks, ec...)
	}

	if len(allChunks) > 0 {
		over := false
		for n, chunk := range allChunks {
			if chunk.Layer == SmNetworkInfo.Layer.Number {
				over = true
				allChunks[n].Desc = time.Now().Format("2006-01-02") + "<br />" + time.Now().Format("15:04:05")
				allChunks[n].Type = 0
				allChunks[n].Name += "【Now】"
			}
		}
		if !over {
			allChunks = append(allChunks, Chunk{
				Name:  "【Now】",
				Desc:  time.Now().Format("2006-01-02") + "<br />" + time.Now().Format("15:04:05"),
				When:  0,
				Layer: SmNetworkInfo.Layer.Number,
				Type:  0,
			})
		}

		sort.Slice(allChunks, func(i, j int) bool {
			return allChunks[i].Layer < allChunks[j].Layer
		})

		htmlData := "<table class=\"block-table\">"
		for i := 0; i < len(allChunks); i += 8 {
			htmlData += "<tr>"
			for j := i; j < i+8 && j < len(allChunks); j++ {
				class := GetBlockColorClass(allChunks[j].Type)
				htmlData += fmt.Sprintf("<td class=\"td-chunk %s\">", class)
				htmlData += fmt.Sprintf("<b>%d</b><br />%s<br />%s", allChunks[j].Layer, allChunks[j].Name, allChunks[j].Desc)
				htmlData += "</td>"
			}
			htmlData += "</tr>"
		}
		htmlData += "</table>"
		return htmlData
	}
	return ""
}
