package main

import "context"

// 安装chrome浏览器
// 导入第三方库
// go get -u github.com/chromedp/chromedp
import (
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
)

func main() {
	// 创建 context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// 生成pdf
	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(`https://www.baidu.com/`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("baidu.pdf", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
// 生成任务列表
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr), // 浏览指定的页面
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx) // 通过cdp执行PrintToPDF
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
