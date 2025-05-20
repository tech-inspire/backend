import torch, clip, numpy as np
from PIL import Image
from functools import lru_cache
from sentence_transformers import SentenceTransformer

from app.services.device import get_device
from app.services.compile import try_compile

device=get_device()

@lru_cache(maxsize=1)
def _load():
    """Download & cache CLIP (image) and multilingual text models, compiling when possible."""
    img_model, preprocess = clip.load("ViT-B/32", device=device)
    img_model.eval()
    img_model = try_compile(img_model)

    txt_model = SentenceTransformer("sentence-transformers/clip-ViT-B-32-multilingual-v1", device=device)
    try:
        txt_model = try_compile(txt_model)
    except Exception:
        # SentenceTransformer isn't always compatible with compile; ignore failures
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