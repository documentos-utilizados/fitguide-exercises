.PHONY: all lint lint_pt check_dupes install metadata workouts admin-build admin-run docker-build docker-up docker-down

sources :=$(wildcard ./database/exercises/en/*.json)
sources_pt :=$(wildcard ./database/exercises/pt/*.json)
sources_workouts :=$(wildcard ./database/workouts/*/*.json)

all: database/dist/exercises_en.json database/dist/exercises_pt.json metadata workouts

metadata: database/dist/exercises_en.json database/dist/exercises_pt.json
	python3 scripts/extract_metadata.py

workouts: database/dist/exercises_en.json database/dist/exercises_pt.json $(sources_workouts)
	python3 scripts/compile_workouts.py

lint:
	check-jsonschema --schemafile ./database/schemas/exercise.schema.json $(sources)

lint_pt:
	check-jsonschema --schemafile ./database/schemas/exercise.schema.json $(sources_pt)

check_dupes:
	@echo "Checking duplicate IDs in English exercises..."
	@python3 -c 'import glob, json, collections; ids = [json.load(open(f))["id"] for f in glob.glob("./database/exercises/en/*.json")]; dupes = [id for id, count in collections.Counter(ids).items() if count > 1]; print("Duplicates found:", dupes if dupes else "None")'
	@echo "Checking duplicate IDs in Portuguese exercises..."
	@python3 -c 'import glob, json, collections; ids = [json.load(open(f))["id"] for f in glob.glob("./database/exercises/pt/*.json")]; dupes = [id for id, count in collections.Counter(ids).items() if count > 1]; print("Duplicates found:", dupes if dupes else "None")'

install:
	pip install check-jsonschema

database/dist/exercises_en.json: $(sources)
	python3 scripts/compile_en.py

database/dist/exercises_pt.json: $(sources_pt)
	python3 scripts/compile_pt.py

admin-build:
	cd admin/web && npm install && npm run build
	cd admin/server && go build -o fitguide-admin .

admin-run: admin-build
	cd admin/server && ./fitguide-admin

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down
