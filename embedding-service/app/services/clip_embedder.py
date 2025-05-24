import clip
import numpy as np
import torch
from PIL import Image
from sentence_transformers import SentenceTransformer

from app.services.compile import try_compile_model
from app.services.device import DEVICE


def normalize(vec: np.ndarray) -> np.ndarray:
    return vec / np.linalg.norm(vec)

class ClipEmbedder:
    """
    Service for loading and using CLIP image and text models.
    Models are loaded immediately upon instantiation.
    """

    def __init__(self, device=DEVICE):
        self.device = device

        # Load and compile CLIP image model
        img_model, preprocess = clip.load("ViT-B/32", device=self.device)
        img_model.eval()
        img_model = try_compile_model(img_model)

        # Load and compile text model
        txt_model = SentenceTransformer(
            "sentence-transformers/clip-ViT-B-32-multilingual-v1",
            device=self.device,
        )
        txt_model = try_compile_model(txt_model)

        # Store models and preprocess function
        self.img_model = img_model
        self.preprocess = preprocess
        self.txt_model = txt_model

    def embed_image(self, pil: Image.Image) -> np.ndarray:
        """Embed and normalize an image using the CLIP model."""
        tensor = self.preprocess(pil).unsqueeze(0).to(self.device)
        with torch.no_grad():
            embedding = self.img_model.encode_image(tensor)
        vec = embedding.cpu().numpy().flatten()
        return vec / np.linalg.norm(vec)  # normalized
    
    def embed_text(self, text: str) -> np.ndarray:
        """Embed and normalize text using the CLIP model."""

        # or use text model
        # vec = self.txt_model.encode([text])[0]

        text_tokenized = clip.tokenize([text]).to(DEVICE)
        with torch.no_grad():
            text_features = self.img_model.encode_text(text_tokenized)
        vec = text_features.detach().cpu().numpy().flatten()
        return vec / np.linalg.norm(vec)  # normalized