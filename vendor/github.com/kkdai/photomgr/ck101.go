package photomgr

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type CK101 struct {
	//Inherit
	baseCrawler

	//To store current CK101 post result
	BaseDir string
}

func NewCK101() *CK101 {
	c := new(CK101)
	c.baseAddress = "https://ck101.com"
	c.entryAddress = "https://ck101.com/forum-1345-1.html"
	return c
}
func (b *CK101) HasValidURL(url string) bool {
	log.Println("url=", url)
	return true
}

func (p *CK101) GetUrlPhotos(target string) []string {
	var resultSlice []string

	doc, err := goquery.NewDocument(target)
	if err != nil {
		panic(err)
	}

	doc.Find("div[itemprop=articleBody] img").Each(func(i int, img *goquery.Selection) {
		imgUrl, _ := img.Attr("file")
		resultSlice = append(resultSlice, imgUrl)
	})
	return resultSlice
}

func (p *CK101) Crawler(target string, workerNum int) {

	doc, err := goquery.NewDocument(target)
	log.Println("Down load target URL=", target)
	if err != nil {
		panic(err)
	}

	title := doc.Find("h1").Text()
	log.Println("[CK101]:", title, " starting downloading...")
	dir := fmt.Sprintf("%v/%v - %v", p.BaseDir, "CK101", title)
	if exist, _ := exists(dir); exist {
		//fmt.Println("Already download")
		return
	}
	os.MkdirAll(dir, 0755)

	linkChan := make(chan string)
	wg := new(sync.WaitGroup)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go p.worker(dir, linkChan, wg)
	}

	doc.Find("div[itemprop=articleBody] img").Each(func(i int, img *goquery.Selection) {
		imgUrl, _ := img.Attr("file")
		linkChan <- imgUrl
	})

	close(linkChan)
	wg.Wait()
}

//Set CK101 board page index, fetch all post and return article count back
func (p *CK101) ParseCK101PageByIndex(page int) int {
	doc, err := goquery.NewDocument(p.entryAddress)
	if err != nil {
		log.Fatal(err)
	}

	urlList := make([]string, 0)
	postList := make([]string, 0)

	var PageWebSide string
	page = page + 1 //one base
	if page > 1 {
		// Find page result
		PageWebSide = fmt.Sprintf("http://ck101.com/forum-1345-%d.html", page)
	} else {
		PageWebSide = p.entryAddress
	}
	//fmt.Println("Page", PageWebSide)

	doc, err = goquery.NewDocument(PageWebSide)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".cl_box").Each(func(i int, s *goquery.Selection) {
		title := ""
		url := ""
		s.Find("a").Each(func(i int, tQ *goquery.Selection) {
			title, _ = tQ.Attr("title")
			goUrl, _ := tQ.Attr("href")
			url = fmt.Sprintf("%s/%s", p.baseAddress, goUrl)
		})

		urlList = append(urlList, url)
		postList = append(postList, title)
	})

	p.storedPostURLList = urlList
	p.storedPostTitleList = postList

	return len(p.storedPostTitleList)
}
