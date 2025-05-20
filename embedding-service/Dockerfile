# Stage 1: Builder
FROM python:3.12-slim AS builder

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create and activate virtual environment
WORKDIR /src
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy and install Python dependencies
COPY requirements.txt ./
RUN pip install --upgrade pip && pip install -r requirements.txt --extra-index-url https://buf.build/gen/python

# Copy application code
COPY . .


# Stage 2: Runtime
FROM python:3.12-slim AS runtime

# Set environment variables for runtime
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    HF_HUB_DISABLE_PROGRESS_BARS=0

# Copy virtual environment and app from builder
COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /src /src

WORKDIR /src
ENV PATH="/opt/venv/bin:$PATH"

# Expose ports
EXPOSE 8000
EXPOSE 50051

# Default command
CMD ["python", "main.py"]
