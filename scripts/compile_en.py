import os
import json
import glob
import sys

def main():
    dest_dir = "./database/exercises"
    schema_path = "./database/schemas/exercise.schema.json"
    dist_path = "./database/dist/exercises_en.json"
    
    if not os.path.exists(dest_dir):
        print(f"Directory {dest_dir} does not exist.")
        sys.exit(1)
        
    json_files = glob.glob(os.path.join(dest_dir, "*", "en.json"))
    print(f"Found {len(json_files)} exercise files in {dest_dir}.")
    
    compiled_exercises = []
    errors = 0
    
    for filepath in json_files:
        filename = os.path.basename(filepath)
        with open(filepath, "r", encoding="utf-8") as f:
            try:
                data = json.load(f)
                required_keys = ["id", "name", "instructions"]
                for key in required_keys:
                    if key not in data:
                        print(f"Validation error in {filename}: Missing key '{key}'")
                        errors += 1
                        
                compiled_exercises.append(data)
            except Exception as e:
                print(f"Error parsing JSON in {filename}: {e}")
                errors += 1
                
    if errors > 0:
        print(f"Validation failed with {errors} errors.")
        sys.exit(1)
        
    compiled_exercises.sort(key=lambda x: x.get("id", ""))
    
    os.makedirs(os.path.dirname(dist_path), exist_ok=True)
    with open(dist_path, "w", encoding="utf-8") as f:
        json.dump(compiled_exercises, f, ensure_ascii=False, indent=2)
        
    print(f"Successfully compiled {len(compiled_exercises)} exercises into {dist_path}")

if __name__ == "__main__":
    main()
