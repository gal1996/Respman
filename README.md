# resp-man
チーム開発でサーバーサイドとフロントサイドが同時に開発している場合、ドキュメントなどを用意しI/Fを決めていると思いますが、
フロントの開発では表示のためにmockデータを作ったり、サーバーサイドはレスポンスが間違ってないか確認したりと言った作業が発生すると思います。
それを簡略化してくれるやつです。

# 機能
指定したjsonに関して以下の二つの機能を持つローカルサーバーを起動します。
- 指定したjsonをそのままレスポンスとして返す `/return_response`
- postしたjsonと指定したjsonのkeyが一致しているかをチェックして結果を返す `/check_response`

# 使い方
1. `response.txt` にチェックしたいレスポンスをそのままドキュメントからコピペ
2. resp-manを起動  
    - `$ cd resp-man`
    - `$ make start`
   
3. 使う
   - レスポンスを受け取りたい場合
      - `curl -X POST 'http://localhost:3030/return_response' | jq`
      - または、フロントのリクエスト先として登録
   - レスポンスのチェックをしたい場合
      - `curl -X POST 'http://localhost:3030/check_response' -d '{
        "key1": 1,
        "key2": "message"
        }'`
      - 改行にまだ対応できてないです。