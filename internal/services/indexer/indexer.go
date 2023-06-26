package indexer

import (
	"bytes"
	"regexp"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

type IndexService interface {
	GetTitle(searchKey string) string
	GetSearchTitle(title string) []string
	GenerateOffsetKey(requestWrapper *services.RequestWrapper) string
	GetOffsetList(key string, size int) []int
	CalculateCurrentIndex(offset int, offsetList []int) int
	UpdateIndexerRequest(index int, searchTitleList []string, offsetList []int, indexerRequest *req.IndexerReq)
	UpdateOffsetList(key string, offsetList []int) []int
	ExecuteNewRequest(requestWrapper *services.RequestWrapper) string
	ExecuteFormatRule(xml string) string
	GetIndexerUrl(path string) string
}

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

type indexer struct {
	baseIndexer
	service IndexService
}

func NewIndexerService(service IndexService) *indexer {
	return &indexer{
		service: service,
	}
}

func (i *indexer) ExecuteNewRequest(requestWrapper *services.RequestWrapper) string {
	path := requestWrapper.GetPath()
	uStr := i.service.GetIndexerUrl(path) + "?" + requestWrapper.GetQueryString()
	log.Debug().Msg(uStr)
	response, err := utils.GetClient(nil).Get(uStr)
	if err != nil {
		log.Err(err).Msg("Failed to send request to indexer")
		return ""
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	xml := buf.String()

	return xml
}

//	func (*indexer) GetIndexerSearchReq(r *http.Request) IndexerSearchReq {
//		var req req.IndexerSearchReq
//		q := r.URL.Query()
//		req.Limit = utils.ParseInt(q.Get("limit"))
//		req.Offset = utils.ParseInt(q.Get("offset"))
//		req.SearchKey = q.Get("searchKey")
//		req.SearchType = q.Get("searchType")
//		req.SeasonNumber = q.Get("seasonNumber")
//		return req
//	}
// func (*indexer) CalculateCurrentIndex(offset int, offsetList []int) int {
// 	for index, value := range offsetList {
// 		if value == 0 || value > offset {
// 			return index
// 		}
// 	}
// 	return 0
// }
//
// func (*indexer) GenerateOffsetKey(queryString string) string {
// 	queryValues, _ := url.ParseQuery(queryString)
// 	queryValues.Del("offset")
// 	queryValues.Del("apikey")
// 	return queryValues.Encode()
// }
