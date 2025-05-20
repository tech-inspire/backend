from fastapi import APIRouter, UploadFile, File, HTTPException
from fastapi.responses import JSONResponse
from PIL import Image
import io
from fastapi import APIRouter, Request
from fastapi.responses import JSONResponse
from PIL import Image
from io import BytesIO

from app.services.captioner import generate_caption

router = APIRouter()


@router.post("/image-bytes")
async def get_image_embedding_bytes(request: Request):
    body = await request.body()
    image = Image.open(io.BytesIO(body)).convert("RGB")
    caption = generate_caption(image)
    return JSONResponse({"caption": caption})
