import os
import json
import unicodedata
import argparse

CANONICAL_MUSCLE_MAP = {
    "abdominais": "ABDOMINAIS",
    "abdominal": "ABDOMINAIS",
    "abdomen": "ABDOMINAIS",
    "abdômen": "ABDOMINAIS",
    "abdutores": "ABDUTORES",
    "abdutor": "ABDUTORES",
    "adutores": "ADUTORES",
    "adutor": "ADUTORES",
    "antebraços": "ANTEBRAÇOS",
    "antebraço": "ANTEBRAÇOS",
    "antebracos": "ANTEBRAÇOS",
    "antebraco": "ANTEBRAÇOS",
    "bíceps": "BÍCEPS",
    "biceps": "BÍCEPS",
    "costas (meio)": "COSTAS (MEIO)",
    "costas_(meio)": "COSTAS (MEIO)",
    "costas meio": "COSTAS (MEIO)",
    "meio das costas": "COSTAS (MEIO)",
    "dorsais": "DORSAIS",
    "dorsal": "DORSAIS",
    "latissimo": "DORSAIS",
    "latíssimo": "DORSAIS",
    "glúteos": "GLÚTEOS",
    "gluteos": "GLÚTEOS",
    "glúteo": "GLÚTEOS",
    "gluteo": "GLÚTEOS",
    "isquiotibiais": "ISQUIOTIBIAIS",
    "isquiotibial": "ISQUIOTIBIAIS",
    "posterior de coxa": "ISQUIOTIBIAIS",
    "posteriordecoxa": "ISQUIOTIBIAIS",
    "lombar": "LOMBAR",
    "lombares": "LOMBAR",
    "ombros": "OMBROS",
    "ombro": "OMBROS",
    "deltoides": "OMBROS",
    "deltóides": "OMBROS",
    "deltoide": "OMBROS",
    "deltóide": "OMBROS",
    "panturrilhas": "PANTURRILHAS",
    "panturrilha": "PANTURRILHAS",
    "gemeos": "PANTURRILHAS",
    "gêmeos": "PANTURRILHAS",
    "peitoral": "PEITORAL",
    "peito": "PEITORAL",
    "pescoço": "PESCOÇO",
    "pescoco": "PESCOÇO",
    "quadríceps": "QUADRÍCEPS",
    "quadriceps": "QUADRÍCEPS",
    "trapézio": "TRAPÉZIO",
    "trapezio": "TRAPÉZIO",
    "tríceps": "TRÍCEPS",
    "triceps": "TRÍCEPS"
}

def strip_accents(text):
    return "".join(c for c in unicodedata.normalize("NFD", text) if unicodedata.category(c) != "Mn")

def normalize_single_muscle(term):
    if not term:
        return ""
    clean = term.strip().lower()
    if clean in CANONICAL_MUSCLE_MAP:
        return CANONICAL_MUSCLE_MAP[clean]
    clean_no_accents = strip_accents(clean)
    if clean_no_accents in CANONICAL_MUSCLE_MAP:
        return CANONICAL_MUSCLE_MAP[clean_no_accents]
    return term.strip().upper()

def normalize_primary_muscles(file_path):
    if not os.path.exists(file_path):
        return []

    with open(file_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    if not isinstance(data, list):
        return []

    normalized_set = set()
    for item in data:
        if isinstance(item, str):
            norm = normalize_single_muscle(item)
            if norm:
                normalized_set.add(norm)

    sorted_list = sorted(list(normalized_set), key=lambda s: strip_accents(s).lower())

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(sorted_list, f, ensure_ascii=False, indent=2)
        f.write("\n")

    return sorted_list

def main():
    parser = argparse.ArgumentParser(description="Normalize primary muscle categories to uppercase and deduplicate")
    parser.add_argument("--base-dir", default="./database", help="Database directory")
    args = parser.parse_args()

    pt_file = os.path.join(args.base_dir, "metadata", "categories", "primary_muscle", "pt.json")
    results = normalize_primary_muscles(pt_file)
    print(f"Successfully normalized {len(results)} primary muscle categories in {pt_file}")

if __name__ == "__main__":
    main()
