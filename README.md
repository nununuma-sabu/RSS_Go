# RSS_Go

RSSフィードを定期的に取得・要約し、Slackへ通知するバッチ処理システム（Go言語製）です。

## 進捗状況（フェーズ）

- [x] **フェーズ 1**: RSS取得基盤の実装 (YAML設定ファイルからのURL読み込みとパース)
- [x] **フェーズ 2**: 過去記事との重複排除ロジック
- [x] **フェーズ 3**: Gemini APIを用いた要約処理
- [x] **フェーズ 4**: GitHub Actionsでの定期実行とMarkdown/JSON生成
- [ ] **フェーズ 5**: Slack通知連携

## ローカルでの実行方法

事前にGoがインストールされている必要があります。また、Gemini APIを利用するため環境変数 `GEMINI_API_KEY` を設定してください。

1. リポジトリをクローンし、ディレクトリに移動します。
2. 依存パッケージをダウンロードします。
   ```bash
   go mod tidy
   ```
3. プログラムを実行します。
   ```bash
   export GEMINI_API_KEY="あなたのAPIキー"
   go run main.go
   ```

### 実行例

```text
$ go run main.go
2026/05/18 15:53:46 Starting RSS fetcher...
2026/05/18 15:53:46 Fetching feed: https://yamadashy.github.io/tech-blog-rss-feed/feeds/rss.xml (tech)
2026/05/18 15:53:46 Successfully fetched: 企業テックブログRSS
2026/05/18 15:53:46 Found 413 items
2026/05/18 15:53:46   - [tech] 岐阜から通うPlatformエンジニアのリアル：中部支店という選択 | Sansan Tech Blog (https://buildersbox.corp-sansan.com/entry/2026/05/18/150000)
...
2026/05/18 15:53:46 RSS fetching completed.
```

## 設定

`config.yaml` に取得したいRSSフィードのURLとカテゴリを追加して実行することで、複数のフィードをまとめて取得可能です。