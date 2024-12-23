package main

func getNodeStatusTableHTML() string {
	//输出Node 状态表
	config := GetConfig()

	htmlData := SmNetworkInfo.GetHtmlString()

	for n := 0; n < len(config.Node); n++ {
		htmlData += config.Node[n].GetNodeStatusTableHTMLString()
	}

	return htmlData
}
