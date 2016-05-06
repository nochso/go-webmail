#!/usr/bin/env bash
rm -f data/mail.sqlite model/*.xo.go
go-bindata asset/model.sql LICENSE
mkdir -p model
cat asset/model.sql | sqlite3 -batch data/mail.sqlite
xo sqlite3://data/mail.sqlite -o model --int32-type int64 --uint32-type uint64
xo sqlite3://data/mail.sqlite -o model -a -N -M -B -T Flag -F GetFlags << ENDSQL
SELECT
  f.id,
  f.name
FROM flag f
ENDSQL

# With tag: 1.0.0[-commits since tag][-dirty]
# Without existing tags: 4a4154a[-dirty]
version=$(git describe --tags --always --dirty)
# 2000-12-13 13:00:00
buildDate=$(date --utc "+%F %H:%M %Z")
go build -ldflags "-X 'main.Version=$version' -X 'main.BuildDate=$buildDate'"
echo Built version $version
