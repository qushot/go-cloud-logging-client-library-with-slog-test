# go-cloud-logging-client-library-with-slog-test

## 概要

Cloud Logging の[クライアントライブラリ](https://pkg.go.dev/cloud.google.com/go/logging)を試してみました。

クライアントライブラリの概要

- 良い点
  - ロガーを生成する際に、 `logging.RedirectAsJSON(os.Stdout)` とすることで、出力先を標準出力に変更できる。
  - \*http.Request を渡すことで TraceID， SpanID を自動で付与してくれる。
- 微妙な点
  - TraceID， SpanID を付与するために \*http.Request を渡す必要がある。
  - jsonPayload の形式がいまいち（トップレベルの `message` に対してオブジェクトが入っている）。Cloud Logging 上での見た目が悪い。
    - 期待する形式
      ```json
      {
        "message": "Hello, Cloud Logging!",
        "name": "Takashi",
        "age": 30,
        "severity": "INFO"
      }
      ```
    - 実際の形式
      ```json
      {
        "message": {
          "message": "Hello, Cloud Logging!",
          "name": "Takashi",
          "age": 30
        },
        "severity": "INFO"
      }
      ```

jsonPayload の形式を期待する形式で出力するために、 `logger.ToLogEntry` から構造体を取得して slog で整形してみたが、 http.Request を渡さないといけない点は変わらないし、自前で TraceID, SpanID を取得する処理を実装し、 slog や zap を使えばいいのでは？という気持ちになった。

## その他

実行などに必要なコマンドは Makefile に書いてあるので、それを参照してください。
