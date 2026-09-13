import os
import json
import glob
import sys
import argparse

try:
    from translator import translate_term, title_case
except ImportError:
    from scripts.translator import translate_term, title_case

def sanitize_key(s):
    if not s:
        return s
    return s.lower().strip().replace(" ", "_").replace("-", "_")

def main():
    parser = argparse.ArgumentParser(description="Extract metadata from exercise datasets with dynamic translation")
    parser.add_argument("--dist-dir", default="./database/dist", help="Directory containing compiled exercises JSONs")
    parser.add_argument("--base-dir", default="./database/metadata", help="Base directory where topic folders will be created")
    args = parser.parse_args()

    dist_dir = args.dist_dir
    base_dir = args.base_dir

    pt_file = os.path.join(dist_dir, "exercises_pt.json")
    en_file = os.path.join(dist_dir, "exercises_en.json")

    if not os.path.exists(pt_file) and not os.path.exists(en_file):
        print(f"Error: Neither {pt_file} nor {en_file} found in {dist_dir}")
        sys.exit(1)

    pt_exercises = []
    if os.path.exists(pt_file):
        with open(pt_file, "r", encoding="utf-8") as f:
            pt_exercises = json.load(f)

    en_exercises = []
    if os.path.exists(en_file):
        with open(en_file, "r", encoding="utf-8") as f:
            en_exercises = json.load(f)

    pt_by_id = {ex["id"]: ex for ex in pt_exercises if "id" in ex}
    en_by_id = {ex["id"]: ex for ex in en_exercises if "id" in ex}
    all_exercise_ids = set(pt_by_id.keys()).union(set(en_by_id.keys()))

    topics = ["categories", "levels", "mechanics", "equipments", "forces", "muscles"]
    for topic in topics:
        os.makedirs(os.path.join(base_dir, topic), exist_ok=True)

    metadata_pt = {t: {} for t in topics}
    metadata_en = {t: {} for t in topics}

    raw_items = {
        "categories": set(),
        "levels": set(),
        "mechanics": set(),
        "equipments": set(),
        "forces": set(),
        "muscles": set()
    }

    pairs = {
        "categories": {},
        "levels": {},
        "mechanics": {},
        "equipments": {},
        "forces": {},
        "muscles": {}
    }

    for topic in topics:
        pt_path = os.path.join(base_dir, topic, "pt.json")
        en_path = os.path.join(base_dir, topic, "en.json")
        existing_pt = {}
        existing_en = {}
        if os.path.exists(pt_path):
            try:
                with open(pt_path, "r", encoding="utf-8") as f:
                    existing_pt = json.load(f)
            except Exception:
                pass
        if os.path.exists(en_path):
            try:
                with open(en_path, "r", encoding="utf-8") as f:
                    existing_en = json.load(f)
            except Exception:
                pass
        all_keys = set(existing_pt.keys()).union(set(existing_en.keys()))
        for k in all_keys:
            raw_items[topic].add(k)
            pairs[topic][k] = (existing_pt.get(k), existing_en.get(k))

    for ex_id in all_exercise_ids:
        ex_pt = pt_by_id.get(ex_id)
        ex_en = en_by_id.get(ex_id)

        for topic, field in [("categories", "category"), ("levels", "level"), ("mechanics", "mechanic"), ("equipments", "equipment"), ("forces", "force")]:
            v_pt = ex_pt.get(field) if ex_pt else None
            v_en = ex_en.get(field) if ex_en else None
            if v_pt or v_en:
                k = sanitize_key(v_pt or v_en)
                raw_items[topic].add(k)
                if v_pt and v_en and v_pt.strip().lower() != v_en.strip().lower():
                    pairs[topic][k] = (v_pt, v_en)
                elif v_pt and k not in pairs[topic]:
                    pairs[topic][k] = (v_pt, None)
                elif v_en and k not in pairs[topic]:
                    pairs[topic][k] = (None, v_en)

        m_pt = (ex_pt.get("primaryMuscles", []) if ex_pt else []) + (ex_pt.get("secondaryMuscles", []) if ex_pt else [])
        m_en = (ex_en.get("primaryMuscles", []) if ex_en else []) + (ex_en.get("secondaryMuscles", []) if ex_en else [])
        for idx, m_val in enumerate(m_pt):
            if m_val:
                k = sanitize_key(m_val)
                raw_items["muscles"].add(k)
                en_match = m_en[idx] if idx < len(m_en) else None
                if k not in pairs["muscles"] or (pairs["muscles"][k][1] is None and en_match):
                    pairs["muscles"][k] = (m_val, en_match)
        for idx, m_val in enumerate(m_en):
            if m_val:
                k = sanitize_key(m_val)
                raw_items["muscles"].add(k)
                if k not in pairs["muscles"]:
                    pairs["muscles"][k] = (None, m_val)

    for topic in topics:
        for k in sorted(raw_items[topic]):
            val_pt_src, val_en_src = pairs[topic].get(k, (None, None))
            if val_pt_src and val_en_src and val_pt_src.strip().lower() != val_en_src.strip().lower():
                metadata_pt[topic][k] = title_case(val_pt_src)
                metadata_en[topic][k] = title_case(val_en_src)
            elif val_pt_src:
                lbl_en = translate_term(val_pt_src, source="pt", target="en")
                lbl_pt = translate_term(val_pt_src, source="en", target="pt")
                metadata_pt[topic][k] = lbl_pt
                metadata_en[topic][k] = lbl_en
            elif val_en_src:
                lbl_en = title_case(val_en_src)
                lbl_pt = translate_term(val_en_src, source="en", target="pt")
                metadata_pt[topic][k] = lbl_pt
                metadata_en[topic][k] = lbl_en

    for topic in topics:
        out_pt = os.path.join(base_dir, topic, "pt.json")
        out_en = os.path.join(base_dir, topic, "en.json")

        with open(out_pt, "w", encoding="utf-8") as f_pt:
            json.dump(dict(sorted(metadata_pt[topic].items())), f_pt, ensure_ascii=False, indent=2)
            f_pt.write("\n")

        with open(out_en, "w", encoding="utf-8") as f_en:
            json.dump(dict(sorted(metadata_en[topic].items())), f_en, ensure_ascii=False, indent=2)
            f_en.write("\n")

    print("Metadata extraction completed with dynamic translation synchronization.")

if __name__ == "__main__":
    main()
