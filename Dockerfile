FROM python:3.7-alpine

WORKDIR /app
COPY Python/requirements.txt .

RUN pip3 install -r requirements.txt

COPY Python/analyticsrelationships.py .

ENTRYPOINT ["python", "analyticsrelationships.py", "-u"]