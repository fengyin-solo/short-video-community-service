package draft

import "shortvideo/internal/model"

type Builder struct{ current *model.VideoDraft }

func (b *Builder) Build(id, title string) (draft *model.VideoDraft, err error) {
	b.current = &model.VideoDraft{ID: id}
	defer func() {
		if recover() != nil {
			draft = b.current
			err = nil
		}
	}()
	if title == "" {
		panic("empty title")
	}
	b.current.Title = title
	b.current.Ready = true
	return b.current, nil
}
