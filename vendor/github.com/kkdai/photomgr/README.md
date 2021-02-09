photomgr: A photo crawler manager for gomobile usage in Golang
======================
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/photomgr/master/LICENSE) [![GoDoc](https://godoc.org/github.com/kkdai/photomgr?status.svg)](https://godoc.org/github.com/kkdai/photomgr)[![Build Status](https://travis-ci.org/kkdai/photomgr.svg)](https://travis-ci.org/kkdai/photomgr)


A photomgr is a Go package to browse some info (such as PTT/CK101.. more)  and download image if any default image in that article. This tool help you to download those photos for your backup, all the photos still own by original creator. 


Install
--------------

    go get github.com/kkdai/photomgr

Usage
---------------------

refer example for more detail.

```go
ptt := NewPTT()

//Set path for download image
ptt.BaseDir = "YOURPATH"

pageIndex := 0

totalPostCount := ptt.ParsePttPageByIndex(pageIndex)

//iterator this sample
for i := 0; i < totalPostCount ; i_++ {
		title := p.GetPostTitleByIndex(i)
		likeCount := p.GetPostStarByIndex(i)
		fmt.Printf("%d:[%d★]%s\n", i, likeCount, title)
		
		//download image 
		url := p.GetPostUrlByIndex(i)
		p.Crawler(url, 25)
}


```

For more detail check example_iloveptt which is a sample program with [iloveptt](http://github.com/kkdai/ilovept)

Gomobile supported
--------------

To let your package support [gomobile](https://godoc.org/golang.org/x/mobile/cmd/gomobile), need note as follow:

- Only support `string`, `int`  (no slice and array)
- Need a constructor for your structure such as `NewYOUROBJ() *YOUROBJ`.
 
 
How to build this package in iOS using Gomobile
---------------

Here is howto to teach you how to use this in iOS.

```
//Get gomobile and init it
//It might take some time
go get golang.org/x/mobile/cmd/gomobile
gomobile init

//Get this package and build package for iOS
go get github.com/kkdai/photomgr
gomobile bind -target=ios github.com/kkdai/photomgr

//It will generate a photomgr.framework in your cd $GOPATH/src/github.com/kkdai/photomgr
//Drag photomgr.framework in your xcode project
```
     
for more detail, check my iOS project [PhotoViewer](https://github.com/kkdai/PhotoViewer)     


TODO
---------------

- [x] PTT
  - [x] gomobile refine
  - [x] download image
- [x] CK101
  - [x] broad travaeral
  - [x] download image
- [ ] Timliao
  - [ ] broad travaeral
  - [ ] download image
- [ ] Gigacircle
  - [ ] broad travaeral
  - [ ] download image



Contribute
---------------

Please open up an issue on GitHub before you put a lot efforts on pull request.
The code submitting to PR must be filtered with `gofmt`


Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/kkdai/photomgr/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

