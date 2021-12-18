package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/ShangRui-hash/hackflow"
	"github.com/sirupsen/logrus"
)

func main() {
	hackflow.SetDebug(true)
	// urlCollectorStdout, err := hackflow.GetUrlCollector().Run(&hackflow.UrlCollectorCofnig{
	// 	Keyword:      ".php?id=1",
	// 	SearchEngine: "baidu",
	// 	Proxy:        "socks://127.0.0.1:7890",
	// 	OuputFormat:  "protocol_domain",
	// })
	// if err != nil {
	// 	logrus.Error("URLCollector.Run failed,err:", err)
	// 	return
	// }
	// // urlCh := make(chan string, 1024)
	// // urlCh <- "https://mahara.org"
	// // urlCh <- "http://www.baidu.com"
	// // urlCh <- "https://www.360doc.com"
	// // close(urlCh)

	dirsearchStdout, err := hackflow.GetDirSearch().Run(strings.NewReader("https://mahara.org\nhttp://www.baidu.com\nhttps://www.360doc.com\n"), hackflow.DirSearchConfig{
		FullURL: true,
	})
	if err != nil {
		logrus.Error("dirSearch.Run failed,err:", err)
		return
	}
	scanner := bufio.NewScanner(dirsearchStdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
