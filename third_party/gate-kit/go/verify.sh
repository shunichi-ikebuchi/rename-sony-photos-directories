#!/usr/bin/env bash
# verify.sh — gate-kit の写し (third_party/gate-kit) が取り込んだ時のまま (sync.sh の出力のまま) かを見る。
#
# 手でコピーした gate は、写した先で少しずつ書き換わって元とずれる (実測: practice-os が
# eir の gate をコピーして lint と vuln が落ち、check-gofumpt.sh は 111 行食い違っていた)。
# 写しの中の file を SHA256SUMS と照合し、書き換え / 追加 / 削除のどれでも落とす。
# 直すなら写しを書き換えずに gate-kit を直して sync し直す。
#
# network も gate-kit への権限も要らない (共同開発者の clone でも回る)。
#
# 使い方: verify.sh [<写しの root>]   (既定: この script の 1 つ上)
# exit: 0 = 一致 / 1 = 食い違い / 2 = 照合できない (SHA256SUMS か REVISION が無い)
set -uo pipefail

ROOT="${1:-$(cd "$(dirname "$0")/.." && pwd)}"
SUMS="$ROOT/SHA256SUMS"

if [ ! -f "$SUMS" ] || [ ! -f "$ROOT/REVISION" ]; then
	echo "✗ gate-kit の写しを照合できない: SHA256SUMS か REVISION が無い ($ROOT)。sync.sh で取り込み直す" >&2
	exit 2
fi

bad=0

# 書き換え / 削除
while read -r want path; do
	if [ ! -f "$ROOT/$path" ]; then
		echo "✗ 写しから消えている: $path"
		bad=1
		continue
	fi
	got="$(shasum -a 256 "$ROOT/$path" | cut -d' ' -f1)"
	if [ "$got" != "$want" ]; then
		echo "✗ 写しが書き換わっている: $path"
		bad=1
	fi
done <"$SUMS"

# 追加
while IFS= read -r path; do
	if ! awk -v p="$path" '$2 == p {found=1} END {exit !found}' "$SUMS"; then
		echo "✗ 写しに無いはずの file がある: $path"
		bad=1
	fi
done < <(cd "$ROOT" && find go LICENSE NOTICE -type f 2>/dev/null | LC_ALL=C sort)

rev="$(cut -c1-12 "$ROOT/REVISION")"
if [ "$bad" -eq 1 ]; then
	echo "✗ gate-kit ${rev} の写しが食い違っている。写しを直さず gate-kit を直して sync し直す" >&2
	exit 1
fi
echo "✓ gate-kit ${rev} の写しは取り込んだ時のまま"
