package indexer

import (
	"bytes"

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

type Indexer struct {
	service IndexService
}

func NewIndexerService(service IndexService) *Indexer {
	return &Indexer{
		service: service,
	}
}

func (i *Indexer) ExecuteNewRequest(requestWrapper *services.RequestWrapper) string {
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

func (i Indexer) GetTitle(key string) string {
	return i.service.GetTitle(key)
}

func (i Indexer) GetSearchTitle(title string) []string {
	return i.service.GetSearchTitle(title)
}

func (i Indexer) GenerateOffsetKey(requestWrapper *services.RequestWrapper) string {
	return i.service.GenerateOffsetKey(requestWrapper)
}

func (i Indexer) CalculateCurrentIndex(offset int, offsetList []int) int {
	return i.service.CalculateCurrentIndex(offset, offsetList)
}

func (i Indexer) UpdateIndexerRequest(index int, searchTitleList []string, offsetList []int, indexerRequest *req.IndexerReq) {
	i.service.UpdateIndexerRequest(index, searchTitleList, offsetList, indexerRequest)
}

func (i Indexer) UpdateOffsetList(key string, offsetList []int) []int {
	return i.service.UpdateOffsetList(key, offsetList)
}

func (i Indexer) GetOffsetList(key string, size int) []int {
	return i.service.GetOffsetList(key, size)
}

func (i Indexer) ExecuteFormatRule(xml string) string {
	return i.service.ExecuteFormatRule(xml)
}

func (i Indexer) GetIndexerUrl(path string) string {
	return i.service.GetIndexerUrl(path)
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
