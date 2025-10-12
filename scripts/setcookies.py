from selenium_driverless import webdriver
from selenium_driverless.types.by import By
import json
import requests
import asyncio
import re
import argparse
from langdetect import detect
import urllib.request
import urllib.parse

async def clean_vtt(content):
    text = re.sub(r'\d{2}:\d{2}:\d{2}\.\d{3} --> .*', '', content)
    text = re.sub(r'WEBVTT.*', '', text)
    text = re.sub(r'\[.*?\]', '', text)  # Remove sound effects
    text = re.sub(r'\n+', '\n', text)    # Collapse multiple newlines
    text = re.sub(r'<.*?>', '', text)    # Remove HTML tags
    return text.strip()

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
                            else:
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


        # try:
        #     if "JP" in args.title:

        #         await driver.add_cdp_listener("Network.responseReceived", on_response)

        #         # print(f"https://vidnest.fun/anime/{args.anilist_id}/{args.anime_episode}/sub")
        #         await driver.get(f'https://vidnest.fun/anime/{args.anilist_id}/{args.anime_episode}/sub', wait_load=True)

        #         try:
        #             await asyncio.wait_for(m3u8_found.wait(), timeout=7)
        #         except asyncio.TimeoutError:
        #             css = await driver.find_element(By.CSS_SELECTOR, ".Container_player__xy_1D > div:nth-child(1) > video:nth-child(1) > source:nth-child(1)", timeout=3)
        #             src = await css.get_attribute("src")
        #             urls.append(src)

        #         await driver.remove_cdp_listener("Network.responseReceived", on_response)
        #         result =  json.dumps({"urls": urls, "subtitles": subtitles})

        #         if len(urls) != 0:
        #             return result


        # except Exception as e:
        #     await driver.add_cdp_listener("Network.responseReceived", on_response)
        #     await driver.get(f'https://player.vidify.top/embed/tv/{args.id}/{args.season}/{args.episode}/?autoplay=true', wait_load=True)
        #     await driver.sleep(8)
        #     await driver.remove_cdp_listener("Network.responseReceived", on_response)

        #     result =  json.dumps({"urls": urls, "subtitles": subtitles})

        #     if len(urls) != 0:
        #         return result
        #     else:
        #         return json.dumps([])
        #
        if args.content == "anime":
            url = requests.get(f'https://heavenscape.vercel.app/api/anime/search/{args.title}/sub/{args.episode}').json().get("direct")
            urls.append(url)

            result =  json.dumps({"urls": urls, "subtitles": subtitles})

            if len(urls) != 0:
                return result
            else:
                return result



        await driver.add_cdp_listener("Network.responseReceived", on_response)
        if args.content == "tv":
            await driver.get(f'https://vidlink.pro/tv/{args.id}/{args.season}/{args.episode}', wait_load=True)
        else:
            await driver.get(f'https://vidlink.pro/movie/{args.id}', wait_load=True)
        await driver.sleep(3)
        await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                await driver.get(f'https://vidsrc.cc/v2/embed/{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
            else:
                await driver.get(f'https://vidsrc.cc/v2/embed/{args.content}/{args.id}', wait_load=True)

            await driver.sleep(2.5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # 2nd Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                await driver.get(f'https://vidsrc.xyz/embed/{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
                try:
                    player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=3)
                    url = await player_iframe.get_attribute("src")

                    await driver.get(url, wait_load=True)
                    play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=3)
                    await driver.execute_script("arguments[0].click();", play_trigger)
                except:
                    return json.dumps([])
            else:
                await driver.get(f'https://vidsrc.xyz/embed/{args.content}/{args.id}', wait_load=True)
                try:
                    player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=3)
                    url = await player_iframe.get_attribute("src")

                    await driver.get(url, wait_load=True)
                    play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=3)
                    await driver.execute_script("arguments[0].click();", play_trigger)
                except Exception as e:
                    print(e)
                    return json.dumps([])


            await driver.sleep(2.5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # 3rd fallback option
        if len(urls) == 0:
            if args.content == "tv":
                await driver.add_cdp_listener("Network.responseReceived", on_response)
                await driver.get(f'https://hexa.watch/watch/{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
                await driver.sleep(10)
                await driver.remove_cdp_listener("Network.responseReceived", on_response)
            else:
                await driver.add_cdp_listener("Network.responseReceived", on_response)
                await driver.get(f'https://hexa.watch/watch/{args.content}/{args.id}', wait_load=True)
                await driver.sleep(10)
                await driver.remove_cdp_listener("Network.responseReceived", on_response)

        result =  json.dumps({"urls": urls, "subtitles": subtitles})
        return result

# I need this print statement if I want the return value to be sent back to golang
print(asyncio.run(main()))
