package model

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

type Set struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	CardSetTypeCode string `json:"card_set_type_code"`
}

// Card is a placeholder — full struct comes later.
type Card struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	TypeCode    string `json:"type_code"`
	FactionCode string `json:"faction_code"`
	PackCode    string `json:"pack_code"`
}

type PackOwnership struct {
	Pack  Pack `json:"pack"`
	Owned bool `json:"owned"`
}
