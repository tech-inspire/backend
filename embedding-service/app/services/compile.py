import torch


def try_compile_model(model: torch.nn.Module):
    if hasattr(torch, "compile"):
        try:
            print("Compiling model")
            model = torch.compile(model, mode="reduce-overhead")
            print("Model compiling completed")
        except Exception:
            pass
    return model
