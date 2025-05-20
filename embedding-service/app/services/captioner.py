from transformers import pipeline
from PIL import Image
import torch
from functools import lru_cache

from app.services.compile import try_compile
from app.services.device import get_device

device = get_device();



@lru_cache(maxsize=1)
def _get():
    """Instantiate and cache the BLIP captioning pipeline (compiled if possible)."""
    pipe = pipeline("image-to-text", model="Salesforce/blip-image-captioning-base", device=device, torch_dtype=torch.float16)
    pipe.model = try_compile(pipe.model)
    return pipe


def preload():
    """Public helper for startup preload."""
    _get()


_GENERATION_CFG = dict(
    max_new_tokens=20,
    num_beams=5,
    repetition_penalty=1.2,
)

def generate_caption(pil: Image.Image) -> str:
    text = _get()(pil, generate_kwargs=_GENERATION_CFG)[0]["generated_text"]
    return text.strip()