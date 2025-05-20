from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.grpc import embedding
from app.api import caption
from app.services import captioner, clip_embedder


@asynccontextmanager
async def lifespan(app: FastAPI):
    print("Preloading captioner...")
    captioner.preload()  # BLIP caption generator (downloads weights if needed)
    print("Preloading clip_embedder...")
    clip_embedder.preload()  # CLIP image & multilingual text encoders

    print("Starting grpc server...")
    grpc_server = embedding.start_grpc_server()  # Start gRPC server in background thread
    yield
    grpc_server.stop(0)  # startup complete – serve requests


app = FastAPI(title="Media AI API", version="0.1.0", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(caption.router, prefix="/caption", tags=["Caption"])


@app.get("/")
def health():
    return {"status": "ok"}
