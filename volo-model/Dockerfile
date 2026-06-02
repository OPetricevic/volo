FROM python:3.11-slim

WORKDIR /workspace

# System deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install PyTorch CPU-only first (avoids downloading 900MB+ of CUDA libs)
RUN pip install --no-cache-dir torch --index-url https://download.pytorch.org/whl/cpu

# Python deps (torch already satisfied, skipped)
COPY requirements.txt .
RUN pip install --no-cache-dir --timeout=300 -r requirements.txt

# Jupyter config
RUN pip install --no-cache-dir jupyterlab ipywidgets

EXPOSE 8888

CMD ["jupyter", "lab", "--ip=0.0.0.0", "--port=8888", "--no-browser", "--allow-root"]
