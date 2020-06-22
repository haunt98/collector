package clock

import (
	"time"
)

func NowDateInSaiGon() (string, error) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return "", err
	}

	now := time.Now().In(loc)

	layout := "2006-01-02"
	return now.Format(layout), nil
}
