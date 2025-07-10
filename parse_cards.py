#!/usr/bin/env python3
"""Dump Marvel Champions card data to CSV, with optional filtering."""
import json, csv, sys, argparse
from pathlib import Path

def load_cards(data_root, pack_code=None, type_code=None):
    pack_dir = Path(data_root) / "pack"
    cards = []
    for f in sorted(pack_dir.glob("*.json")):
        with open(f) as fp:
            for c in json.load(fp):
                if pack_code and c.get("pack_code") != pack_code:
                    continue
                if type_code and c.get("type_code") != type_code:
                    continue
                cards.append(c)
    return cards

FIELDS = ["code","name","type_code","faction_code","pack_code","cost","attack","thwart","defense","health","text"]

def main():
    p = argparse.ArgumentParser()
    p.add_argument("data_root")
    p.add_argument("--output", default="cards.csv")
    p.add_argument("--pack", help="filter by pack code")
    p.add_argument("--type", help="filter by type_code (hero, ally, event, ...)")
    args = p.parse_args()

    cards = load_cards(args.data_root, pack_code=args.pack, type_code=args.type)
    print(f"Loaded {len(cards)} cards", file=sys.stderr)

    with open(args.output, "w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=FIELDS, extrasaction="ignore")
        w.writeheader()
        w.writerows(cards)
    print(f"wrote {args.output} ({len(cards)} rows)")

if __name__ == "__main__":
    main()
