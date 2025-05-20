import io

import aiohttp
from PIL import Image


class ImageDownloader:
    def __init__(self):
        self._session = aiohttp.ClientSession()

    async def download(self, url: str) -> Image.Image:
        async with self._session.get(url) as resp:
            resp.raise_for_status()
            data = await resp.read()
        return Image.open(io.BytesIO(data)).convert("RGB")

    async def close(self):
        if self._session and not self._session.closed:
            await self._session.close()
