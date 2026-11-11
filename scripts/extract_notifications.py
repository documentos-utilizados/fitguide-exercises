import os
import json
import xml.etree.ElementTree as ET

def clean(text):
    if not text:
        return ""
    return text.replace(r"\'", "'").replace("%%", "%").strip()

def parse(xml_path):
    tree = ET.parse(xml_path)
    root = tree.getroot()
    strings = {s.attrib["name"]: clean(s.text) for s in root.findall("string")}
    arrays = {a.attrib["name"]: [clean(item.text) for item in a.findall("item")] for a in root.findall("string-array")}
    return strings, arrays

groups_def = [
    {
        "file": "strings_notif_workout_call.xml",
        "group_id": "workout_call",
        "name_pt": "Chamada ao Treino",
        "name_en": "Workout Call",
        "items": [
            ("WORKOUT_SCHEDULED_TODAY", "notif_workout_scheduled_title", "notif_workout_scheduled_messages"),
            ("WORKOUT_PREFERRED_TIME", "notif_workout_preferred_time_title", "notif_workout_preferred_time_messages"),
            ("WORKOUT_QUICK_CALL", "notif_workout_quick_call_title", "notif_workout_quick_call_messages"),
            ("WORKOUT_DAY_ONE", "notif_workout_day_one_title", "notif_workout_day_one_messages"),
            ("WORKOUT_OPENED_NO_WORKOUT", "notif_workout_opened_no_workout_title", "notif_workout_opened_no_workout_messages"),
            ("POST_WORKOUT_CONGRATS", "notif_post_workout_congrats_title", "notif_post_workout_congrats_messages"),
            ("POST_WORKOUT_ALL_COMPLETED", "notif_post_workout_congrats_title", "notif_post_workout_all_completed_messages"),
            ("POST_WORKOUT_PARTIAL_COMPLETED", "notif_post_workout_congrats_title", "notif_post_workout_partial_completed_messages"),
            ("WORKOUT_ABANDONED", "notif_workout_abandoned_title", "notif_workout_abandoned_messages"),
            ("GYM_WIFI_ARRIVAL", "notif_gym_wifi_title", "notif_gym_wifi_messages"),
            ("ONE_HOUR_BEFORE_HABITUAL", "notif_one_hour_before_title", "notif_one_hour_before_messages"),
            ("MISSED_MORNING_AFTERNOON", "notif_missed_morning_afternoon_title", "notif_missed_morning_afternoon_body"),
            ("MISSED_DAY_NIGHT", "notif_missed_day_night_title", "notif_missed_day_night_body"),
        ]
    },
    {
        "file": "strings_notif_streak_consistency.xml",
        "group_id": "streak_consistency",
        "name_pt": "Ofensiva e Consistência",
        "name_en": "Streak and Consistency",
        "items": [
            ("STREAK_ACTIVE", "notif_streak_active_title", "notif_streak_active_messages"),
            ("STREAK_AT_RISK", "notif_streak_at_risk_title", "notif_streak_at_risk_messages"),
            ("STREAK_NEW_RECORD", "notif_streak_new_record_title", "notif_streak_new_record_messages"),
            ("WEEKLY_GOAL_COMPLETED", "notif_weekly_goal_title", "notif_weekly_goal_messages"),
        ]
    },
    {
        "file": "strings_notif_comeback_missed_days.xml",
        "group_id": "comeback_missed_days",
        "name_pt": "Retomada e Dias Ausentes",
        "name_en": "Comeback and Missed Days",
        "items": [
            ("COMEBACK_SHORT_ABSENCE", "notif_comeback_short_title", "notif_comeback_short_messages"),
            ("COMEBACK_WEEK_ABSENCE", "notif_comeback_week_title", "notif_comeback_week_messages"),
            ("COMEBACK_LONG_ABSENCE", "notif_comeback_long_title", "notif_comeback_long_messages"),
        ]
    },
    {
        "file": "strings_notif_in_workout_performance.xml",
        "group_id": "in_workout_performance",
        "name_pt": "Performance Durante o Treino",
        "name_en": "In-Workout Performance",
        "items": [
            ("IN_WORKOUT_PR", "notif_perf_pr_title", "notif_perf_pr_messages"),
            ("IN_WORKOUT_HALFWAY", "notif_perf_halfway_title", "notif_perf_halfway_messages"),
            ("IN_WORKOUT_REST_OVERRUN", "notif_perf_rest_overrun_title", "notif_perf_rest_overrun_messages"),
            ("IN_WORKOUT_LAST_EXERCISE", "notif_perf_last_exercise_title", "notif_perf_last_exercise_messages"),
            ("IN_WORKOUT_LOAD_PROGRESSION", "notif_perf_load_progression_title", "notif_perf_load_progression_messages"),
        ]
    },
    {
        "file": "strings_notif_muscle_evolution_balance.xml",
        "group_id": "muscle_evolution_balance",
        "name_pt": "Evolução e Equilíbrio Muscular",
        "name_en": "Muscle Evolution and Balance",
        "items": [
            ("EVOLUTION_NEGLECTED_MUSCLE", "notif_evo_neglected_muscle_title", "notif_evo_neglected_muscle_messages"),
            ("EVOLUTION_VOLUME_MILESTONE", "notif_evo_volume_milestone_title", "notif_evo_volume_milestone_messages"),
            ("EVOLUTION_HEALTH_BENEFIT", "notif_evo_health_benefit_title", "notif_evo_health_benefit_messages"),
        ]
    }
]

def main():
    source_pt_dir = "/Volumes/Desenvolvimento_hub/FitGuide/app/src/main/res/values-pt"
    source_en_dir = "/Volumes/Desenvolvimento_hub/FitGuide/app/src/main/res/values-en"
    base_out = "/Volumes/Desenvolvimento_hub/fitguide-exercises/database/notifications"
    os.makedirs(os.path.join(base_out, "ptbr"), exist_ok=True)
    os.makedirs(os.path.join(base_out, "en"), exist_ok=True)

    for g in groups_def:
        pt_s, pt_a = parse(os.path.join(source_pt_dir, g["file"]))
        en_s, en_a = parse(os.path.join(source_en_dir, g["file"]))

        pt_items = []
        en_items = []

        for item_id, title_key, msg_key in g["items"]:
            pt_title = pt_s.get(title_key, "")
            en_title = en_s.get(title_key, "")

            if msg_key in pt_a:
                pt_msgs = pt_a[msg_key]
            elif msg_key in pt_s:
                pt_msgs = [pt_s[msg_key]]
            else:
                pt_msgs = []

            if msg_key in en_a:
                en_msgs = en_a[msg_key]
            elif msg_key in en_s:
                en_msgs = [en_s[msg_key]]
            else:
                en_msgs = []

            pt_items.append({
                "id": item_id,
                "title": pt_title,
                "messages": pt_msgs
            })
            en_items.append({
                "id": item_id,
                "title": en_title,
                "messages": en_msgs
            })

        pt_doc = {
            "group": g["group_id"],
            "name": g["name_pt"],
            "items": pt_items
        }
        en_doc = {
            "group": g["group_id"],
            "name": g["name_en"],
            "items": en_items
        }

        with open(os.path.join(base_out, "ptbr", f"{g['group_id']}.json"), "w", encoding="utf-8") as f:
            json.dump(pt_doc, f, ensure_ascii=False, indent=2)
            f.write("\n")

        with open(os.path.join(base_out, "en", f"{g['group_id']}.json"), "w", encoding="utf-8") as f:
            json.dump(en_doc, f, ensure_ascii=False, indent=2)
            f.write("\n")

    print("OK")

if __name__ == "__main__":
    main()
