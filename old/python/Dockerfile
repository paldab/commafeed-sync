FROM python:3.11-slim

WORKDIR /app
COPY seed-feeds.py .
COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "seed-feeds.py"]

