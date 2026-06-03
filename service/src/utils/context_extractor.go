package utils

import (
	"context"
	"strconv"
)

const USER_ID_CONTEXT_KEY = "UserID"

func GetAuthenticatedUserID(ctx context.Context) (int64, error) {
	raw, ok := ctx.Value(USER_ID_CONTEXT_KEY).(string)
	if !ok {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(raw, 10, 64)
}
