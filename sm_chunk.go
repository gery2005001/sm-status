package main

import (
	"fmt"
	"log"
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
	if SmNetworkInfo.CurrenEpoch <= 0 || SmNetworkInfo.CurrentLayer <= 0 {
		return TsChunks
	}

	//添加Epoch开始Layer块
	whenLayer := int64(SmNetworkInfo.Epoch.LayerStart) - int64(SmNetworkInfo.CurrentLayer)
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
	whenLayer = int64(nextEpochLayerStart) - int64(SmNetworkInfo.CurrentLayer)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  fmt.Sprintf("Epoch %d", SmNetworkInfo.Epoch.Number+1),
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.LayerEnd + 1,
		Type:  1,
	})

	//开始Gap 12H的Layer
	whenLayer = int64(SmNetworkInfo.Epoch.Gap12H.LayerStart) - int64(SmNetworkInfo.CurrentLayer)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "12H Begin",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap12H.LayerStart,
		Type:  2,
	})
	//结束Gap 12H的Layer
	whenLayer = int64(SmNetworkInfo.Epoch.Gap12H.LayerEnd) - int64(SmNetworkInfo.CurrentLayer)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "12H End",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap12H.LayerEnd,
		Type:  2,
	})
	//开始Gap 24L的Layer
	whenLayer = int64(SmNetworkInfo.Epoch.Gap24L.LayerStart) - int64(SmNetworkInfo.CurrentLayer)
	whenDuration = time.Duration(whenLayer * LayerDuration)
	TsChunks = append(TsChunks, Chunk{
		Name:  "24L Begin",
		Desc:  utility.DurationToTimeFormat(whenDuration * time.Second),
		When:  whenDuration,
		Layer: SmNetworkInfo.Epoch.Gap24L.LayerStart,
		Type:  2,
	})
	//结束Gap 12H的Layer
	whenLayer = int64(SmNetworkInfo.Epoch.Gap24L.LayerEnd) - int64(SmNetworkInfo.CurrentLayer)
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
							whenLayer := int64(elg.Layer) - int64(SmNetworkInfo.CurrentLayer)
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
	log.Printf("Current Epoch is: %d,Current Layer is: %d", SmNetworkInfo.CurrenEpoch, SmNetworkInfo.CurrentLayer)
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
			if chunk.Layer == SmNetworkInfo.CurrentLayer {
				over = true
				allChunks[n].Type = 0
			}
		}
		if !over {
			allChunks = append(allChunks, Chunk{
				Name:  "【Now】",
				Desc:  time.Now().Format("2006-01-02") + "<br />" + time.Now().Format("15:04:05"),
				When:  0,
				Layer: SmNetworkInfo.CurrentLayer,
				Type:  0,
			})
		}

		sort.Slice(allChunks, func(i, j int) bool {
			return allChunks[i].Layer < allChunks[j].Layer
		})

		htmlData := "<div class=\"chunks-container\">"
		for i := 0; i < len(allChunks); i++ {
			class := GetBlockColorClass(allChunks[i].Type)
			htmlData += fmt.Sprintf("<div class=\"chunks-block %s\">", class)
			htmlData += fmt.Sprintf("<b>%d</b><br />%s<br />%s", allChunks[i].Layer, allChunks[i].Name, allChunks[i].Desc)
			htmlData += "</div>"
		}
		htmlData += "</div>"

		return htmlData
	}
	return ""
}
