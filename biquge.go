package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

var biqugeURL = "http://www.biquge.com.cn"

// BiqugeIndex 使用说明
func BiqugeIndex(c *gin.Context) {
	m := &routeManual{
		Path: ":bookID",
		Info: map[string]string{":bookID": "小说在笔趣阁的代码"},
	}
	c.JSON(http.StatusOK, m)
}

// BiqugeRSS 返回相应RSS
func BiqugeRSS(c *gin.Context) {
	bookID := c.Param("bookID")
	rss, err := BiqugeChapter(bookID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, rss)
	}
}

// BiqugeChapter 获取章节信息，返回rss
func BiqugeChapter(bookID string) (string, error) {
	href := fmt.Sprintf("%s/book/%s/", biqugeURL, bookID)
	res, err := http.Get(href)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	feed := &feeds.Feed{
		Title: doc.Find("#info > h1").Text(),
		Link:  &feeds.Link{Href: href},
	}

	selection := doc.Find("#list > dl > dd")

	lenth := selection.Length()
	if lenth > 3 {
		lenth = 3
	}

	for i, s := 0, selection.Last(); i < lenth; i++ {
		href, ok := s.Find("a").Attr("href")
		if ok {
			title := s.Find("a").Text()
			item := &feeds.Item{
				Title:       title,
				Link:        &feeds.Link{Href: biqugeURL + href},
				Description: BiqugeChapterCtx(href),
			}
			feed.Items = append(feed.Items, item)
		}
		s = s.Prev()
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
}

// BiqugeChapterCtx 获取并返回章节内容
func BiqugeChapterCtx(href string) string {
	res, err := http.Get(fmt.Sprintf("%s/%s", biqugeURL, href))
	if err != nil {
		return "无法获取该章节！"
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "无法获取该章节！"
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "无法获取该章节！"
	}

	ctx, err := doc.Find("#content").Html()
	if err != nil {
		return "无法获取该章节！"
	}

	return strings.Replace(ctx, "<br><br>", "<br>", -1)
}
