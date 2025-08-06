package db

const schema = `
CREATE TABLE IF NOT EXISTS packs (
    code           TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    cgdb_id        INTEGER,
    octgn_id       TEXT,
    date_release   TEXT,
    pack_type_code TEXT NOT NULL,
    position       INTEGER,
    size           INTEGER
);
CREATE TABLE IF NOT EXISTS sets (
    code               TEXT PRIMARY KEY,
    name               TEXT NOT NULL,
    card_set_type_code TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS cards (
    code         TEXT PRIMARY KEY,
    pack_code    TEXT NOT NULL REFERENCES packs(code),
    name         TEXT NOT NULL,
    type_code    TEXT NOT NULL,
    faction_code TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS collections (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE TABLE IF NOT EXISTS collection_packs (
    collection_id INTEGER NOT NULL REFERENCES collections(id),
    pack_code     TEXT    NOT NULL REFERENCES packs(code),
    owned         INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (collection_id, pack_code)
);
`
