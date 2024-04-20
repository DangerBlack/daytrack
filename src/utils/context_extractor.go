package utils

import (
	"context"
	"strconv"
)

const USER_ID_CONTEXT_KEY = "UserID"

func GetAuthenticatedUserID(ctx context.Context) (int, error) {
	var err error
	var userID int

	if userID, err = strconv.Atoi(ctx.Value(USER_ID_CONTEXT_KEY).(string)); err != nil {
		return 0, err
	}

	return userID, nil
}
