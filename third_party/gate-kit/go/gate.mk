# gate-kit/go/gate.mk — Go repo の gate 一式。消費側は写し (third_party/gate-kit) を include する。
#
#   include third_party/gate-kit/go/gate.mk
#   ci: $(GATE_KIT_GO_GATES) <repo 固有の gate>
#
# 写しは書き換えない (gate-kit-verify が落とす)。gate を外すなら include の前に
#   GATE_KIT_GO_SKIP := go-vuln      # 理由をコメントで書く
# 外したものは `make gate-kit-gates` に skip として出る (黙って消えない)。
#
# 取り込み直し: make -f third_party/gate-kit/go/gate.mk gate-kit-sync REF=<40 桁 sha> GATE_KIT_SRC=<gate-kit の clone>
#
# 解析器 (golangci-lint / govulncheck) は手元の Go で build したものを使う。古い Go で build した
# 解析器は新しい標準ライブラリを読めない (2026-10-09 実測) ので go-tools-fresh で落とす。

GATE_KIT_GO_DIR := $(patsubst %/,%,$(dir $(lastword $(MAKEFILE_LIST))))
GATE_KIT_ROOT := $(patsubst %/,%,$(dir $(GATE_KIT_GO_DIR)))

GATE_KIT_GO_ALL := gate-kit-verify go-toolchain go-tools-fresh go-mod-tidy go-vet go-lint go-vuln
GATE_KIT_GO_SKIP ?=
GATE_KIT_GO_GATES := $(filter-out $(GATE_KIT_GO_SKIP),$(GATE_KIT_GO_ALL))

GATE_KIT_GOLANGCI ?= golangci-lint
GATE_KIT_GOVULNCHECK ?= govulncheck

.PHONY: $(GATE_KIT_GO_ALL) gate-kit-gates gate-kit-sync

gate-kit-gates:
	@for g in $(GATE_KIT_GO_GATES); do echo "gate: $$g"; done
	@for g in $(GATE_KIT_GO_SKIP); do echo "skip: $$g"; done

# 写しが取り込んだ時のまま (手で書き換えていない) か
gate-kit-verify:
	@bash $(GATE_KIT_GO_DIR)/verify.sh $(GATE_KIT_ROOT)

# go.mod / Dockerfile の Go の版と手元の Go が食い違っていないか
go-toolchain:
	@bash $(GATE_KIT_GO_DIR)/check-go-toolchain.sh .

# 解析器を build した Go が手元の Go と同じか (古ければ標準ライブラリを読めず誤検知 / 見逃しになる)
go-tools-fresh:
	@want="$$(go env GOVERSION)"; rc=0; \
	for t in $(GATE_KIT_GOLANGCI) $(GATE_KIT_GOVULNCHECK); do \
	  p="$$(command -v $$t)" || { echo "✗ $$t が無い"; rc=1; continue; }; \
	  got="$$(go version -m "$$p" 2>/dev/null | head -1 | awk '{print $$2}')"; \
	  if [ "$$got" != "$$want" ]; then echo "✗ $$t は $$got で build されている (手元は $$want)。go install し直す"; rc=1; fi; \
	done; exit $$rc

# go mod tidy で go.mod / go.sum が変わらないか (変わるなら tidy して commit する)
go-mod-tidy:
	@go mod tidy
	@git diff --quiet -- go.mod go.sum || { \
	  echo "✗ go mod tidy で go.mod / go.sum が変わった。tidy した結果を commit する" >&2; \
	  git --no-pager diff --stat -- go.mod go.sum >&2; exit 1; }

go-vet:
	@go vet ./...

go-lint:
	@command -v $(GATE_KIT_GOLANGCI) >/dev/null || { \
	  echo "✗ golangci-lint が無い: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@<版>" >&2; exit 1; }
	@$(GATE_KIT_GOLANGCI) run ./...

go-vuln:
	@command -v $(GATE_KIT_GOVULNCHECK) >/dev/null || { \
	  echo "✗ govulncheck が無い: go install golang.org/x/vuln/cmd/govulncheck@<版>" >&2; exit 1; }
	@$(GATE_KIT_GOVULNCHECK) ./...

# gate-kit の clone から写しを取り込み直す (REF は 40 桁の commit sha)
gate-kit-sync:
	@test -n "$(REF)" || { echo "usage: make gate-kit-sync REF=<40 桁 sha> GATE_KIT_SRC=<gate-kit の clone>" >&2; exit 2; }
	@test -n "$(GATE_KIT_SRC)" || { echo "GATE_KIT_SRC (gate-kit の clone) を指定する" >&2; exit 2; }
	@GATE_KIT_SRC="$(GATE_KIT_SRC)" bash "$(GATE_KIT_SRC)/go/sync.sh" . "$(REF)"
