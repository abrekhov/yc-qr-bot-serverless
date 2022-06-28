# YC-QR-BOT

<https://t.me/qr_it_bot>

Bot for qr encoding/decoding any content.

## Prerequires

- YC CLI
- golang
- docker desktop

## Quickstart

```bash
cp .env.example .env

vi .env # edit env values: set name, telegram bot token, yc_image_registry_id, service account id, ydb document endpoint, static eky credentials

make create # create serverless container. Copy ID and container URL. Paste back to .env

# create serversless api gateway
make create_gw # Copy API GW URL and paste it to .env

make deploy # build, push and deploy serverless containers.
```
