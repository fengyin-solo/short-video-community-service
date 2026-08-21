package transcode

type Worker struct{ SideEffects int }

func (w *Worker) Run(key string) { w.SideEffects++ }
