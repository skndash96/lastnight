package db

type TagOptions []struct {
	ID    int32  `json:"id"`
	Value string `json:"value"`
}
type Tag struct {
	KeyID   int32      `json:"key_id"`
	Key     string     `json:"key"`
	// pointers because they can be null
	ValueID *int32     `json:"value_id" extensions:"x-nullable"`
	Value   *string    `json:"value" extensions:"x-nullable"`
	Options TagOptions `json:"options"`
}
