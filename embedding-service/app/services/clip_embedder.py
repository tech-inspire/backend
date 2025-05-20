from functools import lru_cache

import clip
import numpy as np
import torch
from PIL import Image
from sentence_transformers import SentenceTransformer

from app.services.compile import try_compile
from app.services.device import get_device

device = get_device()


@lru_cache(maxsize=1)
def _load():
    img_model, preprocess = clip.load("ViT-B/32", device=device)
    img_model.eval()
    img_model = try_compile(img_model)

    txt_model = SentenceTransformer(
        "sentence-transformers/clip-ViT-B-32-multilingual-v1", device=device
    )
    try:
        txt_model = try_compile(txt_model)
    except Exception:
        pass

    return img_model, preprocess, txt_model


def preload():
    _load()


def embed_image(pil: Image.Image) -> np.ndarray:
    img_model, preprocess, _ = _load()
    tensor = preprocess(pil).unsqueeze(0).to(device)
    with torch.no_grad():
        return img_model.encode_image(tensor).cpu().numpy().flatten()


def embed_text(text: str) -> np.ndarray:
    _, _, txt_model = _load()
    return txt_model.encode([text])[0]
