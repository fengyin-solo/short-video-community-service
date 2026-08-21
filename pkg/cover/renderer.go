package cover

import "shortvideo/internal/model"

type Renderer struct {
	started chan<- string
}

func NewRenderer(started chan<- string) *Renderer {
	return &Renderer{started: started}
}

func (r *Renderer) Render(ctx *model.CoverRenderContext, release <-chan struct{}) <-chan model.CoverRenderResult {
	done := make(chan model.CoverRenderResult, 1)
	go func() {
		if r.started != nil {
			r.started <- ctx.RequestID
		}
		if release != nil {
			<-release
		}
		done <- model.CoverRenderResult{
			RequestID: ctx.RequestID,
			OwnerID:   ctx.OwnerID,
			VideoID:   ctx.VideoID,
		}
	}()
	return done
}
