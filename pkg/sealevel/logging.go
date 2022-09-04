package sealevel

type Logger interface {
	Log(s string)
}

type LogRecorder struct {
	Logs []string
}

func (r *LogRecorder) Log(s string) {
	r.Logs = append(r.Logs, s)
}
