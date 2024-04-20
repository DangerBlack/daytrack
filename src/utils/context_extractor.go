package utils

import (
	"context"
	"strconv"
)

const USER_ID_CONTEXT_KEY = "UserID"

func GetAuthenticatedUserID(ctx context.Context) (int64, error) {
	var err error
	var userID int64

	if userID, err = strconv.ParseInt(ctx.Value(USER_ID_CONTEXT_KEY).(string), 10, 64); err != nil {
		return 0, err
	}

	return userID, nil
}
