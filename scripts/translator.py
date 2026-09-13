import os
import re
import json

CACHE_FILE_PATH = os.path.join(os.path.dirname(__file__), "..", "database", "metadata", "translations_cache.json")

def load_cache():
    if os.path.exists(CACHE_FILE_PATH):
        try:
            with open(CACHE_FILE_PATH, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception:
            return {}
    return {}

def save_cache(cache):
    os.makedirs(os.path.dirname(CACHE_FILE_PATH), exist_ok=True)
    with open(CACHE_FILE_PATH, "w", encoding="utf-8") as f:
        json.dump(cache, f, ensure_ascii=False, indent=2)
        f.write("\n")

def clean_translated_text(text):
    if not text:
        return text
    text = re.sub(r"<[^>]+>", "", text)
    text = re.sub(r"^\((.+)\)$", r"\1", text)
    text = re.sub(r"\s+", " ", text).strip()
    text = text.rstrip(",;.")
    return text

def title_case(s):
    if not s:
        return s
    lower_words = {"de", "do", "da", "dos", "das", "e", "em", "para", "com", "of", "and", "in", "to", "with", "a", "an", "the"}
    words = s.split(" ")
    return " ".join(w.capitalize() if w.lower() not in lower_words else w.lower() for w in words)

def translate_term(text, source="pt", target="en"):
    if not text or not text.strip():
        return text

    clean_source = text.strip()
    cache_key = f"{source}:{target}:{clean_source.lower()}"
    cache = load_cache()

    if cache_key in cache:
        return cache[cache_key]

    translated = None

    try:
        from deep_translator import MyMemoryTranslator
        src_lang = "pt-BR" if source == "pt" else "en-GB"
        tgt_lang = "en-GB" if target == "en" else "pt-BR"
        raw = MyMemoryTranslator(source=src_lang, target=tgt_lang).translate(clean_source)
        cleaned = clean_translated_text(raw)
        if cleaned and cleaned.lower() != clean_source.lower():
            translated = title_case(cleaned)
    except Exception:
        pass

    if not translated:
        try:
            from deep_translator import GoogleTranslator
            raw = GoogleTranslator(source=source, target=target).translate(clean_source)
            cleaned = clean_translated_text(raw)
            if cleaned:
                translated = title_case(cleaned)
        except Exception:
            pass

    if not translated:
        translated = title_case(clean_source)

    cache[cache_key] = translated
    save_cache(cache)
    return translated
