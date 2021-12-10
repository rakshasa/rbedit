package actions

func NewVerifyAnnounceListFunction() ActionFunc {
	return NewVerifyResultIsListContentFunction(NewVerifyResultIsListContentFunction(NewVerifyResultIsURIFunction()))
}

func NewVerifyAnnounceListCategoryFunction() ActionFunc {
	return NewVerifyResultIsListContentFunction(NewVerifyResultIsURIFunction())
}
