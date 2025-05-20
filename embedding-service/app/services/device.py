import torch

if torch.backends.mps.is_built():
    _device = "mps"
elif torch.cuda.is_available():
    _device = "cuda"
else:
    _device = "cpu"

print(f"Using device {_device}")
DEVICE = torch.device(_device)
