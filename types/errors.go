package types

type KeysError interface {
	Error() string
	Keys() []string
	PrependKeys([]string)
}

type KeysLookupError struct {
	msg  string
	keys []string
}

func PrependKeyStringIfKeysError(err error, key string) error {
	return PrependKeysIfKeysError(err, []string{key})
}

func PrependKeysIfKeysError(err error, keys []string) error {
	if keysErr, ok := err.(KeysError); ok {
		keysErr.PrependKeys(keys)
	}

	return err
}

func NewKeysLookupError(msg string, keys []string) KeysError {
	return &KeysLookupError{msg: msg, keys: keys}
}

func (e *KeysLookupError) Error() string {
	return e.msg
}

func (e *KeysLookupError) Keys() []string {
	return e.keys
}

func (e *KeysLookupError) PrependKeys(keys []string) {
	e.keys = append(keys, e.keys...)
}
