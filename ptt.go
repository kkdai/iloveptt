package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	logging "github.com/op/go-logging"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

const (
	BasePttAddress = "https://www.ptt.cc"
	EntryAddress   = "https://www.ptt.cc/bbs/Beauty/index.html"
)

var (
	baseDir  string
	threadId = regexp.MustCompile(`M.(\d*).`)
	imageId  = regexp.MustCompile(`([^\/]+)\.(png|jpg)`)
	log      = logging.MustGetLogger("iloveck101")
)

func worker(destDir string, linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		resp, err := http.Get(target)
		if err != nil {
			log.Debug("Http.Get\nerror: %s\ntarget: %s", err, target)
			continue
		}
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Debug("image.Decode\nerror: %s\ntarget: %s", err, target)
			continue
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgInfo := imageId.FindStringSubmatch(target)
			out, err := os.Create(destDir + "/" + imgInfo[1] + "." + imgInfo[2])
			if err != nil {
				log.Debug("os.Create\nerror: %s", err)
				continue
			}
			defer out.Close()
			switch imgInfo[2] {
			case "jpg":
				jpeg.Encode(out, m, nil)
			case "png":
				png.Encode(out, m)
			}
		}
	}
}

func crawler(target string, workerNum int) {
	doc, err := goquery.NewDocument(target)
	if err != nil {
		log.Fatal(err)
	}

	articleTitle := doc.Find(".article-metaline .article-meta-value").Text()
	//Title and folder
	dir := fmt.Sprintf("%v/%v - %v", baseDir, threadId.FindStringSubmatch(target)[1], articleTitle)
	os.MkdirAll(dir, 0755)

	//Concurrecny
	linkChan := make(chan string)
	wg := new(sync.WaitGroup)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go worker(dir, linkChan, wg)
	}

	//Parse Image, currently support <IMG SRC> only
	doc.Find(".richcontent").Each(func(i int, s *goquery.Selection) {
		imgLink, exist := s.Find("img").Attr("src")
		if exist {
			linkChan <- "http:" + imgLink
		} else {
			fmt.Println("Don't have any image in this article.")
		}
	})

	close(linkChan)
	wg.Wait()
}

func parsePttBoardIndex(page int) (hrefs []string) {
	doc, err := goquery.NewDocument(EntryAddress)
	if err != nil {
		log.Fatal(err)
	}
	hrefs = make([]string, 0)
	maxPageNumberString := ""
	var PageWebSide string
	if page > 0 {
		// Find page result
		doc.Find(".btn-group a").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "上頁") {
				href, exist := s.Attr("href")
				if exist {
					targetString := strings.Split(href, "index")[1]
					targetString = strings.Split(targetString, ".html")[0]
					fmt.Println("total page:", targetString)
					maxPageNumberString = targetString
				}
			}
		})
		pageNum, _ := strconv.Atoi(maxPageNumberString)
		pageNum = pageNum - page
		PageWebSide = fmt.Sprintf("https://www.ptt.cc/bbs/Beauty/index%d.html", pageNum)
	} else {
		PageWebSide = EntryAddress
	}

	doc, err = goquery.NewDocument(PageWebSide)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".r-ent").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".title").Text())
		likeCount, _ := strconv.Atoi(s.Find(".nrec span").Text())
		href, _ := s.Find(".title a").Attr("href")
		link := BasePttAddress + href
		hrefs = append(hrefs, link)
		fmt.Printf("%d:[%d★]%s\n", i, likeCount, title)
	})

	// Print pages
	fmt.Printf("Pages: ")
	for i := page - 3; i <= page+2; i++ {
		if i >= 0 {
			if i == page {
				fmt.Printf("[%v] ", i)
			} else {
				fmt.Printf("%v ", i)
			}
		}
	}
	fmt.Printf("(n:next, p:prev, quit: quit program.)\n")

	return hrefs
}

func main() {
	var format = logging.MustStringFormatter("%{level} %{message}")
	logging.SetFormatter(format)
	logging.SetLevel(logging.INFO, "iloveptt")

	usr, _ := user.Current()
	baseDir = fmt.Sprintf("%v/Pictures/iloveptt", usr.HomeDir)

	var workerNum int
	rootCmd := &cobra.Command{
		Use:   "iloveptt",
		Short: "Download all the images in given post url",
		Run: func(cmd *cobra.Command, args []string) {
			page := 0
			hrefs := parsePttBoardIndex(page)

			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Print("ptt:> ")

				if !scanner.Scan() {
					break
				}

				line := scanner.Text()
				parts := strings.Split(line, " ")
				cmd := parts[0]
				args := parts[1:]

				switch cmd {
				case "quit":
					quit = true
				case "n":
					page = page + 1
					hrefs = parsePttBoardIndex(page)
				case "p":
					if page > 0 {
						page = page - 1
					}
					hrefs = parsePttBoardIndex(page)
				case "s":
					page = 0
					hrefs = parsePttBoardIndex(page)
				case "o":
					open.Run(baseDir)
				case "d":
					index, err := strconv.ParseUint(args[0], 0, 0)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if int(index) >= len(hrefs) {
						fmt.Println("Invalid index")
						continue
					}

					if threadId.Match([]byte(hrefs[index])) {
						crawler(hrefs[index], 25)
						fmt.Println("Done!")
					} else {
						fmt.Println("Unsupport url:", hrefs[index])
					}
				default:
					fmt.Println("Unrecognized command:", cmd, args)
				}
			}
		},
	}

	rootCmd.Flags().IntVarP(&workerNum, "worker", "w", 25, "Number of workers")
	rootCmd.Execute()
}
