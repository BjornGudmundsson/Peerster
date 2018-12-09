package data

import (
	"strings"
	"time"
)

//HalfASecond is a constant that represents half a second in go conde.
const HalfASecond time.Duration = time.Duration(500 * time.Millisecond)

//RecentRequests is map that maps keywords and origins
//to the last time this gossiper received this request
type RecentRequests map[string]time.Time

//NewRecentRequests returns a new instance of
//a RecentRequests type that is empty.
func NewRecentRequests() RecentRequests {
	return RecentRequests(make(map[string]time.Time))
}

//AddSearchRequest updates the RecentRequest type with a new searchrequest.
//If it already has the SearchRequest it will update the current value.
//AddSearchRequest also notifies the caller if more than half a second has
//passed since the last time this request was received.
func (rr RecentRequests) AddSearchRequest(req *SearchRequest) bool {
	now := time.Now()
	src := req.Origin
	keywords := strings.Join(req.Keywords, ",")
	temp := []string{src, keywords}
	key := strings.Join(temp, "-")
	last, ok := rr[key]
	if !ok {
		rr[key] = now
		return true
	}
	if time.Since(last) > HalfASecond {
		rr[key] = now
		return true
	}
	rr[key] = now
	return false
}
