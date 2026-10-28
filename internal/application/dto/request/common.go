package request

import (
	"encoding/json"
	"strconv"
)

type Int64Slice []int64

func (s *Int64Slice) UnmarshalJSON(data []byte) error {
	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err != nil {
		return err
	}

	*s = make([]int64, len(strSlice))
	for i, str := range strSlice {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		(*s)[i] = val
	}

	return nil
}

type DeleteBatchRequest struct {
	IDs Int64Slice `json:"ids" binding:"required,min=1,dive,required"`
}
