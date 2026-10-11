#!/usr/bin/env bash
# check-go-toolchain.sh — Go の版の宣言 (go.mod / Dockerfile) と手元の Go が食い違っていないかを見る。
#
# Go の版は 3 か所で決まり、どれも別の物を決める:
#   - 手元の Go         … gate (vet / lint / govulncheck / test) が見る標準ライブラリ
#   - Dockerfile golang  … 出荷物の Go。手元で通った gate とは別の標準ライブラリで build される
#   - go.mod go / toolchain … 必要な最低版。GOTOOLCHAIN=local では toolchain 行は手元を変えない
#
# 見るもの (三状態):
#   異常 (1)
#     - Dockerfile の golang が go.mod の go より古い (build できない / 宣言が嘘)
#     - Dockerfile の golang と手元の Go の minor が違う (gate が見た標準ライブラリと出荷物が違う)
#     - GOTOOLCHAIN=local なのに toolchain 行が手元より新しい (宣言が効いていない)
#   確かめきれていない (3)
#     - Dockerfile の golang が patch まで固定されていない / 変数 (何が出荷されるか読めない)
#   正常 (0)
#   使い方の誤り (2) — 対象 dir に go.mod が無い
#
# patch だけの違い (手元 1.27.2 / 出荷 1.27.1) は落とさず出力に残す。
#
# 使い方: check-go-toolchain.sh [<repo root>]   (既定: カレント)
# Env (テスト用): GATE_KIT_LOCAL_GO (既定 `go env GOVERSION`) / GATE_KIT_GOTOOLCHAIN (既定 `go env GOTOOLCHAIN`)
set -uo pipefail

ROOT="${1:-.}"
[ -f "$ROOT/go.mod" ] || {
	echo "usage: check-go-toolchain.sh [<repo root>] (go.mod が無い: $ROOT)" >&2
	exit 2
}

LOCAL_GO="${GATE_KIT_LOCAL_GO:-$(go env GOVERSION 2>/dev/null)}"
GOTOOLCHAIN_MODE="${GATE_KIT_GOTOOLCHAIN:-$(go env GOTOOLCHAIN 2>/dev/null)}"
LOCAL_GO="${LOCAL_GO#go}"

bad=0
unsure=0
report() { # report <verdict> <file> <message>
	printf '%-18s %-36s %s\n' "$1" "$2" "$3"
	case "$1" in
	異常) bad=1 ;;
	確かめきれていない) unsure=1 ;;
	esac
}

# ver_cmp <a> <b> → -1 / 0 / 1 (1.27 と 1.27.0 は同じ扱い)
ver_cmp() {
	local IFS=. i a b
	# shellcheck disable=SC2206
	local x=($1) y=($2)
	for i in 0 1 2; do
		a="${x[$i]:-0}"
		b="${y[$i]:-0}"
		if ((10#$a < 10#$b)); then
			echo -1
			return
		elif ((10#$a > 10#$b)); then
			echo 1
			return
		fi
	done
	echo 0
}

minor_of() { printf '%s' "$1" | cut -d. -f1-2; }

# go.mod の go / toolchain 行
go_directive() { sed -nE 's/^go[[:space:]]+([0-9][0-9.]*)[[:space:]]*$/\1/p' "$1" | head -1; }
toolchain_directive() { sed -nE 's/^toolchain[[:space:]]+go([0-9][0-9.]*)[[:space:]]*$/\1/p' "$1" | head -1; }

# <dir> から上に辿って一番近い go.mod
nearest_gomod() {
	local d="$1"
	while :; do
		[ -f "$d/go.mod" ] && {
			echo "$d/go.mod"
			return
		}
		[ "$d" = "$ROOT" ] && return
		d="$(dirname "$d")"
	done
}

rel() { printf '%s' "${1#"$ROOT"/}"; }

# ── go.mod の toolchain 行 ──────────────────────────
while IFS= read -r gomod; do
	tc="$(toolchain_directive "$gomod")"
	[ -n "$tc" ] || continue
	if [ "$GOTOOLCHAIN_MODE" = "local" ] && [ -n "$LOCAL_GO" ] && [ "$(ver_cmp "$tc" "$LOCAL_GO")" = 1 ]; then
		report 異常 "$(rel "$gomod")" "toolchain go${tc} は GOTOOLCHAIN=local では効かない (手元は go${LOCAL_GO})"
	fi
done < <(find "$ROOT" -name go.mod -not -path '*/third_party/*' -not -path '*/node_modules/*' -not -path '*/vendor/*' -not -path '*/.git/*' | LC_ALL=C sort)

# ── Dockerfile の golang ────────────────────────────
n_docker=0
while IFS= read -r df; do
	from="$(grep -iE '^[[:space:]]*FROM[[:space:]]' "$df" | grep -oE '(docker\.io/(library/)?)?golang:[^[:space:]]+' | head -1)"
	[ -n "$from" ] || continue
	n_docker=$((n_docker + 1))
	tag="${from#*golang:}"
	gomod="$(nearest_gomod "$(dirname "$df")")"
	[ -n "$gomod" ] || gomod="$ROOT/go.mod"
	want="$(go_directive "$gomod")"

	ship="$(printf '%s' "$tag" | sed -nE 's/^([0-9]+\.[0-9]+\.[0-9]+)([-@].*)?$/\1/p')"
	if [ -z "$ship" ]; then
		# patch が無くても minor だけで「go.mod より古い」と言い切れる時は異常
		ship_minor="$(printf '%s' "$tag" | sed -nE 's/^([0-9]+\.[0-9]+)([-@].*)?$/\1/p')"
		if [ -n "$ship_minor" ] && [ -n "$want" ] && [ "$(ver_cmp "$ship_minor" "$(minor_of "$want")")" = -1 ]; then
			report 異常 "$(rel "$df")" "golang:${tag} が $(rel "$gomod") の go ${want} より古い"
			continue
		fi
		report 確かめきれていない "$(rel "$df")" "golang:${tag} は patch まで固定されていない (出荷される Go が読めない)"
		continue
	fi
	if [ -n "$want" ] && [ "$(ver_cmp "$ship" "$want")" = -1 ]; then
		report 異常 "$(rel "$df")" "golang:${ship} が $(rel "$gomod") の go ${want} より古い"
		continue
	fi
	if [ -n "$LOCAL_GO" ] && [ "$(minor_of "$ship")" != "$(minor_of "$LOCAL_GO")" ]; then
		report 異常 "$(rel "$df")" "出荷物 go${ship} と手元 go${LOCAL_GO} の minor が違う (gate が見た標準ライブラリと別物を出荷する)"
		continue
	fi
	if [ -n "$LOCAL_GO" ] && [ "$ship" != "$LOCAL_GO" ]; then
		report 正常 "$(rel "$df")" "出荷物 go${ship} / 手元 go${LOCAL_GO} (patch だけの違い)"
	else
		report 正常 "$(rel "$df")" "出荷物 go${ship} / 手元 go${LOCAL_GO}"
	fi
done < <(find "$ROOT" \( -name Dockerfile -o -name 'Dockerfile.*' -o -name '*.Dockerfile' \) -not -path '*/third_party/*' -not -path '*/node_modules/*' -not -path '*/vendor/*' -not -path '*/.git/*' | LC_ALL=C sort)

[ "$n_docker" -eq 0 ] && report 正常 "go.mod" "Dockerfile なし (go ${LOCAL_GO:+/ 手元 go${LOCAL_GO}})"

if [ "$bad" -eq 1 ]; then exit 1; fi
if [ "$unsure" -eq 1 ]; then exit 3; fi
exit 0
