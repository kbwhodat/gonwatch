from selenium_driverless import webdriver
from selenium_driverless.types.by import By
import json
from pahe import AnimePahe
import os
import requests
import asyncio
import re
import argparse
from langdetect import detect
import urllib.request
import urllib.parse
import warnings

USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"

warnings.filterwarnings("ignore", category=UserWarning)

async def clean_vtt(content):
    text = re.sub(r'\d{2}:\d{2}:\d{2}\.\d{3} --> .*', '', content)
    text = re.sub(r'WEBVTT.*', '', text)
    text = re.sub(r'\[.*?\]', '', text)  # Remove sound effects
    text = re.sub(r'\n+', '\n', text)    # Collapse multiple newlines
    text = re.sub(r'<.*?>', '', text)    # Remove HTML tags
    return text.strip()

def decrypt_url(provider_id: str) -> str:
    decrypted = ""
    for hex_value in [provider_id[i : i + 2] for i in range(0, len(provider_id), 2)]:
        dec = int(hex_value, 16)
        xor = dec ^ 56
        oct_value = oct(xor)[2:].zfill(3)
        decrypted += chr(int(oct_value, 8))
    return decrypted

async def get_video_url(show_id: str, episode: str, lang: str = "sub") -> str:
    query = """
    query ($showId: String!, $translationType: VaildTranslationTypeEnumType!, $episodeString: String!) {
        episode(
            showId: $showId
            translationType: $translationType
            episodeString: $episodeString
        ) {
            episodeString,
            sourceUrls
        }
    }
    """
    url = bytes(b ^ 0x42 for b in bytes.fromhex("2a36363231786d6d23322b6c232e2e232c2b2f276c26233b6d23322b")).decode()
    response = requests.get(
        url,
        params={
            "variables": f'{{"showId":"{show_id}","translationType":"{lang}","episodeString":"{episode}"}}',
            "query": query,
        },
        headers={"Referer": bytes([104, 116, 116, 112, 115, 58, 47, 47, 97, 108, 108, 109, 97, 110, 103, 97, 46, 116, 111, 47]).decode()},
    )

    data = response.json()
    sources = data["data"]["episode"]["sourceUrls"]

    for source in sources:
        # if source["sourceName"] == "Yt-mp4":
        #     encrypted = source["sourceUrl"].replace("--", "")
        #     return decrypt_url(encrypted)
        if source["sourceName"] == "S-mp4":
            encrypted = source["sourceUrl"].replace("--", "")
            decrypted = decrypt_url(encrypted).replace("clock", "clock.json")
            url = f"https://{bytes([97, 108, 108, 97, 110, 105, 109, 101, 46, 100, 97, 121]).decode()}{decrypted}"
            response = requests.get(url, headers={"Referer": bytes([104, 116, 116, 112, 115, 58, 47, 47, 97, 108, 108, 109, 97, 110, 103, 97, 46, 116, 111, 47]).decode(), "User-Agent": USER_AGENT})
            if response.text and response.text != "error":
                try:
                    result = response.json()
                    if result and "links" in result:
                        for link in result["links"]:
                            return link['link']
                except:
                    return ""
    return ""

m3u8_found = asyncio.Event()

