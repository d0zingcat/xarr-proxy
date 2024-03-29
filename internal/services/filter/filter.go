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
		count := 0
		log.Debug().
			Any("title", title).
			Any("searchTitleList", searchTitleList).
			Any("offsetKey", offsetKey).
			Any("offsetList", offsetList).
			Any("index", index).
			Msg("first request")
		for {
			if size > 1 && index == size-1 {
				// 已查询到的结果数量少于 8 则去除季集信息尝试查询
				if offset < 8 {
					// 只查询 limit - 1 条记录
					indexerRequest.Limit = indexerRequest.Limit - 1
					if indexerRequest.SeasonNumber != "" {
						indexerRequest.SeasonNumber = ""
						indexerRequest.EpisodeNumber = ""
					} else {
						indexerRequest.SearchKey = utils.RemoveSeasonEpisode(indexerRequest.SearchKey)
					}
				} else {
					break
				}
			}
			// 请求
			i.Service.UpdateIndexerRequest(index, searchTitleList, offsetList, indexerRequest)
			UpdateRequestWrapper(indexerRequest, requestWrapper)
			newXml := i.Service.ExecuteNewRequest(requestWrapper)
			count = utils.XmlCount(newXml)
			if count > 0 || len(xml) == 0 {
				xml = utils.XmlMerge(xml, newXml)
			}
			index++
			if index >= size {
				break
			}
			// 更新参数
			offset = offset + count
			indexerRequest.Offset = offset
			indexerRequest.Limit = indexerRequest.Limit - count
			offsetList[index] = offset
			i.Service.UpdateOffsetList(offsetKey, offsetList)

			if indexerRequest.Limit-count <= 0 {
				break
			}
		}
	} else {
		xml = i.Service.ExecuteNewRequest(requestWrapper)
	}
	// 处理结果
	log.Logger.Debug().Msgf("before: %v", xml)
	xml = i.Service.ExecuteFormatRule(xml)
	log.Logger.Debug().Msgf("after: %v", xml)
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
