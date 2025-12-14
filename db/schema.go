package db

const schema = `
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS packs (
    code            TEXT PRIMARY KEY,
    name            TEXT NOT NULL,
    cgdb_id         INTEGER,
    octgn_id        TEXT,
    date_release    TEXT,
    pack_type_code  TEXT NOT NULL,
    position        INTEGER,
    size            INTEGER
);

CREATE TABLE IF NOT EXISTS sets (
    code                TEXT PRIMARY KEY,
    name                TEXT NOT NULL,
    card_set_type_code  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cards (
    code                    TEXT PRIMARY KEY,
    octgn_id                TEXT,
    pack_code               TEXT NOT NULL REFERENCES packs(code),
    position                INTEGER,
    quantity                INTEGER,
    set_code                TEXT,
    card_set_code           TEXT,
    set_position            INTEGER,
    duplicate_of            TEXT,
    hidden                  INTEGER NOT NULL DEFAULT 0,
    back_link               TEXT,
    back_name               TEXT,
    back_text               TEXT,
    double_sided            INTEGER NOT NULL DEFAULT 0,
    type_code               TEXT NOT NULL,
    faction_code            TEXT NOT NULL,
    traits                  TEXT,
    is_unique               INTEGER NOT NULL DEFAULT 0,
    permanent               INTEGER NOT NULL DEFAULT 0,
    spoiler                 INTEGER NOT NULL DEFAULT 0,
    name                    TEXT NOT NULL,
    subname                 TEXT,
    flavor                  TEXT,
    illustrator             TEXT,
    text                    TEXT,
    errata                  TEXT,
    cost                    INTEGER,
    cost_per_hero           INTEGER NOT NULL DEFAULT 0,
    cost_star               INTEGER NOT NULL DEFAULT 0,
    deck_limit              INTEGER,
    attack                  INTEGER,
    attack_cost             INTEGER,
    attack_star             INTEGER NOT NULL DEFAULT 0,
    thwart                  INTEGER,
    thwart_cost             INTEGER,
    thwart_star             INTEGER NOT NULL DEFAULT 0,
    defense                 INTEGER,
    defense_star            INTEGER NOT NULL DEFAULT 0,
    recover                 INTEGER,
    recover_star            INTEGER NOT NULL DEFAULT 0,
    health                  INTEGER,
    health_per_hero         INTEGER NOT NULL DEFAULT 0,
    health_per_group        INTEGER NOT NULL DEFAULT 0,
    health_star             INTEGER NOT NULL DEFAULT 0,
    hand_size               INTEGER,
    resource_energy         INTEGER,
    resource_mental         INTEGER,
    resource_physical       INTEGER,
    resource_wild           INTEGER,
    scheme                  INTEGER,
    scheme_text             TEXT,
    scheme_star             INTEGER NOT NULL DEFAULT 0,
    boost                   INTEGER,
    boost_star              INTEGER NOT NULL DEFAULT 0,
    base_threat             INTEGER,
    base_threat_fixed       INTEGER NOT NULL DEFAULT 0,
    base_threat_per_group   INTEGER NOT NULL DEFAULT 0,
    threat                  INTEGER,
    threat_fixed            INTEGER NOT NULL DEFAULT 0,
    threat_per_group        INTEGER NOT NULL DEFAULT 0,
    threat_star             INTEGER NOT NULL DEFAULT 0,
    escalation_threat               INTEGER,
    escalation_threat_fixed         INTEGER NOT NULL DEFAULT 0,
    escalation_threat_star          INTEGER NOT NULL DEFAULT 0,
    stage                   TEXT,
    scheme_acceleration     INTEGER,
    scheme_hazard           INTEGER,
    scheme_crisis           INTEGER,
    scheme_amplify          INTEGER,
    deck_options            TEXT,
    deck_requirements       TEXT,
    meta                    TEXT
);

CREATE INDEX IF NOT EXISTS idx_cards_octgn_id  ON cards (octgn_id);
CREATE INDEX IF NOT EXISTS idx_cards_pack_code ON cards (pack_code);
CREATE INDEX IF NOT EXISTS idx_cards_type_code ON cards (type_code);
CREATE INDEX IF NOT EXISTS idx_cards_name      ON cards (name COLLATE NOCASE);

CREATE TABLE IF NOT EXISTS collections (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS collection_packs (
    collection_id   INTEGER NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    pack_code       TEXT    NOT NULL REFERENCES packs(code),
    owned           INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (collection_id, pack_code)
);

CREATE INDEX IF NOT EXISTS idx_collection_packs_owned ON collection_packs (collection_id, owned);

CREATE TABLE IF NOT EXISTS _meta (
    key     TEXT PRIMARY KEY,
    value   TEXT NOT NULL
);
`
