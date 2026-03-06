#!/usr/bin/env bash

CURL="curl_chrome124"
M3U8_URL="$1"
REFERER="${2:-https://embedsports.top/}"
MEDIA_PLAYLIST_URL=""

declare -A played_segments

# Resolve nested playlists to find the media playlist
resolve_playlist() {
    local url="$1"
    local playlist
    local first_line
    
    playlist=$($CURL -s -H "Referer: $REFERER" -H "Origin: $REFERER" "$url")
    first_line=$(echo "$playlist" | grep -v '^#' | grep -v '^$' | head -1)
    
    if [[ "$first_line" == *.m3u8 ]]; then
        # It's a nested playlist reference
        if [[ "$first_line" == http* ]]; then
            resolve_playlist "$first_line"
        else
            local base
            base=$(echo "$url" | sed 's|[^/]*$||')
            resolve_playlist "${base}${first_line}"
        fi
    else
        # Found the media playlist
        MEDIA_PLAYLIST_URL="$url"
    fi
}

resolve_playlist "$M3U8_URL"

while true; do
    playlist=$($CURL -s -H "Referer: $REFERER" -H "Origin: $REFERER" "$MEDIA_PLAYLIST_URL")
    
    segments=$(echo "$playlist" | grep -v '^#' | grep -v '^$')
    
    for segment in $segments; do
        if [[ -z "${played_segments[$segment]}" ]]; then
            played_segments[$segment]=1
            
            if [[ "$segment" == http* ]]; then
                segment_url="$segment"
            elif [[ "$segment" == /* ]]; then
                # Absolute path from root - prepend host only
                host=$(echo "$MEDIA_PLAYLIST_URL" | sed 's|\(https://[^/]*\).*|\1|')
                segment_url="${host}${segment}"
            else
                # Relative path - prepend base directory
                base=$(echo "$MEDIA_PLAYLIST_URL" | sed 's|[^/]*$||')
                segment_url="${base}${segment}"
            fi
            
            if [[ "$segment_url" == *"vercel-storage.com"* ]]; then
                curl -s "$segment_url"
            else
                $CURL -s -H "Referer: $REFERER" -H "Origin: $REFERER" "$segment_url"
            fi
        fi
    done
    
    sleep 2
done
