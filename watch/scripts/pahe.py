
import asyncio
import re
import tempfile
import os
import json
from urllib.parse import unquote, quote
from bs4 import BeautifulSoup
from typing import Optional
import tls_client


# -------------------- Constants --------------------
USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"


# Credit goas to Pal-droid. Made some changes to his code but still great work for writing the logic. Not as easy as it looks...
class AnimePahe:
    def __init__(self):
        self.base = b"\x68\x74\x74\x70\x73\x3a\x2f\x2f\x61\x6e\x69\x6d\x65\x70\x61\x68\x65\x2e\x73\x69".decode()
        self.headers = {
            "User-Agent": USER_AGENT,
            "Referer": b"\x68\x74\x74\x70\x73\x3a\x2f\x2f\x61\x6e\x69\x6d\x65\x70\x61\x68\x65\x2e\x73\x69\x2f".decode(),
        }
        self.session = tls_client.Session(client_identifier="chrome_120")

        try:
            with open("/tmp/cookies.json", "r") as f:
                cookies = json.load(f)
                for cookie in cookies:
                    self.session.cookies.set(
                        cookie["name"],
                        cookie["value"],
                        domain=cookie.get("domain", "animepahe.si"),
                        path=cookie.get("path", "/"),
                    )
        except FileNotFoundError:
            print("Warning: /tmp/cookies.json not found")

    async def get(self, url: str):
        """Run tls-client GET asynchronously"""

        def _req():
            r = self.session.get(url, headers=self.headers)
            return r

        return await asyncio.to_thread(_req)

    async def search(self, query: str):
        query = quote(query)
        url = f"{self.base}/api?m=search&q={query}"
        r = await self.get(url)
        data = r.json()
        results = []
        for a in data.get("data", []):
            results.append(
                {
                    "id": a["id"],
                    "title": a["title"],
                    "url": f"{self.base}/anime/{a['session']}",
                    "year": a.get("year"),
                    "poster": a.get("poster"),
                    "type": a.get("type"),
                    "session": a.get("session"),
                }
            )

        for i in results:
            if i.get("title") == unquote(query):
                return i.get("session")

    async def get_episodes(self, anime_session: str):
        html = (await self.get(f"{self.base}/anime/{anime_session}")).text
        soup = BeautifulSoup(html, "html.parser")
        meta = soup.find("meta", {"property": "og:url"})
        if not meta:
            raise Exception("Could not find session ID in meta tag")
        temp_id = meta["content"].split("/")[-1]

        first_page_json = (
            await self.get(
                f"{self.base}/api?m=release&id={temp_id}&sort=episode_asc&page=1"
            )
        ).json()
        episodes = first_page_json.get("data", [])
        last_page = first_page_json.get("last_page", 1)

        async def fetch_page(p):
            return (
                (
                    await self.get(
                        f"{self.base}/api?m=release&id={temp_id}&sort=episode_asc&page={p}"
                    )
                )
                .json()
                .get("data", [])
            )

        tasks = [fetch_page(p) for p in range(2, last_page + 1)]
        for pages in await asyncio.gather(*tasks):
            episodes.extend(pages)

        return [
            {
                "id": e["id"],
                "number": e["episode"],
                "title": e.get("title") or f"Episode {e['episode']}",
                "session": e["session"],
            }
            for e in episodes
        ]

    async def get_sources(self, anime_session: str, episode_session: str):
        html = (
            await self.get(f"{self.base}/play/{anime_session}/{episode_session}")
        ).text

        buttons = re.findall(
            r'<button[^>]+data-src="([^"]+)"[^>]+data-fansub="([^"]+)"[^>]+data-resolution="([^"]+)"[^>]+data-audio="([^"]+)"[^>]*>',
            html,
        )

        sources = []
        for src, fansub, resolution, audio in buttons:
            if src.startswith("https://kwik."):
                sources.append(
                    {
                        "url": src,
                        "quality": f"{resolution}p",
                        "fansub": fansub,
                        "audio": audio,
                    }
                )

        if not sources:
            kwik_links = re.findall(r"https:\/\/kwik\.(si|cx|link)\/e\/\w+", html)
            sources = [
                {"url": link, "quality": None, "fansub": None, "audio": None}
                for link in kwik_links
            ]

        unique_sources = list({s["url"]: s for s in sources}.values())

        def sort_key(s):
            try:
                return int(s["quality"].replace("p", "")) if s["quality"] else 0
            except Exception:
                return 0

        unique_sources.sort(key=sort_key, reverse=True)

        if not unique_sources:
            raise Exception("No kwik links found on play page")

        return unique_sources

    async def resolve_kwik_with_node(
        self, kwik_url: str, node_bin: str = "node"
    ) -> str:
        """Use tls-client and Node.js to extract .m3u8"""
        resp = await self.get(kwik_url)

        html = resp.text
        m3u8_direct = re.search(r"https?://[^'\"\s<>]+\.m3u8", html)
        if m3u8_direct:
            return m3u8_direct.group(0)

        scripts = re.findall(r"(<script[^>]*>[\s\S]*?</script>)", html, re.IGNORECASE)
        script_block, largest_eval_script, max_len = None, None, 0
        for s in scripts:
            if "eval(" in s:
                if "Plyr" in s or ".m3u8" in s or "source" in s or "uwu" in s:
                    script_block = s
                    break
                if len(s) > max_len:
                    max_len = len(s)
                    largest_eval_script = s
        if not script_block:
            script_block = largest_eval_script
        if not script_block:
            raise Exception("No candidate <script> block found")

        inner_js = re.sub(
            r"^<script[^>]*>", "", script_block, flags=re.IGNORECASE
        ).strip()
        inner_js = re.sub(r"</script>$", "", inner_js, flags=re.IGNORECASE).strip()

        wrapper = r"""
globalThis.window = { location: {} };
globalThis.document = { cookie: '' };
globalThis.navigator = { userAgent: 'mozilla' };
const __captured = [];
const origLog = console.log;
console.log = (...args)=>{__captured.push(args.join(' '));origLog(...args);};
(function(){
  const origEval = eval;
  eval = (x)=>{__captured.push('[EVAL]' + x);return origEval(x);};
})();
"""
        final_js = (
            wrapper
            + "\n"
            + inner_js
            + "\n"
            + (
                "setTimeout(()=>{for(const c of __captured){console.log('__CAPTURED__START__');"
                "console.log(c);console.log('__CAPTURED__END__');}process.exit(0)},300);"
            )
        )

        with tempfile.NamedTemporaryFile(
            "w", suffix=".js", delete=False, encoding="utf-8"
        ) as tf:
            tmp_path = tf.name
            tf.write(final_js)
            tf.flush()

        try:
            proc = await asyncio.create_subprocess_exec(
                node_bin,
                tmp_path,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
            )
            stdout, stderr = await proc.communicate()
        finally:
            os.unlink(tmp_path)

        out = stdout.decode(errors="ignore") + stderr.decode(errors="ignore")
        m = re.search(r"https?://[^'\"\s]+\.m3u8[^\s'\"\)]*", out)
        if m:
            return m.group(0)

        raise Exception(f"Could not resolve .m3u8. Node output:\n{out[:1000]}")
