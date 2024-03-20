package rollingwindow

type Policy interface {
	Append(data float64)
	Reduce(fu func(window Window) float64)
}
