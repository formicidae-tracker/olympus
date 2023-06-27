#!/bin/sh

MTX_PATH=$1


DIR=$(/usr/bin/dirname $1)

mkdir -p /srv/thumbnails/$DIR

while true
do
	ffmpeg -i rtmp://localhost/$MTX_PATH -frames:v 1 /srv/thumbnails/$MTX_PATH.jpg 2>/dev/null
	/bin/sleep 600
done
