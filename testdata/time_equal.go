package pkg

import "time"

func t() bool {
	t := time.Now()
	u := t

	if !t.After(u) {
		return t == u // MATCH /use t.Equal(u) instead of "==" operator/
	}

	return t != u // MATCH /use !t.Equal(u) instead of "!=" operator/
}

// issue #846
func isNow(t time.Time) bool    { return t == time.Now() } // MATCH /use t.Equal(time.Now()) instead of "==" operator/
func isNotNow(t time.Time) bool { return time.Now() != t } // MATCH /use !time.Now().Equal(t) instead of "!=" operator/
