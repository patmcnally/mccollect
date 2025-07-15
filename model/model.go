package model

// Pack represents a release pack.
type Pack struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	CgdbID       *int    `json:"cgdb_id,omitempty"`
	OctgnID      *string `json:"octgn_id,omitempty"`
	DateRelease  *string `json:"date_release,omitempty"`
	PackTypeCode string  `json:"pack_type_code"`
	Position     *int    `json:"position,omitempty"`
	Size         *int    `json:"size,omitempty"`
}

// Set represents a card set.
type Set struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	CardSetTypeCode string `json:"card_set_type_code"`
}
