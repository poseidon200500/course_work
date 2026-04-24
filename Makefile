APP=bench

# ===== PATHS =====
BENCH_CMD=./cmd/benchmarks
GEN_CMD=./cmd/generator
ANALYSIS_SCRIPT=./internal/analysis/analyze_benchmarks_diploma_v2.py

DATA_DIR=datastore
OUT_CSV=benchmark_results.csv
PLOTS_DIR=plots_diploma

# ===== DEFAULT PARAMS =====
TOTAL?=1000000
DUP?=40
MAXLEN?=8
DIST?=uniform
STORAGE?=all

# =========================
# ===== BENCHMARKS ========
# =========================

run:
	go run $(BENCH_CMD) \
		-name=TEST \
		-total=$(TOTAL) \
		-dup=$(DUP) \
		-maxlen=$(MAXLEN) \
		-dist=$(DIST) \
		-storage=$(STORAGE) \
		-out=$(OUT_CSV)

quick:
	make run TOTAL=50000 DUP=40 MAXLEN=8

# ===== SCENARIOS =====

dup10:
	make run DUP=10

dup40:
	make run DUP=40

dup80:
	make run DUP=80

len4:
	make run MAXLEN=4

len8:
	make run MAXLEN=8

len16:
	make run MAXLEN=16

unique:
	make run DUP=95 MAXLEN=32 DIST=zipf

# ===== STORAGE =====

base:
	make run STORAGE=base

intern:
	make run STORAGE=intern

v1:
	make run STORAGE=v1

v2:
	make run STORAGE=v2

# =========================
# ===== GENERATOR =========
# =========================

generate:
	go run $(GEN_CMD)

# =========================
# ===== ANALYSIS ==========
# =========================

analyze:
	python3 $(ANALYSIS_SCRIPT)

# =========================
# ===== PIPELINE ==========
# =========================

full:
	make run
	make analyze

# =========================
# ===== CLEAN =============
# =========================

clean:
	rm -rf $(DATA_DIR)/*
	rm -f $(OUT_CSV)
	rm -rf $(PLOTS_DIR)

# =========================
# ===== MATRIX ============
# =========================

matrix:
	make dup10
	make dup40
	make dup80
	make len4
	make len8
	make len16
	make unique