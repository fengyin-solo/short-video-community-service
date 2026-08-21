package tagindex

type Job struct {
	Tags    []string
	Release <-chan struct{}
	Done    chan<- []string
}

func (j Job) Run() {
	<-j.Release
	j.Done <- j.Tags
}
