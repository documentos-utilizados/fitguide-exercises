import os
import json
import glob
import sys
import argparse

try:
    from translator import title_case
except ImportError:
    try:
        from scripts.translator import title_case
    except ImportError:
        def title_case(s):
            return s.capitalize() if s else s

try:
    from normalize_primary_muscle_categories import normalize_primary_muscles
except ImportError:
    try:
        from scripts.normalize_primary_muscle_categories import normalize_primary_muscles
    except ImportError:
        def normalize_primary_muscles(p):
            return []

def normalize_muscle_name(name):
    if not name:
        return ""
    name = name.strip()
    if name.lower() == "costas (meio)" or name.lower() == "costas_(meio)":
        return "Costas (meio)"
    return title_case(name)

def main():
    parser = argparse.ArgumentParser(description="Extract unique primary muscles categories as a plain list")
    parser.add_argument("--base-dir", default="./database", help="Database directory")
    args = parser.parse_args()

    base_dir = args.base_dir
    workouts_dir = os.path.join(base_dir, "workouts")
    exercises_dir = os.path.join(base_dir, "exercises")
    output_dir = os.path.join(base_dir, "metadata", "categories", "primary_muscle")

    os.makedirs(output_dir, exist_ok=True)
    output_pt = os.path.join(output_dir, "pt.json")

    primary_muscles = set()

    exercise_files = glob.glob(os.path.join(exercises_dir, "*", "pt.json"))
    exercises_by_id = {}
    for ef in exercise_files:
        try:
            with open(ef, "r", encoding="utf-8") as f:
                data = json.load(f)
                ex_id = data.get("id")
                if ex_id:
                    exercises_by_id[ex_id] = data
                for m in data.get("primaryMuscles", []):
                    norm = normalize_muscle_name(m)
                    if norm:
                        primary_muscles.add(norm)
        except Exception:
            continue

    workout_files = glob.glob(os.path.join(workouts_dir, "*", "*.json"))
    for wf in workout_files:
        try:
            with open(wf, "r", encoding="utf-8") as f:
                wdata = json.load(f)
                for ex_item in wdata.get("exercises", []):
                    ex_id = ex_item.get("exercise_id")
                    if ex_id and ex_id in exercises_by_id:
                        for m in exercises_by_id[ex_id].get("primaryMuscles", []):
                            norm = normalize_muscle_name(m)
                            if norm:
                                primary_muscles.add(norm)
                    for m in ex_item.get("primaryMuscles", []):
                        norm = normalize_muscle_name(m)
                        if norm:
                            primary_muscles.add(norm)
        except Exception:
            continue

    sorted_muscles = sorted(list(primary_muscles), key=lambda s: s.lower())

    with open(output_pt, "w", encoding="utf-8") as f:
        json.dump(sorted_muscles, f, ensure_ascii=False, indent=2)
        f.write("\n")

    normalized = normalize_primary_muscles(output_pt)
    final_count = len(normalized) if normalized else len(sorted_muscles)

    print(f"Successfully extracted and normalized {final_count} primary muscle categories to {output_pt}")

if __name__ == "__main__":
    main()
