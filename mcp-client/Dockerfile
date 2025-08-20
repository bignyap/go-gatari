# Stage 1: Builder
FROM python:3.12.11-slim AS builder

WORKDIR /app 

# Install dependencies
COPY requirements.txt .
RUN pip install --upgrade pip setuptools
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Stage 2: Final Image
FROM python:3.12.11-slim

# Install curl & netcat (Debian-based command, not apk)
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy installed packages from builder
COPY --from=builder /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages

# Copy application code
COPY --from=builder /app .

# Expose the port the app runs on
EXPOSE 9000

# Set the command to run the application
CMD ["python", "main.py"]
