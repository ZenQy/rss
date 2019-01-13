package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

var dotaURL = "http://dota2.uuu9.com/"

// DotaIndex 使用说明
func DotaIndex(c *gin.Context) {
	m := &routeManual{
		Path: ":listID",
		Info: map[string]string{":listID": "UUU9的dota2新闻分类ID"},
	}
	c.JSON(http.StatusOK, m)
}

// DotaRSS 返回相应RSS
func DotaRSS(c *gin.Context) {
	listID := c.Param("listID")
	rss, err := DotaChapter(listID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, rss)
	}
}

// DotaChapter 获取章节信息，返回rss
func DotaChapter(listID string) (string, error) {
	href := fmt.Sprintf("%sList/List_%s.shtml", dotaURL, listID)
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
	title := doc.Find("body > div.all > div.main > div > div.w680 > div.title > h2 > a:last-child").Text()
	feed := &feeds.Feed{
		Title: "UUU9 DOTA2-" + Decode(title),
		Link:  &feeds.Link{Href: href},
	}

	selection := doc.Find("body > div.all > div.main > div > div.w680 > div.con.p10 > ul > li")

	selection.Each(func(i int, s *goquery.Selection) {
		if i > 2 {
			return
		}
		href, ok := s.Find("a").Attr("href")
		if ok {
			title := s.Find("a").Text()
			item := &feeds.Item{
				Title:       Decode(title),
				Link:        &feeds.Link{Href: href},
				Description: DotaChapterCtx(href),
			}
			feed.Items = append(feed.Items, item)
		}
	})

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
}

// DotaChapterCtx 获取并返回章节内容
func DotaChapterCtx(href string) string {
	res, err := http.Get(href)
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

	return Decode(strings.Replace(ctx, "<br><br>", "<br>", -1))
}
