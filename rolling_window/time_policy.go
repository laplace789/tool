package rollingwindow

import (
	"sync"
	"time"
)

type TimePolicy struct {
	bucketSize        time.Duration //每個bucket 可以裝多少時限內資料
	bucketSizeNano    int64
	numberOfBuckets   int //有幾個bucket
	numberOfBuckets64 int64
	window            Window
	lastWindowOffset  int
	lastWindowTime    int64
	lock              sync.RWMutex
}

func NewTimePolicy(window Window, bucketDuration time.Duration) *TimePolicy {
	return &TimePolicy{
		bucketSize:        bucketDuration,
		bucketSizeNano:    bucketDuration.Nanoseconds(),
		numberOfBuckets:   len(window),
		numberOfBuckets64: int64(len(window)),
		window:            window,
	}
}

// selectBucket 要加入window之前 都要先選擇在哪個bucket
func (w *TimePolicy) selectBucket(currentTime time.Time) (int64, int) {

	currentTimeNano := currentTime.UnixNano()

	var adjustedTime = currentTimeNano / w.bucketSizeNano

	//決定內部bucket的index,如果windowOffset < lastWindowOffset -> 開始複寫就資料
	var windowOffset = int(adjustedTime % w.numberOfBuckets64)
	return adjustedTime, windowOffset
}

// resetWindow 重置整個window
func (w *TimePolicy) resetWindow() {
	for offset := range w.window {
		w.window[offset] = w.window[offset][:0]
	}
}

func (w *TimePolicy) resetBuckets(windowOffset int) {
	var distance = windowOffset - w.lastWindowOffset
	// If the distance between current and last is negative then we've wrapped
	// around the ring. Recalculate the distance.
	if distance < 0 {
		distance = (w.numberOfBuckets - w.lastWindowOffset) + windowOffset
	}
	for counter := 1; counter < distance; counter = counter + 1 {
		var offset = (counter + w.lastWindowOffset) % w.numberOfBuckets
		w.window[offset] = w.window[offset][:0]
	}
}

func (w *TimePolicy) keepConsistent(adjustedTime int64, windowOffset int) {
	if adjustedTime-w.lastWindowTime > w.numberOfBuckets64 {
		w.resetWindow()
	}

	// When one or more buckets are missed we need to zero them out.
	if adjustedTime != w.lastWindowTime && adjustedTime-w.lastWindowTime < w.numberOfBuckets64 {
		w.resetBuckets(windowOffset)
	}
}

// AppendWithTimestamp same as Append but with timestamp as parameter
func (w *TimePolicy) AppendWithTimestamp(value float64, timestamp time.Time) {
	w.lock.Lock()
	defer w.lock.Unlock()

	var adjustedTime, windowOffset = w.selectBucket(timestamp)
	w.keepConsistent(adjustedTime, windowOffset)
	if w.lastWindowOffset != windowOffset {
		w.window[windowOffset] = []float64{value}
	} else {
		w.window[windowOffset] = append(w.window[windowOffset], value)
	}
	w.lastWindowTime = adjustedTime
	w.lastWindowOffset = windowOffset
}

// Append a value to the window using a time bucketing strategy.
func (w *TimePolicy) Append(value float64) {
	w.AppendWithTimestamp(value, time.Now())
}

// Reduce the window to a single value using a reduction function.
func (w *TimePolicy) Reduce(f func(Window) float64) float64 {
	w.lock.Lock()
	defer w.lock.Unlock()

	var adjustedTime, windowOffset = w.selectBucket(time.Now())
	w.keepConsistent(adjustedTime, windowOffset)
	return f(w.window)
}
