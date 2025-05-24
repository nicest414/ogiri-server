# 大喜利サーバー

大喜利のお題や回答を管理するためのシンプルなGoサーバーです。

## 機能

- お題の作成・取得・更新・削除
- 回答の投稿・取得・更新・削除
- メモリ内データストア（再起動するとデータは消えます）

## セットアップと実行方法

### 前提条件

- [Go](https://golang.org/dl/) 1.18以上がインストールされていること

### インストールと実行

1. リポジトリをクローン

```bash
git clone https://github.com/yourusername/ogiri-server.git
cd ogiri-server
```

2. 依存関係のダウンロード

```bash
go mod download
go mod init ogiri-server
go get github.com/gorilla/mux
```

3. サーバーの実行

```bash
go run cmd/api/main.go
```

または、Windowsの場合は以下のバッチファイルを実行：

```
run.bat
```

サーバーは http://localhost:8080 で起動します。

## API エンドポイント

### お題関連

- `GET /api/themes` - すべてのお題を取得
- `POST /api/themes` - 新しいお題を作成
- `GET /api/themes/{id}` - 特定のお題を取得
- `PUT /api/themes/{id}` - お題を更新
- `DELETE /api/themes/{id}` - お題を削除

### 回答関連

- `GET /api/themes/{themeID}/answers` - お題に対するすべての回答を取得
- `POST /api/themes/{themeID}/answers` - お題に対して新しい回答を投稿
- `GET /api/themes/{themeID}/answers/{id}` - 特定の回答を取得
- `PUT /api/themes/{themeID}/answers/{id}` - 回答を更新
- `DELETE /api/themes/{themeID}/answers/{id}` - 回答を削除

## リクエスト/レスポンス例

### お題の作成

リクエスト:
```json
POST /api/themes
{
  "title": "猫と和解する方法",
  "description": "怒っている猫と仲直りするユニークな方法を考えてください",
  "created_by": "管理者"
}
```

レスポンス:
```json
{
  "id": "a1b2c3d4",
  "title": "猫と和解する方法",
  "description": "怒っている猫と仲直りするユニークな方法を考えてください",
  "created_at": "2023-06-15T12:34:56Z",
  "updated_at": "2023-06-15T12:34:56Z",
  "created_by": "管理者",
  "active": true
}
```

### 回答の投稿

リクエスト:
```json
POST /api/themes/a1b2c3d4/answers
{
  "content": "猫に「ごめんね」と言いながら、自分も床で寝転がって目をウインクする",
  "created_by": "ねこ好き"
}
```

レスポンス:
```json
{
  "id": "e5f6g7h8",
  "theme_id": "a1b2c3d4",
  "content": "猫に「ごめんね」と言いながら、自分も床で寝転がって目をウインクする",
  "created_at": "2023-06-15T13:45:12Z",
  "updated_at": "2023-06-15T13:45:12Z",
  "created_by": "ねこ好き",
  "likes": 0
}
```
