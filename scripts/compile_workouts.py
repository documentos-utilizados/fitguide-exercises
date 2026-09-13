import os
import json
import glob
import sys
import argparse

def main():
    parser = argparse.ArgumentParser(description="Compile and hydrate workouts for multiple languages")
    parser.add_argument("--dist-dir", default="./database/dist", help="Directory containing compiled exercises JSONs")
    parser.add_argument("--workouts-dir", default="./database/workouts", help="Directory containing workout definitions")
    args = parser.parse_args()

    dist_dir = args.dist_dir
    workouts_dir = args.workouts_dir

    seed_files = sorted(glob.glob(os.path.join(workouts_dir, "seeds", "*.json")))
    template_files = sorted(glob.glob(os.path.join(workouts_dir, "templates", "*.json")))
    all_workout_files = seed_files + template_files

    if not all_workout_files:
        print(f"Error: No workout files found in {workouts_dir}")
        sys.exit(1)

    raw_workouts = []
    for wf in all_workout_files:
        with open(wf, "r", encoding="utf-8") as f:
            raw_workouts.append(json.load(f))

    lang_files = glob.glob(os.path.join(dist_dir, "exercises_*.json"))
    if not lang_files:
        print(f"Error: No exercises datasets found in {dist_dir}")
        sys.exit(1)

    for lf in lang_files:
        basename = os.path.basename(lf)
        lang_code = basename.replace("exercises_", "").replace(".json", "")

        with open(lf, "r", encoding="utf-8") as f:
            exercises_list = json.load(f)

        exercises_by_id = {ex["id"]: ex for ex in exercises_list if "id" in ex}

        hydrated_workouts = []
        errors = 0

        for workout in raw_workouts:
            w_copy = dict(workout)
            w_copy["language"] = lang_code
            hydrated_exercises = []

            for ex_config in workout.get("exercises", []):
                ex_id = ex_config.get("exercise_id")
                if not ex_id or ex_id not in exercises_by_id:
                    print(f"Error [{lang_code}]: Exercise ID '{ex_id}' in workout '{workout.get('id')}' not found in dataset!")
                    errors += 1
                    continue

                full_ex = exercises_by_id[ex_id]

                hydrated_ex = {
                    "exercise_id": ex_id,
                    "name": full_ex.get("name", ""),
                    "category": full_ex.get("category", ""),
                    "level": full_ex.get("level", ""),
                    "primaryMuscles": full_ex.get("primaryMuscles", []),
                    "secondaryMuscles": full_ex.get("secondaryMuscles", []),
                    "equipment": full_ex.get("equipment"),
                    "mechanic": full_ex.get("mechanic"),
                    "force": full_ex.get("force"),
                    "instructions": full_ex.get("instructions", []),
                    "images": full_ex.get("images", []),
                    "sets": ex_config.get("sets", 3),
                    "reps": ex_config.get("reps", "10-12"),
                    "weight_kg": ex_config.get("weight_kg", 0.0),
                    "rest_time_seconds": ex_config.get("rest_time_seconds", 60),
                    "is_time_based": ex_config.get("is_time_based", False),
                    "time_seconds": ex_config.get("time_seconds", 0)
                }
                hydrated_exercises.append(hydrated_ex)

            w_copy["exercises"] = hydrated_exercises
            hydrated_workouts.append(w_copy)

        if errors > 0:
            print(f"Encountered {errors} missing exercise references for language '{lang_code}'.")
            sys.exit(1)

        out_path = os.path.join(dist_dir, f"workouts_{lang_code}.json")
        os.makedirs(os.path.dirname(out_path), exist_ok=True)
        with open(out_path, "w", encoding="utf-8") as out:
            json.dump(hydrated_workouts, out, ensure_ascii=False, indent=2)
            out.write("\n")

        print(f"Successfully compiled {len(hydrated_workouts)} hydrated workouts into {out_path}")

if __name__ == "__main__":
    main()
