package photomgr

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type FBAlbum struct {
	//Inherit
	baseCrawler

	//To store current FBAlbum post result
	BaseDir string
}

func NewFBAlbum() *FBAlbum {
	c := new(FBAlbum)
	c.baseAddress = "https://www.FBAlbum.com"
	c.entryAddress = "http://FBAlbum.com/forum-3465-1.html"
	return c
}
func (b *FBAlbum) HasValidURL(url string) bool {
	log.Println("url=", url)
	return true
}

func (p *FBAlbum) Crawler(target string, workerNum int) {

	doc, err := goquery.NewDocument(target)
	if err != nil {
		panic(err)
	}

	title := doc.Find("h1#thread_subject").Text()

	log.Println("[FBAlbum]:", title, " starting downloading...")
	dir := fmt.Sprintf("%v/%v - %v", p.BaseDir, "FBAlbum", title)
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

//Set FBAlbum board page index, fetch all post and return article count back
func (p *FBAlbum) ParseFBAlbumPageByIndex(page int) int {
	doc, err := goquery.NewDocument(p.entryAddress)
	if err != nil {
		log.Fatal(err)
	}

	urlList := make([]string, 0)
	postList := make([]string, 0)
	starList := make([]int, 0)

	var PageWebSide string
	page = page + 1 //one base
	if page > 1 {
		// Find page result
		PageWebSide = fmt.Sprintf("http://FBAlbum.com/forum-3465-%d.html", page)
	} else {
		PageWebSide = p.entryAddress
	}
	//fmt.Println("Page", PageWebSide)

	doc, err = goquery.NewDocument(PageWebSide)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".titleBox").Each(func(i int, s *goquery.Selection) {

		star := ""
		title := ""
		url := ""
		starInt := 0
		s.Find(".blockTitle a").Each(func(i int, tQ *goquery.Selection) {
			title, _ = tQ.Attr("title")
			url, _ = tQ.Attr("href")
		})
		s.Find(".icoPage img").Each(func(i int, starC *goquery.Selection) {
			star_c, _ := starC.Attr("title")
			if strings.Contains(star_c, "熱度") {
				star = strings.TrimPrefix(star_c, "熱度:")
				star = strings.TrimSpace(star)
				starInt, _ = strconv.Atoi(star)
			}
			//}
		})
		urlList = append(urlList, url)
		starList = append(starList, starInt)
		postList = append(postList, title)
	})

	p.storedPostURLList = urlList
	p.storedStarList = starList
	p.storedPostTitleList = postList

	return len(p.storedPostTitleList)
}
