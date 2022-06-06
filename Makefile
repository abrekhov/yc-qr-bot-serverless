include .env

create:
	yc serverless container create --name $(SERVERLESS_CONTAINER_NAME)
	yc serverless container allow-unauthenticated-invoke --name  $(SERVERLESS_CONTAINER_NAME)

create_gw_spec:
	$(shell sed "s/SERVERLESS_CONTAINER_ID/${SERVERLESS_CONTAINER_ID}/;s/SERVICE_ACCOUNT_ID/${SERVICE_ACCOUNT_ID}/" api-gw.yaml.example > api-gw.yaml)
create_gw: create_gw_spec
	yc serverless api-gateway create --name $(SERVERLESS_CONTAINER_NAME) --spec api-gw.yaml
gitlab_env_delete:
	curl   https://$(GITLAB_HOST)/api/v4/projects/$(GITLAB_PROJECT_NAME)/variables -H "private-token: $(GITLAB_TOKEN)" -X GET | jq -rMc '.[] | .key' | xargs -L1 bash -c 'curl   https://$(GITLAB_HOST)/api/v4/projects/$(GITLAB_PROJECT_NAME)/variables/$$0 -H "private-token: $(GITLAB_TOKEN)" -X DELETE '
gitlab_env_push: gitlab_env_delete
	cat .env | awk -F'=' '{print $$1" "$$2}' | xargs  -L1 bash -c ' curl   https://$(GITLAB_HOST)/api/v4/projects/$(GITLAB_PROJECT_NAME)/variables -H "private-token: $(GITLAB_TOKEN)" -X POST --form "key=$$0" --form "value=$$1" '
gitlab_env_pull:
	rm .env.test
	curl   https://$(GITLAB_HOST)/api/v4/projects/$(GITLAB_PROJECT_NAME)/variables -H "private-token: $(GITLAB_TOKEN)" -X GET | jq -rMc '.[] | .key+" "+.value' | xargs -L1 bash -c 'echo $$0=$$1 >> .env.test'
webhook_info:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/getWebhookInfo"

webhook_delete:
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/deleteWebhook"

webhook_create: webhook_delete
	curl --request POST --url "https://api.telegram.org/bot$(TELEGRAM_APITOKEN)/setWebhook" --header 'content-type: application/json' --data "{\"url\": \"$(SERVERLESS_APIGW_URL)\"}"

build: webhook_create
	docker build -t cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME) .

push: build
	docker push cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME)

deploy: push
	$(shell sed 's/=.*/=/' .env > .env.example)
	yc serverless container revision deploy --container-name $(SERVERLESS_CONTAINER_NAME) --image 'cr.yandex/$(YC_IMAGE_REGISTRY_ID)/$(SERVERLESS_CONTAINER_NAME):latest' --service-account-id $(SERVICE_ACCOUNT_ID)  --environment='$(shell tr '\n' ',' < .env)' --core-fraction 5 --execution-timeout $(SERVERLESS_CONTAINER_EXEC_TIMEOUT)

all: deploy