import torch


def try_compile(model: torch.nn.Module):
    if hasattr(torch, "compile"):
        try:
            print("compiling model")
            model = torch.compile(model, mode="reduce-overhead")
            print("compiling completed")
        except Exception:
            pass
    return model
