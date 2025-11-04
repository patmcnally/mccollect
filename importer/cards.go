package importer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/patmcnally/mccollect/model"
)

// rawCard uses json.RawMessage for flexible boolean/int parsing and JSON blobs.
type rawCard struct {
	Code        string           `json:"code"`
	OctgnID     *string          `json:"octgn_id"`
	PackCode    string           `json:"pack_code"`
	Position    *int             `json:"position"`
	Quantity    *int             `json:"quantity"`
	SetCode     *string          `json:"set_code"`
	CardSetCode *string          `json:"card_set_code"`
	SetPosition *int             `json:"set_position"`
	DuplicateOf *string          `json:"duplicate_of"`
	Hidden      json.RawMessage  `json:"hidden"`
	BackLink    *string          `json:"back_link"`
	BackName    *string          `json:"back_name"`
	BackText    *string          `json:"back_text"`
	DoubleSided json.RawMessage  `json:"double_sided"`
	TypeCode    *string          `json:"type_code"`
	FactionCode *string          `json:"faction_code"`
	Traits      *string          `json:"traits"`
	IsUnique    json.RawMessage  `json:"is_unique"`
	Permanent   json.RawMessage  `json:"permanent"`
	Spoiler     json.RawMessage  `json:"spoiler"`
	Name        *string          `json:"name"`
	Subname     *string          `json:"subname"`
	Flavor      *string          `json:"flavor"`
	Illustrator *string          `json:"illustrator"`
	Text        *string          `json:"text"`
	Errata      *string          `json:"errata"`
	Cost        *int             `json:"cost"`
	CostPerHero json.RawMessage  `json:"cost_per_hero"`
	CostStar    json.RawMessage  `json:"cost_star"`
	DeckLimit   *int             `json:"deck_limit"`

	Attack     *int            `json:"attack"`
	AttackCost *int            `json:"attack_cost"`
	AttackStar json.RawMessage `json:"attack_star"`
	Thwart     *int            `json:"thwart"`
	ThwartCost *int            `json:"thwart_cost"`
	ThwartStar json.RawMessage `json:"thwart_star"`
	Defense    *int            `json:"defense"`
	DefenseStar json.RawMessage `json:"defense_star"`
	Recover    *int            `json:"recover"`
	RecoverStar json.RawMessage `json:"recover_star"`
	Health     *int            `json:"health"`
	HealthPerHero  json.RawMessage `json:"health_per_hero"`
	HealthPerGroup json.RawMessage `json:"health_per_group"`
	HealthStar     json.RawMessage `json:"health_star"`
	HandSize   *int            `json:"hand_size"`

	ResourceEnergy   *int `json:"resource_energy"`
	ResourceMental   *int `json:"resource_mental"`
	ResourcePhysical *int `json:"resource_physical"`
	ResourceWild     *int `json:"resource_wild"`

	Scheme     *int            `json:"scheme"`
	SchemeText *string         `json:"scheme text"` // note: space in JSON key
	SchemeStar json.RawMessage `json:"scheme_star"`
	Boost      *int            `json:"boost"`
	BoostStar  json.RawMessage `json:"boost_star"`

	BaseThreat         *int            `json:"base_threat"`
	BaseThreatFixed    json.RawMessage `json:"base_threat_fixed"`
	BaseThreatPerGroup json.RawMessage `json:"base_threat_per_group"`
	Threat             *int            `json:"threat"`
	ThreatFixed        json.RawMessage `json:"threat_fixed"`
	ThreatPerGroup     json.RawMessage `json:"threat_per_group"`
	ThreatStar         json.RawMessage `json:"threat_star"`
	EscalationThreat      *int            `json:"escalation_threat"`
	EscalationThreatFixed json.RawMessage `json:"escalation_threat_fixed"`
	EscalationThreatStar  json.RawMessage `json:"escalation_threat_star"`

	Stage              *string `json:"stage"`
	SchemeAcceleration *int `json:"scheme_acceleration"`
	SchemeHazard       *int `json:"scheme_hazard"`
	SchemeCrisis       *int `json:"scheme_crisis"`
	SchemeAmplify      *int `json:"scheme_amplify"`

	DeckOptions      json.RawMessage `json:"deck_options"`
	DeckRequirements json.RawMessage `json:"deck_requirements"`
	Meta             json.RawMessage `json:"meta"`
}

// flag interprets a json.RawMessage as a boolean flag (0 or 1).
// Handles true, false, 1, 0, null, and absent values.
func flag(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 0
	}
	s := strings.TrimSpace(string(raw))
	if s == "true" || s == "1" {
		return 1
	}
	return 0
}

