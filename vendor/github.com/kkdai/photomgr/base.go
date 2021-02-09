package photomgr

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type baseCrawler struct {

	//Init on inherit class
	baseAddress  string
	entryAddress string

	//To store current baseCrawler post result
	storedPostURLList   []string
	storedPostTitleList []string
	storedStarList      []int
}

var (
	threadId = regexp.MustCompile(`M.(\d*).`)
	imageId  = regexp.MustCompile(`([^\/]+)\.(png|jpg)`)
)

func (b *baseCrawler) HasValidURL(url string) bool {
	return threadId.Match([]byte(url))
}

// Return parse page result count, it will be 0 if you still not parse any page
func (b *baseCrawler) GetCurrentPageResultCount() int {
	return len(b.storedPostTitleList)
}

// Get post title by index in current parsed page
func (b *baseCrawler) GetPostTitleByIndex(postIndex int) string {
	if postIndex >= len(b.storedPostTitleList) {
		return ""
	}
	return b.storedPostTitleList[postIndex]
}

// Get post URL by index in current parsed page
func (b *baseCrawler) GetPostUrlByIndex(postIndex int) string {
	if postIndex >= len(b.storedPostURLList) {
		return ""
	}

	return b.storedPostURLList[postIndex]
}

// Get post like count by index in current parsed page
func (b *baseCrawler) GetPostStarByIndex(postIndex int) int {
	if postIndex >= len(b.storedStarList) {
		return 0
	}
	return b.storedStarList[postIndex]
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (b *baseCrawler) worker(destDir string, linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		resp, err := http.Get(target)
		if err != nil {
			log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
			continue
		}
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			m, err = png.Decode(resp.Body)
			if err != nil {
				log.Printf("image.Decode\nerror: %s\ntarget: %s", err, target)
				continue
			}
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgInfo := imageId.FindStringSubmatch(target)
			finalPath := destDir + "/" + imgInfo[1] + "." + imgInfo[2]
			out, err := os.Create(filepath.FromSlash(finalPath))
			if err != nil {
				log.Printf("os.Create\nerror: %s\n", err)
				continue
			}
			defer out.Close()
			switch imgInfo[2] {
			case "jpg":
				jpeg.Encode(out, m, nil)
			case "png":
				png.Encode(out, m)
			case "gif":
				gif.Encode(out, m, nil)
			}
		}
	}
}
