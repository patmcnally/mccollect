#!/usr/bin/env python3
"""Quick script to dump Marvel Champions card data to CSV."""
import json, csv, sys, argparse
from pathlib import Path

def load_cards(data_root):
    pack_dir = Path(data_root) / "pack"
    cards = []
    for f in sorted(pack_dir.glob("*.json")):
        with open(f) as fp:
            cards.extend(json.load(fp))
    return cards

def main():
    p = argparse.ArgumentParser()
    p.add_argument("data_root", help="path to marvelsdb-json-data clone")
    p.add_argument("--output", default="cards.csv")
    args = p.parse_args()

    cards = load_cards(args.data_root)
    print(f"Loaded {len(cards)} cards", file=sys.stderr)

    with open(args.output, "w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=["code","name","type_code","pack_code","cost"])
        w.writeheader()
        for c in cards:
            w.writerow({k: c.get(k,"") for k in ["code","name","type_code","pack_code","cost"]})
    print(f"wrote {args.output}")

if __name__ == "__main__":
    main()
