package filter

import (
	"net/http"
	"regexp"
	"strconv"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/services/indexer"
	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

type (
	IndexerFilter struct {
		Service indexer.IndexService
	}
)

func NewIndexerFilter(service indexer.IndexService) *IndexerFilter {
	return &IndexerFilter{
		Service: service,
	}
}

func (i *IndexerFilter) DoFilter(w http.ResponseWriter, r *http.Request) {
	requestWrapper := services.NewRequestWrapper(r)
	indexerRequest := GetIndexerRequest(requestWrapper)

	// 处理查询
	xml := ""
	searchKey := indexerRequest.SearchKey
	if searchKey != "" {
		// 无绝对集数，去除 00
		r := regexp.MustCompile(" 00$")
		searchKey = r.ReplaceAllString(searchKey, "")
		indexerRequest.SearchKey = searchKey
		// 获取所有待查询标题
		title := i.Service.GetTitle(indexerRequest.SearchKey)
		searchTitleList := i.Service.GetSearchTitle(title)
		// 计算 offset
		offset := indexerRequest.Offset
		size := len(searchTitleList)
		offsetKey := i.Service.GenerateOffsetKey(requestWrapper)
		offsetList := i.Service.GetOffsetList(offsetKey, size)
		// 计算当前标题下标
		index := i.Service.CalculateCurrentIndex(offset, offsetList)
		// 更新参数
		i.Service.UpdateIndexerRequest(index, searchTitleList, offsetList, indexerRequest)
		UpdateRequestWrapper(indexerRequest, requestWrapper)
		// 请求
		xml = i.Service.ExecuteNewRequest(requestWrapper)
		count := utils.XmlCount(xml)
		log.Debug().
			Any("title", title).
			Any("searchTitleList", searchTitleList).
			Any("offsetKey", offsetKey).
			Any("offsetList", offsetList).
			Any("index", index).
			Any("xml", xml[:3000]).Msg("first request")
		index++
		for index < size && indexerRequest.Limit-count > 0 {
			log.Debug().Msgf("index %v < size %v, limit: %v, count", index, size, indexerRequest.Limit, count)
			// 更新参数
			offset = offset + count
			if index == size-1 {
				if indexerRequest.SeasonNumber != "" && offsetList[index-1] < indexerRequest.Limit {
					indexerRequest.SeasonNumber = ""
					indexerRequest.EpisodeNumber = ""
				} else {
					break
				}
			}
			indexerRequest.Offset = offset
			indexerRequest.Limit = indexerRequest.Limit - count
			offsetList[index-1] = offset
			i.Service.UpdateOffsetList(offsetKey, offsetList)
			i.Service.UpdateIndexerRequest(index, searchTitleList, offsetList, indexerRequest)
			UpdateRequestWrapper(indexerRequest, requestWrapper)
			// 重新请求
			newXml := i.Service.ExecuteNewRequest(requestWrapper)
			count = utils.XmlCount(newXml)
			if count > 0 {
				xml = utils.XmlMerge(xml, newXml)
			}
			index++
		}
	} else {
		xml = i.Service.ExecuteNewRequest(requestWrapper)
	}
	// 处理结果
	log.Logger.Debug().Msgf("before: %v", xml[:3000])
	xml = i.Service.ExecuteFormatRule(xml)
	log.Logger.Debug().Msgf("after: %v", xml[:3000])
	// 返回
	w.Write([]byte(xml))
}

func GetIndexerRequest(requestWrapper *services.RequestWrapper) *req.IndexerReq {
	indexerRequest := &req.IndexerReq{}
	indexerRequest.SearchKey = requestWrapper.GetParameter(consts.INDEXER_SEARCH_KEY)
	indexerRequest.SearchType = requestWrapper.GetParameter(consts.INDEXER_SEARCH_TYPE)
	seasonNumber := requestWrapper.GetParameter(consts.INDEXER_SEASON_NUMBER)
	indexerRequest.SeasonNumber = seasonNumber
	episodeNumber := requestWrapper.GetParameter(consts.INDEXER_EPISODE_NUMBER)
	indexerRequest.EpisodeNumber = episodeNumber
	offset := requestWrapper.GetParameter(consts.INDEXER_OFFSET)
	indexerRequest.Offset = utils.ParseInt(offset)
	limit := requestWrapper.GetParameter(consts.INDEXER_LIMIT)
	indexerRequest.Limit = utils.ParseInt(limit)
	return indexerRequest
}

func UpdateRequestWrapper(indexerRequest *req.IndexerReq, requestWrapper *services.RequestWrapper) {
	if indexerRequest.SearchKey != "" {
		requestWrapper.SetParameter(consts.INDEXER_SEARCH_KEY, indexerRequest.SearchKey)
	}
	if indexerRequest.SearchType != "" {
		requestWrapper.SetParameter(consts.INDEXER_SEARCH_TYPE, indexerRequest.SearchType)
	}
	if indexerRequest.SeasonNumber != "" {
		requestWrapper.SetParameter(consts.INDEXER_SEASON_NUMBER, indexerRequest.SeasonNumber)
	}
	if indexerRequest.EpisodeNumber != "" {
		requestWrapper.SetParameter(consts.INDEXER_EPISODE_NUMBER, indexerRequest.EpisodeNumber)
	}
	if indexerRequest.Offset != 0 {
		requestWrapper.SetParameter(consts.INDEXER_OFFSET, strconv.Itoa(indexerRequest.Offset))
	}
	if indexerRequest.Limit != 0 {
		requestWrapper.SetParameter(consts.INDEXER_LIMIT, strconv.Itoa(indexerRequest.Limit))
	}
}
