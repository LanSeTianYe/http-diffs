package task

import (
	"strconv"
	"sync/atomic"
	"time"
)

type StatisticsInfo struct {
	// totalCount 总请求数量
	totalCount int64

	// startTime 任务开始时间
	startTime time.Time

	//failedCount 失败的数量
	failedCount *atomic.Int64

	//diffCount 有diff的数量
	diffCount *atomic.Int64

	//sameCount 没有diff的数量
	sameCount *atomic.Int64

	// lastStatisticsTime 上次统计时间
	lastStatisticsTime time.Time
	// lastStatisticsCount 上次统计的数量
	lastStatisticsCount int64
}

func NewStatisticsInfo(totalCount int) *StatisticsInfo {

	s := &StatisticsInfo{
		totalCount:  int64(totalCount),
		startTime:   time.Now(),
		failedCount: &atomic.Int64{},
		diffCount:   &atomic.Int64{},
		sameCount:   &atomic.Int64{},
	}

	s.failedCount.Store(0)
	s.diffCount.Store(0)
	s.sameCount.Store(0)

	return s
}

func (s *StatisticsInfo) AddFailed() {
	s.failedCount.Add(1)
}

func (s *StatisticsInfo) AddDiff() {
	s.diffCount.Add(1)
}

func (s *StatisticsInfo) AddSame() {
	s.sameCount.Add(1)
}

func (s *StatisticsInfo) GetTotalCount() int64 {
	return s.totalCount
}

func (s *StatisticsInfo) GetTimeCost() string {
	return strconv.FormatFloat(time.Since(s.startTime).Seconds(), 'f', 0, 64) + "s"
}

func (s *StatisticsInfo) GetTimeLeft() string {
	runTime := time.Since(s.startTime).Seconds()
	timeLeft := float64(s.GetTotalCount()-s.GetProcessedCount()) * runTime / float64(s.GetProcessedCount())

	return strconv.FormatFloat(timeLeft, 'f', 0, 64) + "s"
}

func (s *StatisticsInfo) GetFailedCount() int64 {
	return s.failedCount.Load()
}

func (s *StatisticsInfo) GetDiffCount() int64 {
	return s.diffCount.Load()
}

func (s *StatisticsInfo) GetSameCount() int64 {
	return s.sameCount.Load()
}

func (s *StatisticsInfo) GetProcessedCount() int64 {
	return s.GetSameCount() + s.GetFailedCount() + s.GetDiffCount()
}

func (s *StatisticsInfo) GetProgress() string {
	progressFloat := float64(s.GetProcessedCount()) / float64(s.GetTotalCount())
	return strconv.FormatFloat(progressFloat*100, 'f', 2, 64) + "%"
}

func (s *StatisticsInfo) GetRate() string {
	rateFloat := float64(s.GetProcessedCount()-s.lastStatisticsCount) / time.Since(s.lastStatisticsTime).Seconds()
	return strconv.FormatFloat(rateFloat, 'f', 0, 64) + " req/s"
}

func (s *StatisticsInfo) ResetLastStatisticsInfo() {
	s.lastStatisticsCount = s.GetProcessedCount()
	s.lastStatisticsTime = time.Now()
}
