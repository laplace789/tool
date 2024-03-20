package rollingwindow

import (
	"testing"
	"time"
)

func TestTimeWindow(t *testing.T) {
	var bucketSize = time.Millisecond * 100
	var numberBuckets = 10
	var w = NewWindow(numberBuckets)
	var p = NewTimePolicy(w, bucketSize)
	for x := 0; x < numberBuckets; x = x + 1 {
		p.Append(1)
		time.Sleep(bucketSize)
	}
	var final = p.Reduce(func(w Window) float64 {
		var result float64
		for _, bucket := range w {
			for _, point := range bucket {
				result = result + point
			}
		}
		return result
	})
	if final != float64(numberBuckets) {
		t.Fatal(final)
	}

	for x := 0; x < numberBuckets; x = x + 1 {
		p.Append(2)
		time.Sleep(bucketSize)
	}

	final = p.Reduce(func(w Window) float64 {
		var result float64
		for _, bucket := range w {
			for _, point := range bucket {
				result = result + point
			}
		}
		return result
	})
	if final != 2*float64(numberBuckets) {
		t.Fatal("got", final, "expected", 2*float64(numberBuckets))
	}
}
