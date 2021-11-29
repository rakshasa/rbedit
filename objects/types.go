package objects

import (
	"fmt"
)

type AnnounceCategory struct {
	uris []string
}

type AnnounceList struct {
	categories []*AnnounceCategory
}

func NewAnnounceList(announceListObj interface{}) (AnnounceList, error) {
	categories, ok := AsList(announceListObj)
	if !ok {
		return AnnounceList{}, fmt.Errorf("object is not a list")
	}

	announceList := AnnounceList{}

	for cIdx, categoryObj := range categories {
		uriList, ok := AsList(categoryObj)
		if !ok {
			return AnnounceList{}, fmt.Errorf("category %d is not a list", cIdx)
		}

		category := AnnounceCategory{}

		for uIdx, uriObj := range uriList {
			uri, ok := AsAbsoluteURI(uriObj)
			if !ok {
				return AnnounceList{}, fmt.Errorf("index %d:%d is not a valid URI", cIdx, uIdx)
			}

			category.uris = append(category.uris, uri)
		}

		announceList.categories = append(announceList.categories, &category)
	}

	return announceList, nil
}

func (a *AnnounceCategory) URIs() []string {
	return a.uris
}

func (a *AnnounceCategory) AppendURI(uri string) {
	a.uris = append(a.uris, uri)
}

func (a *AnnounceCategory) ToListObject() []interface{} {
	var uriListObj []interface{}

	for _, uri := range a.uris {
		uriListObj = append(uriListObj, uri)
	}

	return uriListObj
}

func (a *AnnounceList) Categories() []*AnnounceCategory {
	return a.categories
}

func (a *AnnounceList) ToListObject() []interface{} {
	var announceListObj []interface{}

	for _, category := range a.categories {
		announceListObj = append(announceListObj, category.ToListObject())
	}

	return announceListObj
}
