package db

import (
	"database/sql"
	"strings"

	"github.com/patmcnally/mccollect/model"
)

// UpsertCards inserts or replaces cards within a transaction.
func (d *DB) UpsertCards(tx *sql.Tx, cards []model.Card) error {
	const cols = 72
	placeholders := "(" + strings.Repeat("?,", cols-1) + "?)"
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO cards VALUES " + placeholders)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, c := range cards {
		if _, err := stmt.Exec(cardRow(c)...); err != nil {
			return err
		}
	}
	return nil
}

// DeleteCardsByPack removes all cards for a given pack_code within a transaction.
func (d *DB) DeleteCardsByPack(tx *sql.Tx, packCode string) error {
	_, err := tx.Exec("DELETE FROM cards WHERE pack_code = ?", packCode)
	return err
}

// cardRow maps a Card struct to the 72-column insert parameter slice.
func cardRow(c model.Card) []any {
	return []any{
		c.Code, c.OctgnID, c.PackCode, c.Position, c.Quantity,
		c.SetCode, c.CardSetCode, c.SetPosition, c.DuplicateOf,
		c.Hidden, c.BackLink, c.BackName, c.BackText, c.DoubleSided,
		c.TypeCode, c.FactionCode, c.Traits, c.IsUnique, c.Permanent, c.Spoiler,
		c.Name, c.Subname, c.Flavor, c.Illustrator, c.Text, c.Errata,
		c.Cost, c.CostPerHero, c.CostStar, c.DeckLimit,
		// player stats
		c.Attack, c.AttackCost, c.AttackStar,
		c.Thwart, c.ThwartCost, c.ThwartStar,
		c.Defense, c.DefenseStar,
		c.Recover, c.RecoverStar,
		c.Health, c.HealthPerHero, c.HealthPerGroup, c.HealthStar,
		c.HandSize,
		c.ResourceEnergy, c.ResourceMental, c.ResourcePhysical, c.ResourceWild,
		// encounter stats
		c.Scheme, c.SchemeText, c.SchemeStar,
		c.Boost, c.BoostStar,
		c.BaseThreat, c.BaseThreatFixed, c.BaseThreatPerGroup,
		c.Threat, c.ThreatFixed, c.ThreatPerGroup, c.ThreatStar,
		c.EscalationThreat, c.EscalationThreatFixed, c.EscalationThreatStar,
		c.Stage,
		c.SchemeAcceleration, c.SchemeHazard, c.SchemeCrisis, c.SchemeAmplify,
		// json blobs
		c.DeckOptions, c.DeckRequirements, c.Meta,
	}
}
