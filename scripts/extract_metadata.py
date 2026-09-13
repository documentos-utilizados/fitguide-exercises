import os
import json
import glob
import sys
import argparse

CATEGORY_TRANSLATIONS_FALLBACK = {
    "pt": {
        "strength": "Força",
        "stretching": "Alongamento",
        "plyometrics": "Pliometria",
        "powerlifting": "Powerlifting",
        "strongman": "Strongman",
        "cardio": "Cardio",
        "olympic weightlifting": "Levantamento de Peso Olímpico"
    },
    "en": {
        "strength": "Strength",
        "stretching": "Stretching",
        "plyometrics": "Plyometrics",
        "powerlifting": "Powerlifting",
        "strongman": "Strongman",
        "cardio": "Cardio",
        "olympic weightlifting": "Olympic Weightlifting"
    }
}

def title_case(s):
    if not s:
        return s
    words = s.split(" ")
    return " ".join(w.capitalize() if w not in ["de", "do", "da", "dos", "das", "e", "em", "para", "com", "of", "and", "in", "to", "with", "a", "an", "the"] else w for w in words)

def main():
    parser = argparse.ArgumentParser(description="Extract metadata from exercise datasets")
    parser.add_argument("--dist-dir", default="./database/dist", help="Directory containing compiled exercises JSONs")
    parser.add_argument("--base-dir", default="./database/metadata", help="Base directory where topic folders will be created")
    args = parser.parse_args()

    dist_dir = args.dist_dir
    base_dir = args.base_dir

    en_file = os.path.join(dist_dir, "exercises_en.json")
    if not os.path.exists(en_file):
        print(f"Error: Base English dataset not found at {en_file}")
        sys.exit(1)

    with open(en_file, "r", encoding="utf-8") as f:
        en_exercises = json.load(f)

    en_by_id = {ex["id"]: ex for ex in en_exercises if "id" in ex}

    lang_files = glob.glob(os.path.join(dist_dir, "exercises_*.json"))
    languages = []
    for lf in lang_files:
        basename = os.path.basename(lf)
        lang_code = basename.replace("exercises_", "").replace(".json", "")
        languages.append((lang_code, lf))

    if not languages:
        print(f"No exercise files found in {dist_dir}")
        sys.exit(1)

    for topic in ["categories", "levels", "mechanics", "equipments", "forces", "muscles"]:
        os.makedirs(os.path.join(base_dir, topic), exist_ok=True)

    for lang, filepath in languages:
        with open(filepath, "r", encoding="utf-8") as f:
            lang_exercises = json.load(f)

        categories = {}
        levels = {}
        mechanics = {}
        equipments = {}
        forces = {}
        muscles = {}

        for lang_ex in lang_exercises:
            ex_id = lang_ex.get("id")
            if not ex_id or ex_id not in en_by_id:
                continue
            en_ex = en_by_id[ex_id]

            cat_en = en_ex.get("category")
            cat_lang = lang_ex.get("category")
            if cat_en:
                if cat_lang:
                    val = cat_lang
                else:
                    val = CATEGORY_TRANSLATIONS_FALLBACK.get(lang, {}).get(cat_en, cat_en)
                categories[cat_en] = title_case(val)

            lvl_en = en_ex.get("level")
            lvl_lang = lang_ex.get("level")
            if lvl_en:
                val = lvl_lang if lvl_lang else lvl_en
                levels[lvl_en] = title_case(val)

            mech_en = en_ex.get("mechanic")
            mech_lang = lang_ex.get("mechanic")
            if mech_en:
                val = mech_lang if mech_lang else mech_en
                mechanics[mech_en] = title_case(val)

            eq_en = en_ex.get("equipment")
            eq_lang = lang_ex.get("equipment")
            if eq_en:
                val = eq_lang if eq_lang else eq_en
                equipments[eq_en] = title_case(val)

            force_en = en_ex.get("force")
            force_lang = lang_ex.get("force")
            if force_en:
                val = force_lang if force_lang else force_en
                forces[force_en] = title_case(val)

            en_prim = en_ex.get("primaryMuscles", [])
            lang_prim = lang_ex.get("primaryMuscles", [])
            for idx, m_en in enumerate(en_prim):
                if m_en:
                    m_lang = lang_prim[idx] if idx < len(lang_prim) else m_en
                    muscles[m_en] = title_case(m_lang)

            en_sec = en_ex.get("secondaryMuscles", [])
            lang_sec = lang_ex.get("secondaryMuscles", [])
            for idx, m_en in enumerate(en_sec):
                if m_en:
                    m_lang = lang_sec[idx] if idx < len(lang_sec) else m_en
                    muscles[m_en] = title_case(m_lang)

        topic_data = {
            "categories": dict(sorted(categories.items())),
            "levels": dict(sorted(levels.items())),
            "mechanics": dict(sorted(mechanics.items())),
            "equipments": dict(sorted(equipments.items())),
            "forces": dict(sorted(forces.items())),
            "muscles": dict(sorted(muscles.items()))
        }

        for topic_folder, data in topic_data.items():
            out_file = os.path.join(base_dir, topic_folder, f"{lang}.json")
            with open(out_file, "w", encoding="utf-8") as out:
                json.dump(data, out, ensure_ascii=False, indent=2)
                out.write("\n")

        print(f"Extracted metadata for language '{lang}' into dedicated topic directories")

if __name__ == "__main__":
    main()
