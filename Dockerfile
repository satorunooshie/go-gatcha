FROM golang:1.15.2-alpine
# アップデートとgitのインストール
RUN apk update && apk add git
# appディレクトリの作成
RUN mkdir /go/src/app
# ワーキングディレクトリの設定
WORKDIR /go/src/app
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/app
RUN go get -u github.com/oxequa/realize \
    # sqlを使うためのモジュール
    && go get github.com/go-sql-driver/mysql
CMD ["realize", "start"]