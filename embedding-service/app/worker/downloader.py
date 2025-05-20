import io

import aiohttp
from PIL import Image

_http_session = None


async def get_http_session():
    global _http_session
    if _http_session is None:
        _http_session = aiohttp.ClientSession()
    return _http_session


async def download_image(url: str) -> Image.Image:
    session = await get_http_session()
    async with session.get(url) as resp:
        resp.raise_for_status()
        data = await resp.read()
    return Image.open(io.BytesIO(data)).convert("RGB")
