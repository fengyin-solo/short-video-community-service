package store

type SubtitleTransaction struct {
	staged  []string
	visible []string
}

func NewSubtitleTransaction() *SubtitleTransaction { return &SubtitleTransaction{} }

func (t *SubtitleTransaction) Stage(cue string) { t.staged = append(t.staged, cue) }

func (t *SubtitleTransaction) Commit() error {
	t.visible = append(t.visible, t.staged...)
	t.staged = nil
	return nil
}

func (t *SubtitleTransaction) Visible() []string {
	return append([]string(nil), t.visible...)
}
