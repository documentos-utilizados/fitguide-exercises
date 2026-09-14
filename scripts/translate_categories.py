import os
import json
import sys
import argparse

try:
    from translator import translate_term, title_case
except ImportError:
    from scripts.translator import translate_term, title_case

MUSCLE_TRANSLATIONS_MAP = {
    "abdominais": "ABDOMINALS",
    "abdutores": "ABDUCTORS",
    "adutores": "ADDUCTORS",
    "antebraços": "FOREARMS",
    "bíceps": "BICEPS",
    "costas (meio)": "MIDDLE BACK",
    "dorsais": "LATS",
    "glúteos": "GLUTES",
    "isquiotibiais": "HAMSTRINGS",
    "lombar": "LOWER BACK",
    "ombros": "SHOULDERS",
    "panturrilhas": "CALVES",
    "peitoral": "CHEST",
    "pescoço": "NECK",
    "quadríceps": "QUADRICEPS",
    "trapézio": "TRAPS",
    "tríceps": "TRICEPS"
}

def translate_muscle(term, source="pt", target="en"):
    if not term:
        return term
    clean = term.strip().lower()
    if target == "en" and clean in MUSCLE_TRANSLATIONS_MAP:
        res = MUSCLE_TRANSLATIONS_MAP[clean]
        return res if term.isupper() else title_case(res)
    translated = translate_term(term, source=source, target=target)
    return translated.upper() if term.isupper() else translated

def main():
    parser = argparse.ArgumentParser(description="Translate category JSON between languages")
    parser.add_argument("--base-dir", default="./database", help="Database directory")
    parser.add_argument("--category-type", default="primary_muscle", help="Category type folder")
    parser.add_argument("--source", default="pt", help="Source language code")
    parser.add_argument("--target", default="en", help="Target language code")
    args = parser.parse_args()

    cat_dir = os.path.join(args.base_dir, "metadata", "categories", args.category_type)
    src_file = os.path.join(cat_dir, f"{args.source}.json")
    tgt_file = os.path.join(cat_dir, f"{args.target}.json")

    if not os.path.exists(src_file):
        print(f"Error: Source file {src_file} does not exist.")
        sys.exit(1)

    with open(src_file, "r", encoding="utf-8") as f:
        src_data = json.load(f)

    if isinstance(src_data, list):
        translated_data = []
        for item in src_data:
            if args.category_type == "primary_muscle":
                translated = translate_muscle(item, source=args.source, target=args.target)
            else:
                translated = translate_term(item, source=args.source, target=args.target)
            translated_data.append(translated)
    elif isinstance(src_data, dict):
        translated_data = {}
        for key, val in src_data.items():
            translated = translate_term(val, source=args.source, target=args.target)
            translated_data[key] = translated
    else:
        print("Error: Unsupported JSON format")
        sys.exit(1)

    with open(tgt_file, "w", encoding="utf-8") as f:
        json.dump(translated_data, f, ensure_ascii=False, indent=2)
        f.write("\n")

    print(f"Successfully translated {args.category_type} from {args.source} to {args.target} at {tgt_file}")

if __name__ == "__main__":
    main()
