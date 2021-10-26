package bestSeller

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

var Sel = ""

func LinkBestSeller(ctx context.Context) {

	chromedp.Run(ctx, GoBestSeller())
}

func GoBestSeller() chromedp.Tasks {
	log.Print(Sel)
	res := []byte{}
	return chromedp.Tasks{
		chromedp.Click(Sel, chromedp.NodeVisible),
		chromedp.Sleep(10),
		chromedp.ActionFunc(func(c context.Context) error {
			log.Print("--------------------------->")
			return nil
		}),
		chromedp.CaptureScreenshot(&res),
		chromedp.ActionFunc(func(c context.Context) error {
			if err := ioutil.WriteFile("screenshot1.png", res, 0o644); err != nil {
				log.Fatal(err)
			}
			return nil
		}),
		//chromedp.NavigateBack(),
	}
}

// func analysisHtml() {
// 	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
// 	dom.Find(`div[class="a-section a-spacing-large"]`).Each(func(i int, sel *goquery.Selection) {
// 		log.Print(sel.Find(`div[class="a-column a-span8"]>h2`).Text())
// 		sel.Find(`li[class="a-carousel-card"]`).Each(func(i int, sel *goquery.Selection) {
// 			href, exists := sel.Find(`a[class="a-link-normal"]`).Attr("href")
// 			if exists {
// 				log.Printf("href:%s", href)
// 				mianKey := strings.SplitAfter(href, "pd_rd_i=")
// 				if mianKey[1] != "" {
// 					log.Printf("mianKey:%s", mianKey)
// 					hrefs = append(hrefs, mianKey[1])
// 				}
// 			}
// 			log.Print(sel.Text())
// 		})
// 	})
// }

func GoBestGoods(key, href string) chromedp.Tasks {
	sel := fmt.Sprintf(`//div[@class="zg-carousel-general-faceout"]/div/div[@class="a-row"]/a[contains(@href,"pd_rd_i=%s")]`, href)
	log.Printf("sel:%s", sel)
	res := []byte{}
	return chromedp.Tasks{
		chromedp.Click(sel),
		chromedp.Sleep(10),
		chromedp.ActionFunc(func(c context.Context) error {
			log.Print("--------------------------->")
			return nil
		}),
		//chromedp.WaitVisible(`div[id="cm-cr-dp-review-list"]`, chromedp.ByID),
		//chromedp.OuterHTML(`div[id="cm-cr-dp-review-list"]`, &res, chromedp.ByID),
		chromedp.CaptureScreenshot(&res),
		chromedp.ActionFunc(func(c context.Context) error {
			//SerachInfo(&res)
			//log.Printf("res : %s", res)
			if err := ioutil.WriteFile("screenshot1.png", res, 0o644); err != nil {
				log.Fatal(err)
			}
			return nil
		}),
		chromedp.NavigateBack(),
	}
}
