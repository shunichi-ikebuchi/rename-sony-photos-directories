#!/usr/bin/env bash
# sync.sh — gate-kit の go/ を消費側の repo に commit sha で固定して写す。
#
# Go repo は npm を使わないので、npm 版と同じ「sha で固定して取り込む」を make 断片の
# vendor で行う。写しは消費側の repo に commit するので、clone した人 (共同開発者を
# 含む) は gate-kit への権限も network も無しに同じ gate を回せる。
#
# 写すもの: go/ (tests を除く) + LICENSE + NOTICE。添えるもの:
#   REVISION     取り込んだ commit sha
#   SHA256SUMS   写した file の hash。verify.sh が照合する (手で書き換えた写しを落とす)
#
# 使い方: sync.sh <consumer repo root> <40 桁の commit sha>
# Env: GATE_KIT_SRC  gate-kit の clone (既定: この script の 1 つ上 = gate-kit の checkout)
set -euo pipefail

usage() {
	echo "usage: sync.sh <consumer repo root> <40 桁の commit sha>" >&2
	exit 2
}

[ $# -eq 2 ] || usage
CONSUMER="$1"
REF="$2"
# branch 名や短縮 sha は受けない: branch は上流の変更が予告なく入り、短縮 sha は将来曖昧になる
[[ "$REF" =~ ^[0-9a-f]{40}$ ]] || {
	echo "REF は 40 桁の commit sha で指定する (branch / tag / 短縮 sha は不可): $REF" >&2
	usage
}
[ -d "$CONSUMER" ] || {
	echo "consumer の dir が無い: $CONSUMER" >&2
	exit 2
}

SRC="${GATE_KIT_SRC:-$(cd "$(dirname "$0")/.." && pwd)}"
git -C "$SRC" cat-file -e "${REF}^{commit}" 2>/dev/null || {
	echo "gate-kit に commit ${REF} が無い (git -C $SRC fetch してから)" >&2
	exit 1
}

DEST="$CONSUMER/third_party/gate-kit"
stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT

# 作業 tree ではなく commit の中身を写す (手元の未 commit の変更を混ぜない)
git -C "$SRC" archive "$REF" go LICENSE NOTICE | tar -x -C "$stage"
rm -rf "$stage/go/tests"
echo "$REF" >"$stage/REVISION"
(cd "$stage" && find go LICENSE NOTICE -type f | LC_ALL=C sort | xargs shasum -a 256) >"$stage/SHA256SUMS"

# 前の写しは丸ごと置き換える (消えた file が残らないように)
rm -rf "$DEST"
mkdir -p "$(dirname "$DEST")"
cp -R "$stage" "$DEST"
chmod +x "$DEST"/go/*.sh

echo "✓ gate-kit ${REF:0:12} を $(basename "$CONSUMER")/third_party/gate-kit に写した ($(wc -l <"$DEST/SHA256SUMS" | tr -d ' ') files)"
