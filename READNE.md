# rename-sony-photos-directores

このプロジェクトはSonyのデジタルカメラの[日付形式のフォルダ](https://www.sony.jp/ServiceArea/impdf/pdf/44879440M.w-JP/jp/contents/TP0000220296.html)をyyyy-mm-dd形式に変更するためのプログラムとシェルスクリプトが含まれています。

# 前提条件

- macOS
  - このレポジトリに含まれているシェルスクリプトはmacOSでのみ有効なコマンドが含まれています
- 写真を保存したいSDカードの名前は`1-1`になっていること
- 写真を保存したいSDカードのパスは`/Volumes/1-1`になっていること
- バックアップ用のSDカードの名前は`1-2`になっていること
- バックアップ用のSDカードのパスは`/Volumes/1-2`になっていること



# 使い方

## 事前準備

`go install`してPATHの通った場所に`rename-sony-photos-directores`を配置してください。



## フォルダ名を変更して写真をコピーする

```
./copy-rename-and-delete.sh
```

## バックアップ用のSDカードの写真を削除する
```
./delete.sh
```