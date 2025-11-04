package model

// Pack represents a release pack (core set, hero pack, scenario pack, etc.)
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

// Set represents a card set (hero set, villain set, modular set, etc.)
type Set struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	CardSetTypeCode string `json:"card_set_type_code"`
}

// Card represents a single card with all 72 database columns.
// JSON blob fields (deck_options, deck_requirements, meta) are stored as raw JSON strings.
type Card struct {
	Code        string  `json:"code"`
	OctgnID     *string `json:"octgn_id,omitempty"`
	PackCode    string  `json:"pack_code"`
	Position    *int    `json:"position,omitempty"`
	Quantity    *int    `json:"quantity,omitempty"`
	SetCode     *string `json:"set_code,omitempty"`
	CardSetCode *string `json:"card_set_code,omitempty"`
	SetPosition *int    `json:"set_position,omitempty"`
	DuplicateOf *string `json:"duplicate_of,omitempty"`
	Hidden      int     `json:"hidden"`
	BackLink    *string `json:"back_link,omitempty"`
	BackName    *string `json:"back_name,omitempty"`
	BackText    *string `json:"back_text,omitempty"`
	DoubleSided int     `json:"double_sided"`

	TypeCode    string  `json:"type_code"`
	FactionCode string  `json:"faction_code"`
	Traits      *string `json:"traits,omitempty"`
	IsUnique    int     `json:"is_unique"`
	Permanent   int     `json:"permanent"`
	Spoiler     int     `json:"spoiler"`

	Name        string  `json:"name"`
	Subname     *string `json:"subname,omitempty"`
	Flavor      *string `json:"flavor,omitempty"`
	Illustrator *string `json:"illustrator,omitempty"`
	Text        *string `json:"text,omitempty"`
	Errata      *string `json:"errata,omitempty"`
	Cost        *int    `json:"cost,omitempty"`
	CostPerHero int     `json:"cost_per_hero"`
	CostStar    int     `json:"cost_star"`
	DeckLimit   *int    `json:"deck_limit,omitempty"`

	// Player stats
	Attack       *int `json:"attack,omitempty"`
	AttackCost   *int `json:"attack_cost,omitempty"`
	AttackStar   int  `json:"attack_star"`
	Thwart       *int `json:"thwart,omitempty"`
	ThwartCost   *int `json:"thwart_cost,omitempty"`
	ThwartStar   int  `json:"thwart_star"`
	Defense      *int `json:"defense,omitempty"`
	DefenseStar  int  `json:"defense_star"`
	Recover      *int `json:"recover,omitempty"`
	RecoverStar  int  `json:"recover_star"`
	Health       *int `json:"health,omitempty"`
	HealthPerHero  int `json:"health_per_hero"`
	HealthPerGroup int `json:"health_per_group"`
	HealthStar   int  `json:"health_star"`
	HandSize     *int `json:"hand_size,omitempty"`
	ResourceEnergy   *int `json:"resource_energy,omitempty"`
	ResourceMental   *int `json:"resource_mental,omitempty"`
	ResourcePhysical *int `json:"resource_physical,omitempty"`
	ResourceWild     *int `json:"resource_wild,omitempty"`

	// Encounter stats
	Scheme      *int    `json:"scheme,omitempty"`
	SchemeText  *string `json:"scheme_text,omitempty"`
	SchemeStar  int     `json:"scheme_star"`
	Boost       *int    `json:"boost,omitempty"`
	BoostStar   int     `json:"boost_star"`
	BaseThreat         *int `json:"base_threat,omitempty"`
	BaseThreatFixed    int  `json:"base_threat_fixed"`
	BaseThreatPerGroup int  `json:"base_threat_per_group"`
	Threat             *int `json:"threat,omitempty"`
	ThreatFixed        int  `json:"threat_fixed"`
	ThreatPerGroup     int  `json:"threat_per_group"`
	ThreatStar         int  `json:"threat_star"`
	EscalationThreat      *int `json:"escalation_threat,omitempty"`
	EscalationThreatFixed int  `json:"escalation_threat_fixed"`
	EscalationThreatStar  int     `json:"escalation_threat_star"`
	Stage               *string `json:"stage,omitempty"`
	SchemeAcceleration   *int `json:"scheme_acceleration,omitempty"`
	SchemeHazard         *int `json:"scheme_hazard,omitempty"`
	SchemeCrisis         *int `json:"scheme_crisis,omitempty"`
	SchemeAmplify        *int `json:"scheme_amplify,omitempty"`

	// JSON blobs stored as raw strings
	DeckOptions      *string `json:"deck_options,omitempty"`
	DeckRequirements *string `json:"deck_requirements,omitempty"`
	Meta             *string `json:"meta,omitempty"`
}

// Collection represents a named collection (e.g. one per player).
type Collection struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// PackOwnership pairs a pack with its ownership status for display.
type PackOwnership struct {
	Pack  Pack `json:"pack"`
	Owned bool `json:"owned"`
}
