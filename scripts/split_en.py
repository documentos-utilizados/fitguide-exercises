import os
import json
import sys

def main():
    src_file = "./database/dist/exercises_en.json"
    if len(sys.argv) > 1:
        src_file = sys.argv[1]
    
    if not os.path.exists(src_file):
        print(f"Error: file not found at {src_file}")
        sys.exit(1)
        
    dest_dir = "./database/exercises"
    os.makedirs(dest_dir, exist_ok=True)
    
    with open(src_file, "r", encoding="utf-8") as f:
        exercises = json.load(f)
        
    for ex in exercises:
        ex_id = ex["id"]
        ex_folder = os.path.join(dest_dir, ex_id)
        os.makedirs(ex_folder, exist_ok=True)
        filepath = os.path.join(ex_folder, "en.json")
        with open(filepath, "w", encoding="utf-8") as out:
            json.dump(ex, out, ensure_ascii=False, indent=2)
            out.write("\n")
            
    print(f"Successfully extracted {len(exercises)} files to {dest_dir}")

if __name__ == "__main__":
    main()
