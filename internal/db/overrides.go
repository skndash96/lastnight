package db

type TagOptions []struct {
	ID    int32  `json:"id"`
	Value string `json:"value"`
}
type Tag struct {
	KeyID   int32      `json:"key_id"`
	Key     string     `json:"key"`
	ValueID int32      `json:"value_id,omitempty" extensions:"x-nullable"`
	Value   string     `json:"value,omitempty" extensions:"x-nullable"`
	Options TagOptions `json:"options,omitempty"`
}
