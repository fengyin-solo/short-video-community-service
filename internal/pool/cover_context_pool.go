package pool

import "shortvideo/internal/model"

type CoverContextPool struct {
	available chan *model.CoverRenderContext
}

func NewCoverContextPool() *CoverContextPool {
	return &CoverContextPool{available: make(chan *model.CoverRenderContext, 1)}
}

func (p *CoverContextPool) Acquire(requestID, ownerID, videoID string) *model.CoverRenderContext {
	var ctx *model.CoverRenderContext
	select {
	case ctx = <-p.available:
	default:
		ctx = &model.CoverRenderContext{}
	}
	ctx.RequestID = requestID
	ctx.OwnerID = ownerID
	ctx.VideoID = videoID
	return ctx
}

func (p *CoverContextPool) Release(ctx *model.CoverRenderContext) {
	select {
	case p.available <- ctx:
	default:
	}
}
