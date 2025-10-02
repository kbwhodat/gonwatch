from selenium_driverless import webdriver
from selenium_driverless.types.by import By
import json
import asyncio
import re
import argparse

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
    args = parser.parse_args()

    async with webdriver.Chrome(options=chrome_options) as driver:

        await driver.execute_cdp_cmd("Network.enable", {})

        urls = []
        subtitles = []
        async def on_response(event):
            m3u8 = re.compile(".*m3u8")
            vtt = re.compile(".*vtt")
            if ".m3u8" in event.get("response").get("url"):
                m = m3u8.match(event.get("response").get("url"))
                if ".gif" not in m.group():
                    urls.append(m.group())
            if ".vtt" in event.get("response").get("url"):
                v = vtt.match(event.get("response").get("url"))
                if "eng" in v.group().lower() or "spa" in v.group().lower() or "fre" in v.group().lower():
                    subtitles.append(v.group())
                else:
                    subtitles.append(v.group())

        await driver.add_cdp_listener("Network.responseReceived", on_response)
        if args.content == "tv":
            await driver.get(f'https://vidsrc.cc/v2/embed/{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
        else:
            await driver.get(f'https://vidsrc.cc/v2/embed/{args.content}/{args.id}', wait_load=True)

        await driver.sleep(2.5)
        await driver.remove_cdp_listener("Network.responseReceived", on_response)

        # Fallback option
        if len(urls) == 0:
            await driver.add_cdp_listener("Network.responseReceived", on_response)
            if args.content == "tv":
                await driver.get(f'https://vidsrc.xyz/embed/{args.content}/{args.id}/{args.season}/{args.episode}', wait_load=True)
                player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=15)
                url = await player_iframe.get_attribute("src")
                await driver.get(url)
                await driver.sleep(1)
                play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=15)
                await driver.execute_script("arguments[0].click();", play_trigger)
            else:
                await driver.get(f'https://vidsrc.xyz/embed/{args.content}/{args.id}', wait_load=True)
                player_iframe = await driver.find_element(By.ID, "player_iframe", timeout=15)
                url = await player_iframe.get_attribute("src")
                await driver.get(url)
                await driver.sleep(1)
                play_trigger = await driver.find_element(By.CSS_SELECTOR, "i#pl_but, #pl_but, #pl_but_background", timeout=15)
                await driver.execute_script("arguments[0].click();", play_trigger)


            await driver.sleep(2.5)
            await driver.remove_cdp_listener("Network.responseReceived", on_response)

        result =  json.dumps({"urls": urls, "subtitles": subtitles})
        return result

print(asyncio.run(main()))
