#!/usr/bin/env bash
rm data/mail.sqlite models/*.xo.go
go-bindata model.sql LICENSE
cat model.sql | sqlite3 -batch data/mail.sqlite
xo sqlite3://data/mail.sqlite -o models --int32-type int64 --uint32-type uint64
xo sqlite3://data/mail.sqlite -o models -a -N -M -B -T Flag -F GetFlags << ENDSQL
SELECT
  f.id,
  f.name
FROM flag f
ENDSQL

# With tag: 1.0.0[-commits since tag][-dirty]
# Without existing tags: 4a4154a[-dirty]
version=$(git describe --tags --always --dirty)
# 2000-12-13 13:00:00
buildDate=$(date --utc "+%F %T")
go build -ldflags "-X 'main.Version=$version' -X 'main.BuildDate=$buildDate'"
echo Built version $version
