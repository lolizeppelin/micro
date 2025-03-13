package log

var metric LoggingMetric

type LoggingMetric interface {
	Error()
	Warn()
}

type fakeMetric struct {
}

func (fakeMetric) Error() {

}

func (fakeMetric) Warn() {

}

func init() {
	metric = fakeMetric{}
}
