package cli

import "github.com/mebiusashan/beaker/internal/common"

func responseString(data interface{}, field string) string {
	value, ok := data.(string)
	if !ok || value == "" {
		common.ErrCode(common.ErrorCodeInvalidResponse, "missing or invalid "+field)
	}
	return value
}
