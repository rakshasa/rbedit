package actions

func NewVerifyAnnounceList() ActionFunc {
	return NewVerifyResultIsListContent(NewVerifyResultIsListContent(NewVerifyResultIsURI()))
}

func NewVerifyAnnounceListCategory() ActionFunc {
	return NewVerifyResultIsListContent(NewVerifyResultIsURI())
}
