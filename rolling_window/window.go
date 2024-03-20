package rollingwindow

//source: https://github.com/asecurityteam/rolling

// Window  policy底下的基礎結構
type Window [][]float64

func NewWindow(buckets int) Window {
	return make([][]float64, buckets)
}
