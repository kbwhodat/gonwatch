#!/usr/bin/env bash

CURL="curl_chrome124"
M3U8_URL="$1"
REFERER="https://embedsports.top/"

declare -A played_segments

while true; do
    playlist=$($CURL -s -H "Referer: $REFERER" -H "Origin: $REFERER" "$M3U8_URL")

    segments=$(echo "$playlist" | grep "\.ts$")

    for segment in $segments; do
        if [[ -z "${played_segments[$segment]}" ]]; then
            played_segments[$segment]=1
            host=$(echo "$M3U8_URL" | sed 's|\(https://[^/]*\).*|\1|')
            $CURL -s -H "Referer: $REFERER" -H "Origin: $REFERER" "${host}${segment}"
        fi
    done

    sleep 2
done
