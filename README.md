# RSS_Go

RSSフィードを定期的に取得し、新着記事をGemini API（AI）で要約してMarkdown保存・Slack通知を行うバッチ処理システムです。
GitHub Actionsを利用することで、完全無料かつサーバーレスで自動運用が可能です。

## システム構成・フロー

1. **GitHub Actions（定期実行）**: 毎日指定時間（デフォルトは日本時間 深夜0:00）に起動
2. **Goバッチ処理**: `config.yaml` に定義されたRSSリストをパースして記事を取得
3. **重複排除**: `fetched_articles.json` を参照し、未読の新着記事のみを抽出
4. **AI要約**: `gemini-2.5-pro` モデルを用いて、各記事を3行で簡潔に要約
5. **Markdown保存**: 日付ごとのファイル（例：`summaries/2026-05-18.md`）に要約を書き出し、リポジトリに自動Push
6. **Slack通知**: SlackのIncoming Webhook経由で指定チャンネルに要約テキストを自動送信

## 使い方：RSSフィードの追加方法

新しいRSSフィードを追加したい場合は、プロジェクト直下にある `config.yaml` を編集するだけです。
カテゴリ名（`category`）は自由に設定でき、MarkdownやSlack通知の際にタグとして表示されます。

```yaml
feeds:
  - url: "https://yamadashy.github.io/tech-blog-rss-feed/feeds/rss.xml"
    category: "tech"
  - url: "https://example.com/fitness/rss"
    category: "fitness"
  - url: "https://example.com/nutrition/feed"
    category: "nutrition"
```

## 必要な環境変数・シークレット

システムを稼働させるためには、以下の環境変数（GitHub Actionsの場合はRepository Secrets）が必要です。

| 変数名 | 必須 | 説明 |
|---|---|---|
| `GEMINI_API_KEY` | **必須** | Google AI Studioから取得したGeminiのAPIキー |
| `SLACK_WEBHOOK_URL` | 任意 | SlackのIncoming Webhook URL（未設定の場合は通知をスキップ） |

## GitHub Actionsでの運用方法（本番環境）

1. リポジトリの **Settings > Secrets and variables > Actions** にて、上記の2つの変数を登録します。
2. **Settings > Actions > General > Workflow permissions** にて **`Read and write permissions`** を有効にします。
3. あとは毎日自動で実行され、`summaries/` フォルダにMarkdownが蓄積されていきます。
   ※ 手動で実行したい場合は、`Actions` タブから `Run workflow` をクリックしてください。

## ローカルでの実行方法（開発環境）

Go 1.26以上がインストールされている必要があります。

```bash
# 依存パッケージのダウンロード
go mod tidy

# 環境変数のセット
export GEMINI_API_KEY="あなたのAPIキー"
export SLACK_WEBHOOK_URL="あなたのWebhook URL" # (任意)

# 実行
go run main.go
```

## 開発の進捗（完了フェーズ）

- [x] **フェーズ 1**: RSS取得基盤の実装 (YAML設定ファイルからのURL読み込みとパース)
- [x] **フェーズ 2**: 過去記事との重複排除ロジック
- [x] **フェーズ 3**: Gemini APIを用いた要約処理 (`gemini-2.5-pro` 対応)
- [x] **フェーズ 4**: GitHub Actionsでの定期実行とMarkdown/JSON生成
- [x] **フェーズ 5**: Slack通知連携