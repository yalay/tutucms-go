package controllers

import (
	"conf"
	"log"
	"strings"

	"github.com/huichen/sego"
)

const (
	MinKeywordLen = 4
	MaxKeywordNum = 5
)

var gKeywordsHandler *KeywordsHandler

func init() {
	gKeywordsHandler = newKeywordsHandler(conf.GetDictPath())
	if gKeywordsHandler == nil {
		log.Panicf("KeywordsHandler is nil")
	}
}

type KeywordsHandler struct {
	seg sego.Segmenter
}

func newKeywordsHandler(dict string) *KeywordsHandler {
	if dict == "" {
		return nil
	}

	var seg sego.Segmenter
	seg.LoadDictionary(dict)
	segDict := seg.Dictionary()
	if segDict == nil || segDict.NumTokens() == 0 {
		return nil
	}

	return &KeywordsHandler{
		seg: seg,
	}
}

func GetKeywords(text string) string {
	return gKeywordsHandler.GetKeywords(text)
}

func (handler *KeywordsHandler) GetKeywords(text string) string {
	segments := handler.seg.Segment([]byte(text))
	segStr := sego.SegmentsToString(segments, false)
	if segStr == "" {
		return ""
	}

	keywords := make([]string, 0)
	existFlag := make(map[string]bool)
	segStrs := strings.Fields(segStr)
	for i, keywordAttr := range segStrs {
		if !isNouns(keywordAttr) {
			continue
		}

		keyword := removeTail(keywordAttr)
		if i > 0 {
			lastkeywordAttr := segStrs[i-1]
			if isAdjectiveWord(lastkeywordAttr) || isVerb(lastkeywordAttr) {
				keyword = removeTail(lastkeywordAttr) + removeTail(keyword)
			}
		}

		if len(keyword) < MinKeywordLen {
			continue
		}

		if existFlag[keyword] {
			continue
		} else {
			existFlag[keyword] = true
		}

		keywords = append(keywords, keyword)
		if len(keywords) > MaxKeywordNum {
			break
		}
	}

	return strings.Join(keywords, ",")
}

func isNouns(keyword string) bool {
	if strings.HasSuffix(keyword, "/n") {
		return true
	}
	return false
}

func isAdjectiveWord(keyword string) bool {
	if strings.HasSuffix(keyword, "/a") {
		return true
	}
	return false
}

func isVerb(keyword string) bool {
	if strings.HasSuffix(keyword, "/v") {
		return true
	}
	return false
}

func splitKeyword(keyword string) (string, string) {
	keywordFields := strings.Split(keyword, "/")
	if len(keywordFields) < 2 {
		return keyword, ""
	}
	return keywordFields[0], keywordFields[len(keywordFields)]
}

func removeTail(keyword string) string {
	tailIdx := strings.Index(keyword, "/")
	if tailIdx > 0 {
		return keyword[:strings.Index(keyword, "/")]
	}
	return keyword
}
