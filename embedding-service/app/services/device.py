import torch

def get_device():
    device = "mps" if torch.backends.mps.is_built() else "cuda" if torch.cuda.is_available() else "cpu"
    print(f'Using device {device}')
    return device