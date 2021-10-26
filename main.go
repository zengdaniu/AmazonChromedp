package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/examples/amazon/bestSeller"
)

var cks string
var req *http.Request
var url = "https://www.amazon.com/ref=nav_logo"
var html string

func main() {
	flag.StringVar(&cks, "cookies", "", "cookies")
	flag.Parse()
	if cks == "" {
		cks = "aws-target-static-id=1594024662116-604718; s_fid=7036EC6DEF860921-2A21A16BCFA8F920; regStatus=pre-register; s_dslv=1594024668671; session-id=133-7661309-7893419; session-id-time=2082787201l; i18n-prefs=USD; ubid-main=135-4348489-7503348; lc-main=en_US; session-token=mA4N7h9iVdKALGWi9P3jGxegJR5IVwLZCsM651gP7DaCGyziWFz9D2KvLSXK2mFDMo5lixQV5M6AH4rISIlFDNi7uKPQSPnTAyVUPV86w4TJDpL58VRr4nz91dKJl3QxaoucmEz0OHfazE38yLTcWwGj1PFYLFPNI/UTZ/XopSoLFsP6ry7xIBoVL9SRjQpu; csm-hit=tb:FR3SR8E9QFSRJABD3RHR+b-77C4X5BYMP32T9F6KY72|1635214214080&t:1635214214080&adb:adblk_no"
	}
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36`),
		chromedp.WindowSize(1280, 720),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()

	header := http.Header{}
	header.Add("Cookie", cks)
	req = &http.Request{Header: header}

	// 监听得到第二个tab页的target ID
	ch := make(chan target.ID, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*target.EventTargetCreated); ok &&
			// if OpenerID == "", this is the first tab.
			e.TargetInfo.OpenerID != "" {
			ch <- e.TargetInfo.TargetID
			log.Print("EventTargetCreated")
		} else if e, ok := ev.(*target.EventTargetInfoChanged); ok &&
			// if OpenerID == "", this is the first tab.
			e.TargetInfo.OpenerID != "" {
			ch <- e.TargetInfo.TargetID
			log.Print("EventTargetInfoChanged")
		} else if e, ok := ev.(*target.EventAttachedToTarget); ok &&
			// if OpenerID == "", this is the first tab.
			e.TargetInfo.OpenerID != "" {
			ch <- e.TargetInfo.TargetID
			log.Print("EventAttachedToTarget")
		} else if e, ok := ev.(*page.EventWindowOpen); ok {
			log.Printf("EventWindowOpen :%s", e.WindowName)
		}
	})

	ctx, timeoutcancel := context.WithTimeout(ctx, 50*time.Second)
	defer timeoutcancel()
	chromedp.Run(ctx, setCookie())
	bestSeller.LinkBestSeller(ctx)
	//xshps(ctx)
	//analysisHtml()
	// for _, href := range hrefs {
	// 	c, cnl := context.WithTimeout(ctx, 20*time.Second)
	// 	defer cnl()
	// 	chromedp.Run(c, CheckCheck("a", href))
	// 	time.Sleep(10 * time.Second)
	// 	break
	// }
	log.Print("+++++++++++++++")
	// 第二个tab页
	newCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(<-ch))
	defer cancel()
	if err := chromedp.Run(
		newCtx,
		chromedp.ActionFunc(func(c context.Context) error {
			log.Print("new Event")
			return nil
		}),
		chromedp.Sleep(1*time.Second),
		chromedp.OuterHTML("*", &html),
	); err != nil {
		log.Printf("[scrapeNewArticle] chromedp Run fail,err: %s", err.Error())
		return
	}

}

func setCookie() chromedp.Tasks {
	Nodes := []*cdp.Node{}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			cookies := req.Cookies()
			// add cookies to chrome
			for i := 0; i < len(cookies); i++ {
				err := network.SetCookie(cookies[i].Name, cookies[i].Value).
					WithExpires(&expr).
					WithDomain("www.amazon.com").
					WithHTTPOnly(true).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.Sleep(5),
		chromedp.WaitVisible(`#nav-xshop-container`, chromedp.ByID),
		chromedp.Nodes(`//div[@id="nav-xshop-container"]/div/a`, &Nodes),
		chromedp.ActionFunc(func(c context.Context) error {
			//ioutil.WriteFile("cur.html", []byte(html), 064)
			printNodes(Nodes, 0)
			return nil
		}),
	}
}

func xshps(c context.Context) {
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	dom.Find(`div[id="nav-xshop"]>a`).Each(func(i int, sel *goquery.Selection) {
		key := sel.Text()
		if key == "Best Sellers" {
			log.Print("SearchBestSeller")
			href, _ := sel.Attr("href")
			bestSeller.Sel = fmt.Sprintf(`a[@href="%s")]`, href)
			bestSeller.LinkBestSeller(c)
			return
		}

	})
}

func printNodes(nodes []*cdp.Node, indent int) {
	spaces := strings.Repeat(" ", indent)
	for _, node := range nodes {
		fmt.Print(spaces)
		var extra interface{}
		if node.NodeName == "#text" {
			extra = node.NodeValue
		} else {
			extra = node.Attributes
		}
		fmt.Printf("%s: %q\n", node.NodeName, extra)
		if node.NodeName == "#text" {
			log.Print(node.FullXPath())
			if extra == "Best Sellers" {
				bestSeller.Sel = node.Parent.FullXPath()
			}
		}
		if node.ChildNodeCount > 0 {
			printNodes(node.Children, indent+4)
		}
	}
}

func SerachInfo(res *string) {
	log.Print(*res)
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(*res))
	dom.Find(`div[class="a-row"]`).Each(func(i int, sel *goquery.Selection) {
		log.Print(sel.Text())
	})
}

func CheckCheck(key, href string) chromedp.Tasks {
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
