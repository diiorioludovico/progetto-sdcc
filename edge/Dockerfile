FROM python:3.12.9-slim

ENV PYTHONUNBUFFERED=1

WORKDIR /app

# Copia requirements.txt nella root del container
COPY requirements.txt .

# Installa le dipendenze
RUN pip install --no-cache-dir -r requirements.txt

# Copia la cartella con il codice
COPY python/ ./python
COPY config.json .

# Imposta la directory di lavoro sul codice Python
WORKDIR /app/python

# Avvia lo script edge.py
CMD ["python", "edge.py"]
