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

COPY requirements.txt ./

# cache pip downloads in BuildKit’s layer cache
RUN --mount=type=cache,id=pip-cache,target=/root/.cache/pip \
    pip install --upgrade pip && \
    pip install -r requirements.txt --extra-index-url https://buf.build/gen/python

# Copy application code
COPY . .

ARG TARGETOS
ARG TARGETARCH

ARG GRPC_CURL_VERSION='1.9.3'
ADD https://github.com/fullstorydev/grpcurl/releases/download/v${GRPC_CURL_VERSION}/grpcurl_${GRPC_CURL_VERSION}_${TARGETOS}_${TARGETARCH}.tar.gz /tmp/grpcurl.tar.gz
RUN tar -xzf /tmp/grpcurl.tar.gz -C /tmp && \
    mv /tmp/grpcurl /usr/local/bin/grpcurl && \
    chmod +x /usr/local/bin/grpcurl && \
    rm /tmp/grpcurl.tar.gz


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

COPY --from=builder /usr/local/bin/grpcurl /usr/local/bin/grpcurl
RUN chmod +x /usr/local/bin/grpcurl

COPY healthcheck.sh /usr/local/bin/healthcheck.sh
RUN chmod +x /usr/local/bin/healthcheck.sh

HEALTHCHECK \
  --interval=30s \
  --start-interval=2s \
  --timeout=3s \
  --start-period=5s \
  --retries=5 \
  CMD /usr/local/bin/healthcheck.sh

# Default command
CMD ["python", "main.py"]
