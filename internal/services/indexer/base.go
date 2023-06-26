package indexer

import (
	"encoding/xml"
	"regexp"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/utils"
)

type baseIndexer struct{}

func (b *baseIndexer) GetTitle(key string) string {
	return utils.RemoveEpisode(key)
}

func (b *baseIndexer) GetSearchTitle(title string) []string {
	return []string{}
}

func (b *baseIndexer) GenerateOffsetKey(requestWrapper *services.RequestWrapper) string {
	builder := strings.Builder{}
	builder.WriteString(requestWrapper.GetPath())
	queryString := requestWrapper.GetQueryString()
	m := regexp.MustCompile(`(offset=\d+|apikey=\w+)`)
	queryString = m.ReplaceAllString(queryString, "")
	builder.WriteString("?" + queryString)
	return builder.String()
}

func (b *baseIndexer) CalculateCurrentIndex(offset int, offsetList []int) int {
	for index, offsetValue := range offsetList {
		if offsetValue == 0 || offsetValue > offset {
			return index
		}
	}
	return 0
}

func (b *baseIndexer) UpdateIndexerRequest(index int, searchTitleList []string, offsetList []int, indexerRequest *req.IndexerReq) {
	if index > 0 {
		lastIndex := index - 1
		title := searchTitleList[index]
		searchKey := indexerRequest.SearchKey
		searchKey = strings.ReplaceAll(searchKey, searchTitleList[0], title)
		if lastIndex > 0 {
			searchKey = strings.ReplaceAll(searchKey, searchTitleList[lastIndex], title)
		}
		indexerRequest.SearchKey = searchKey
		indexerRequest.Offset = indexerRequest.Offset - offsetList[lastIndex]
	}
}

func (b *baseIndexer) UpdateOffsetList(key string, offsetList []int) []int {
	return offsetList
}

func (b *baseIndexer) GetOffsetList(key string, size int) []int {
	offsetList := make([]int, size)
	return offsetList
}

func (b *baseIndexer) ExecuteFormatRule(xml string) string {
	return xml
}

func (b *baseIndexer) GetIndexerUrl(path string) string {
	return path
}

func (b *baseIndexer) ExecuteNewRequest(requestWrapper *services.RequestWrapper) string {
	return ""
}

type TorznabRss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Torznab string   `xml:"xmlns:torznab,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title    string   `xml:"title"`
	AtomLink AtomLink `xml:"atom:link"`
	Items    []Item   `xml:"item"`
}

type AtomLink struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Guid        string    `xml:"guid"`
	Prowlarr    string    `xml:"prowlarrindexer"`
	Comments    string    `xml:"comments"`
	PubDate     string    `xml:"pubDate"`
	Size        int       `xml:"size"`
	Link        string    `xml:"link"`
	Categories  []int     `xml:"category"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

//
// type TorznabRss struct {
// 	RSS RSS `xml:"rss"`
// }
//
// type RSS struct {
// 	Channel      Channel `xml:"channel"`
// 	XmlnsAtom    string  `xml:"_xmlns:atom"`
// 	XmlnsTorznab string  `xml:"_xmlns:torznab"`
// 	Version      string  `xml:"_version"`
// }
//
// type Channel struct {
// 	Link  Link   `xml:"link"`
// 	Title string `xml:"title"`
// 	Items []Item `xml:"item"`
// }
//
// type Item struct {
// 	Title           string          `xml:"title"`
// 	Description     string          `xml:"description"`
// 	GUID            string          `xml:"guid"`
// 	Prowlarrindexer Prowlarrindexer `xml:"prowlarrindexer"`
// 	Comments        string          `xml:"comments"`
// 	PubDate         string          `xml:"pubDate"`
// 	Size            string          `xml:"size"`
// 	Link            string          `xml:"link"`
// 	Category        []string        `xml:"category"`
// 	Enclosure       Enclosure       `xml:"enclosure"`
// 	Attr            []Attr          `xml:"attr"`
// }
//
// type Attr struct {
// 	Name   string `xml:"_name"`
// 	Value  string `xml:"_value"`
// 	Prefix string `xml:"__prefix"`
// }
//
// type Enclosure struct {
// 	URL    string `xml:"_url"`
// 	Length string `xml:"_length"`
// 	Type   string `xml:"_type"`
// }
//
// type Prowlarrindexer struct {
// 	ID   string `xml:"_id"`
// 	Text string `xml:"__text"`
// }
//
// type Link struct {
// 	Rel    string `xml:"_rel"`
// 	Type   string `xml:"_type"`
// 	Prefix string `xml:"__prefix"`
// }