// jsonBlob converts a RawMessage to a *string for storage, or nil if empty/null.
func jsonBlob(raw json.RawMessage) *string {
	if len(raw) == 0 {
		return nil
	}
	s := strings.TrimSpace(string(raw))
	if s == "null" || s == "" || s == "[]" || s == "{}" {
		return nil
	}
	return &s
}

func toCard(r rawCard) model.Card {
	return model.Card{
		Code:        r.Code,
		OctgnID:     r.OctgnID,
		PackCode:    r.PackCode,
		Position:    r.Position,
		Quantity:    r.Quantity,
		SetCode:     r.SetCode,
		CardSetCode: r.CardSetCode,
		SetPosition: r.SetPosition,
		DuplicateOf: r.DuplicateOf,
		Hidden:      flag(r.Hidden),
		BackLink:    r.BackLink,
		BackName:    r.BackName,
		BackText:    r.BackText,
		DoubleSided: flag(r.DoubleSided),
		TypeCode:    deref(r.TypeCode),
		FactionCode: deref(r.FactionCode),
		Traits:      r.Traits,
		IsUnique:    flag(r.IsUnique),
		Permanent:   flag(r.Permanent),
		Spoiler:     flag(r.Spoiler),
		Name:        deref(r.Name),
		Subname:     r.Subname,
		Flavor:      r.Flavor,
		Illustrator: r.Illustrator,
		Text:        r.Text,
		Errata:      r.Errata,
		Cost:        r.Cost,
		CostPerHero: flag(r.CostPerHero),
		CostStar:    flag(r.CostStar),
		DeckLimit:   r.DeckLimit,

		Attack:       r.Attack,
		AttackCost:   r.AttackCost,
		AttackStar:   flag(r.AttackStar),
		Thwart:       r.Thwart,
		ThwartCost:   r.ThwartCost,
		ThwartStar:   flag(r.ThwartStar),
		Defense:      r.Defense,
		DefenseStar:  flag(r.DefenseStar),
		Recover:      r.Recover,
		RecoverStar:  flag(r.RecoverStar),
		Health:       r.Health,
		HealthPerHero:  flag(r.HealthPerHero),
		HealthPerGroup: flag(r.HealthPerGroup),
		HealthStar:     flag(r.HealthStar),
		HandSize:     r.HandSize,

		ResourceEnergy:   r.ResourceEnergy,
		ResourceMental:   r.ResourceMental,
		ResourcePhysical: r.ResourcePhysical,
		ResourceWild:     r.ResourceWild,

		Scheme:     r.Scheme,
		SchemeText: r.SchemeText,
		SchemeStar: flag(r.SchemeStar),
		Boost:      r.Boost,
		BoostStar:  flag(r.BoostStar),

		BaseThreat:         r.BaseThreat,
		BaseThreatFixed:    flag(r.BaseThreatFixed),
		BaseThreatPerGroup: flag(r.BaseThreatPerGroup),
		Threat:             r.Threat,
		ThreatFixed:        flag(r.ThreatFixed),
		ThreatPerGroup:     flag(r.ThreatPerGroup),
		ThreatStar:         flag(r.ThreatStar),
		EscalationThreat:      r.EscalationThreat,
		EscalationThreatFixed: flag(r.EscalationThreatFixed),
		EscalationThreatStar:  flag(r.EscalationThreatStar),

		Stage:              r.Stage,
		SchemeAcceleration: r.SchemeAcceleration,
		SchemeHazard:       r.SchemeHazard,
		SchemeCrisis:       r.SchemeCrisis,
		SchemeAmplify:      r.SchemeAmplify,

		DeckOptions:      jsonBlob(r.DeckOptions),
		DeckRequirements: jsonBlob(r.DeckRequirements),
		Meta:             jsonBlob(r.Meta),
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// LoadPackFile reads a single pack JSON file and returns the parsed cards.
// Skips reprint stubs that have no type_code or name.
func LoadPackFile(path string) ([]model.Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw []rawCard
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var cards []model.Card
	for _, r := range raw {
		// Skip reprint stubs
		if r.TypeCode == nil || r.Name == nil {
			continue
		}
		cards = append(cards, toCard(r))
	}
	return cards, nil
}

// LoadAllCards reads all pack/*.json files from the data root.
func LoadAllCards(dataRoot string) ([]model.Card, error) {
	packDir := filepath.Join(dataRoot, "pack")
	entries, err := os.ReadDir(packDir)
	if err != nil {
		return nil, err
	}

	// Sort for deterministic order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var allCards []model.Card
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		cards, err := LoadPackFile(filepath.Join(packDir, e.Name()))
		if err != nil {
			return nil, err
		}
		allCards = append(allCards, cards...)
	}
	return allCards, nil
}
