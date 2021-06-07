package data

import (
	"fmt"
	"strconv"
)

type Pages int32

func (r Pages) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d pages", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