async def main() -> str:

    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument("--start-maximized")
    chrome_options.add_argument("--window-size=1920,1080")
    chrome_options.add_argument("--headless=new")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--mute-audio")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument("--disable-features=DisableLoadExtensionCommandLineSwitch")
    chrome_options.add_experimental_option("prefs", {
        "profile.default_content_setting_values.popups": 0,  # Allow pop-ups
        "profile.default_content_setting_values.cookies": 1,
        "profile.cookie_controls_mode": 0,
        "profile.block_third_party_cookies": False,
    })

    parser = argparse.ArgumentParser()
    parser.add_argument("content")
    parser.add_argument("id")
    parser.add_argument("season", nargs='?')
    parser.add_argument("episode", nargs='?')
    parser.add_argument("title", nargs='?')
    parser.add_argument("anilist_id", nargs='?')
    parser.add_argument("anime_episode", nargs='?')
    parser.add_argument("stream_url", nargs='?')
    parser.add_argument("anime_title", nargs='?')
    args = parser.parse_args()

    async with webdriver.Chrome(options=chrome_options) as driver:

        await driver.execute_cdp_cmd("Network.enable", {})

        urls = []
        subtitles = []
        langs = []
        async def on_response(event):
            try:
                url = event.get("response").get("url")
                m3u8 = re.compile(".*m3u8")
                vtt = re.compile(".*vtt")

                # print(url)


                if url not in urls:

                    if "test.vidify.top" in url:
                        urls.append(url)

                    if ".m3u8" in url:
                        if ".gif" not in url:
                            if "vidnest.fun" in url:
                                decoded_url = urllib.parse.unquote(urllib.parse.urlparse(url).query.split("url=")[1])
                                m = m3u8.match(decoded_url)
                                urls.append(m.group())
                                m3u8_found.set()
                            elif "mono.ts.m3u8" in url:
                                urls.append(url)
                            elif "strmd.top" in url or "gg.poocloud.in" in url or "vdcast.live" in url:
                                urls.append(url)
                            else:
                                urls.append(url)
                    elif "cf-master" in url:
                        urls.append(url)

                if ".vtt" in url:
                    # I can add more languages here as needed if end users want. But right now muliple vtt files are being loaded to mpv slowing it down so I'm being more selective with the languages.
                    if "es" not in langs or "en" not in langs:
                        v = vtt.match(url)
                        response = urllib.request.urlopen(v.group())
                        data = response.read()
                        cleanvtt = await clean_vtt(data.decode("utf-8"))
                        lang = detect(cleanvtt)
                        if lang == "en" or lang == "es":
                            langs.append(lang)
                            subtitles.append(v.group())

            except ValueError as e:
                return json.dumps([])

        if args.content == "stream":
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            await driver.get(args.stream_url, wait_load=True)
            await driver.sleep(4)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)
            # urls = [args.stream_url]

            result =  json.dumps({"urls": urls, "subtitles": subtitles})
            return result


        if args.content == "anime":
            location = "/tmp/cookies.json"

            if not os.path.isfile(location):
                await driver.add_cdp_listener("Network.responseReceived", on_response)
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 97, 110, 105, 109, 101, 112, 97, 104, 101, 46, 115, 105, 47])
                await driver.get(url_bytes.decode(), wait_load=True)
                await driver.sleep(0.5)
                await driver.wait_for_cdp("Page.domContentEventFired", timeout=15)
                await driver.sleep(5)

                await driver.remove_cdp_listener("Network.responseReceived", on_response)

                all_cookies = await driver.get_cookies()
                important_names = {"laravel_session", "XSRF-TOKEN", "ddg_last_challenge"}
                cookie_list = [
                    c
                    for c in all_cookies
                    if c["name"].startswith("__ddg") or c["name"] in important_names
                ]

                with open(location, "w+") as filehandler:
                    json.dump(cookie_list, filehandler, indent=4)

            pahe = AnimePahe()
            anime_session = await pahe.search(args.anime_title)
            if anime_session:
                episodes = await pahe.get_episodes(anime_session)
                episode_num = int(args.episode)
                target_episode = episodes[episode_num - 1] if 0 <= (episode_num - 1) < len(episodes) else None

                if target_episode:
                    episode_session = target_episode['session']

                    sources = await pahe.get_sources(anime_session, episode_session)

                    if sources:
                        m3u8_url = await pahe.resolve_kwik_with_node(sources[0]['url'])
                        urls.append(m3u8_url)
            else:
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 104, 101, 97, 118, 101, 110, 115, 99, 97, 112, 101, 46, 118, 101, 114, 99, 101, 108, 46, 97, 112, 112, 47, 97, 112, 105, 47, 97, 110, 105, 109, 101, 47, 115, 101, 97, 114, 99, 104, 47])
                url = requests.get(url_bytes.decode() + f'{args.title}/sub/{args.episode}').json().get("direct")

                if len(url) == 0:
                    url = await get_video_url(args.title, args.episode)
                urls.append(url)

            result =  json.dumps({"urls": urls, "subtitles": subtitles})
            return result


        await driver.add_cdp_listener("Network.responseReceived", on_response)
        if args.content == "tv":
            url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 108, 105, 110, 107, 46, 112, 114, 111, 47, 116, 118, 47])
            await driver.get(url_bytes.decode() + f'{args.id}/{args.season}/{args.episode}', wait_load=True)
        else:
            url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 108, 105, 110, 107, 46, 112, 114, 111, 47, 109, 111, 118, 105, 101, 47])
            await driver.get(url_bytes.decode() + f'{args.id}', wait_load=True)
        await driver.sleep(3)
        await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 115, 114, 99, 46, 99, 99, 47, 118, 51, 47, 101, 109, 98, 101, 100, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
            else:
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 115, 114, 99, 46, 99, 99, 47, 118, 51, 47, 101, 109, 98, 101, 100, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}', wait_load=True)

            await driver.sleep(2.5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # 2nd Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 110, 101, 115, 116, 46, 102, 117, 110, 47, 116, 118, 47])
                await driver.get(url_bytes.decode() + f'{args.id}/{args.season}/{args.episode}', wait_load=True)
            else:
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 110, 101, 115, 116, 46, 102, 117, 110, 47, 109, 111, 118, 105, 101, 47])
                await driver.get(url_bytes.decode() + f'{args.id}/{args.season}/{args.episode}', wait_load=True)
            await driver.sleep(5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # 3rd Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 115, 114, 99, 45, 101, 109, 98, 101, 100, 46, 114, 117, 47, 101, 109, 98, 101, 100, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
                try:
                    player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=5)
                    url = await player_iframe.get_attribute("src")

                    await driver.get(url, wait_load=True)
                    play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=10)
                    await driver.sleep(1)
                    await driver.execute_script("arguments[0].click();", play_trigger)

                except:
                    result =  json.dumps({"urls": urls, "subtitles": subtitles})
            else:
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 118, 105, 100, 115, 114, 99, 45, 101, 109, 98, 101, 100, 46, 114, 117, 47, 101, 109, 98, 101, 100, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}', wait_load=True)
                try:
                    player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=3)
                    url = await player_iframe.get_attribute("src")

                    await driver.get(url, wait_load=True)
                    play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=3)
                    await driver.sleep(1)
                    await driver.execute_script("arguments[0].click();", play_trigger)
                except Exception as e:
                    result =  json.dumps({"urls": urls, "subtitles": subtitles})

            await driver.sleep(2.5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # 4th fallback option
        if len(urls) == 0:
            if args.content == "tv":
                await driver.add_cdp_listener("Network.responseReceived", on_response)
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 104, 101, 120, 97, 46, 115, 117, 47, 119, 97, 116, 99, 104, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
                await driver.sleep(10)
                await driver.remove_cdp_listener("Network.responseReceived", on_response)
            else:
                await driver.add_cdp_listener("Network.responseReceived", on_response)
                url_bytes = bytes([104, 116, 116, 112, 115, 58, 47, 47, 104, 101, 120, 97, 46, 115, 117, 47, 119, 97, 116, 99, 104, 47])
                await driver.get(url_bytes.decode() + f'{args.content}/{args.id}', wait_load=True)
                await driver.sleep(10)
                await driver.remove_cdp_listener("Network.responseReceived", on_response)

        result =  json.dumps({"urls": urls, "subtitles": subtitles})
        return result

# I need this print statement if I want the return value to be sent back to golang
print(asyncio.run(main()))
