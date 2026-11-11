import os
import re
import json
import argparse
from deep_translator import MyMemoryTranslator, GoogleTranslator

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CACHE_FILE = os.path.join(BASE_DIR, "database", "metadata", "translations_cache.json")

def load_cache():
    if os.path.exists(CACHE_FILE):
        try:
            with open(CACHE_FILE, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception:
            return {}
    return {}

def save_cache(cache):
    os.makedirs(os.path.dirname(CACHE_FILE), exist_ok=True)
    with open(CACHE_FILE, "w", encoding="utf-8") as f:
        json.dump(cache, f, ensure_ascii=False, indent=2)
        f.write("\n")

def clean_text(text):
    if not text:
        return ""
    text = re.sub(r"<[^>]+>", "", text)
    text = re.sub(r"\s+", " ", text).strip()
    return text

def translate_sentence(text, source="pt", target="en"):
    if not text or not text.strip():
        return text

    clean_source = clean_text(text)
    cache = load_cache()
    cache_key = f"{source}:{target}:{clean_source}"

    if cache_key in cache:
        return cache[cache_key]

    placeholders = re.findall(r"%[0-9]+\$[a-zA-Z]|%[a-zA-Z]", clean_source)
    temp_text = clean_source
    for i, p in enumerate(placeholders):
        temp_text = temp_text.replace(p, f"__PH{i}__")

    translated = None
    try:
        src_lang = "pt-BR" if source == "pt" else "en-GB"
        tgt_lang = "en-GB" if target == "en" else "pt-BR"
        raw = MyMemoryTranslator(source=src_lang, target=tgt_lang).translate(temp_text)
        if raw and raw.lower() != temp_text.lower():
            translated = raw.strip()
    except Exception:
        pass

    if not translated:
        try:
            raw = GoogleTranslator(source=source, target=target).translate(temp_text)
            if raw:
                translated = raw.strip()
        except Exception:
            pass

    if not translated:
        translated = temp_text

    for i, p in enumerate(placeholders):
        translated = translated.replace(f"__PH{i}__", p)
        translated = translated.replace(f"__ph{i}__", p)

    cache[cache_key] = translated
    save_cache(cache)
    return translated

def translate_group(base_dir, group_name):
    pt_path = os.path.join(base_dir, "database", "notifications", "ptbr", f"{group_name}.json")
    en_path = os.path.join(base_dir, "database", "notifications", "en", f"{group_name}.json")

    if not os.path.exists(pt_path):
        return

    with open(pt_path, "r", encoding="utf-8") as f:
        pt_doc = json.load(f)

    en_doc = {"group": pt_doc.get("group", group_name), "name": pt_doc.get("name", ""), "items": []}
    if os.path.exists(en_path):
        try:
            with open(en_path, "r", encoding="utf-8") as f:
                en_doc = json.load(f)
        except Exception:
            pass

    en_name = en_doc.get("name")
    if not en_name or en_name == pt_doc.get("name"):
        en_doc["name"] = translate_sentence(pt_doc.get("name", ""), "pt", "en")

    en_items_map = {item["id"]: item for item in en_doc.get("items", [])}

    updated_items = []
    for pt_item in pt_doc.get("items", []):
        item_id = pt_item["id"]
        pt_title = pt_item.get("title", "")
        pt_messages = pt_item.get("messages", [])

        if item_id in en_items_map:
            en_item = en_items_map[item_id]
            if not en_item.get("title"):
                en_item["title"] = translate_sentence(pt_title, "pt", "en")

            en_msgs = list(en_item.get("messages", []))
            if len(en_msgs) < len(pt_messages):
                for idx in range(len(en_msgs), len(pt_messages)):
                    tr = translate_sentence(pt_messages[idx], "pt", "en")
                    en_msgs.append(tr)
            en_item["messages"] = en_msgs
            updated_items.append(en_item)
        else:
            en_title = translate_sentence(pt_title, "pt", "en")
            en_msgs = [translate_sentence(msg, "pt", "en") for msg in pt_messages]
            updated_items.append({
                "id": item_id,
                "title": en_title,
                "messages": en_msgs
            })

    en_doc["items"] = updated_items
    os.makedirs(os.path.dirname(en_path), exist_ok=True)
    with open(en_path, "w", encoding="utf-8") as f:
        json.dump(en_doc, f, ensure_ascii=False, indent=2)
        f.write("\n")

    config_path = os.path.join(base_dir, "database", "notifications", "config_notifications.json")
    from datetime import datetime
    now_iso = datetime.now().astimezone().isoformat()
    with open(config_path, "w", encoding="utf-8") as f:
        json.dump({"last_modified": now_iso, "last_modified_group": group_name}, f, ensure_ascii=False, indent=2)
        f.write("\n")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-dir", default=BASE_DIR)
    parser.add_argument("--group", default="")
    parser.add_argument("--text", default="")
    args = parser.parse_args()

    if args.text:
        print(translate_sentence(args.text, "pt", "en"), end="")
        return

    pt_dir = os.path.join(args.base_dir, "database", "notifications", "ptbr")
    if not os.path.exists(pt_dir):
        return

    if args.group:
        translate_group(args.base_dir, args.group)
    else:
        for f in os.listdir(pt_dir):
            if f.endswith(".json"):
                grp = f[:-5]
                translate_group(args.base_dir, grp)
    print("SUCCESS")

if __name__ == "__main__":
    main()
