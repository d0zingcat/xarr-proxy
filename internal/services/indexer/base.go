package indexer

import (
	"regexp"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

type baseIndexer struct{}

func (b *baseIndexer) GetTitle(key string) string {
	return utils.RemoveEpisode(key)
}

func (b *baseIndexer) GetSearchTitle(title string) []string {
	log.Info().Msg("2222")
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
	RSS RSS `json:"rss"`
}

type RSS struct {
	Channel      Channel `json:"channel"`
	XmlnsAtom    string  `json:"_xmlns:atom"`
	XmlnsTorznab string  `json:"_xmlns:torznab"`
	Version      string  `json:"_version"`
}

type Channel struct {
	Link  Link   `json:"link"`
	Title string `json:"title"`
	Items []Item `json:"item"`
}

type Item struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	GUID            string          `json:"guid"`
	Prowlarrindexer Prowlarrindexer `json:"prowlarrindexer"`
	Comments        string          `json:"comments"`
	PubDate         string          `json:"pubDate"`
	Size            string          `json:"size"`
	Link            string          `json:"link"`
	Category        []string        `json:"category"`
	Enclosure       Enclosure       `json:"enclosure"`
	Attr            []Attr          `json:"attr"`
}

type Attr struct {
	Name   string `json:"_name"`
	Value  string `json:"_value"`
	Prefix string `json:"__prefix"`
}

type Enclosure struct {
	URL    string `json:"_url"`
	Length string `json:"_length"`
	Type   string `json:"_type"`
}

type Prowlarrindexer struct {
	ID   string `json:"_id"`
	Text string `json:"__text"`
}

type Link struct {
	Rel    string `json:"_rel"`
	Type   string `json:"_type"`
	Prefix string `json:"__prefix"`
}

const (
	Category             string = "category"
	Downloadvolumefactor string = "downloadvolumefactor"
	Grabs                string = "grabs"
	Infohash             string = "infohash"
	Peers                string = "peers"
	Seeders              string = "seeders"
	Tag                  string = "tag"
	Uploadvolumefactor   string = "uploadvolumefactor"
)
