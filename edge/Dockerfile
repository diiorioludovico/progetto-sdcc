# Stage 1: Build
FROM python:3.9-slim AS builder

# Installa le dipendenze di sistema necessarie per compilare grpcio-tools
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copia il file requirements.txt e installa le dipendenze
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Stage 2: Runtime
FROM python:3.9-slim

# Imposta la directory di lavoro nel container
WORKDIR /app

# Copia le dipendenze Python installate dallo stage precedente
COPY --from=builder /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# Copia tutti i file dell'applicazione
COPY . .

WORKDIR /app/python

# Crea la directory __pycache__ se non esiste (per evitare problemi di permessi)
RUN mkdir -p __pycache__

# Comando per avviare l'applicazione
# Modifica questo comando in base al tuo file principale
CMD ["python", "edge.py"]