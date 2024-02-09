ifeq ($(OS),Windows_NT)
    # Windows
    DOWNLOAD_HASURA = @powershell -Command "if (-not (Test-Path "..\hasura\hasura.exe")) { $$ProgressPreference = 'SilentlyContinue'; Invoke-WebRequest -Uri "https://github.com/hasura/graphql-engine/releases/download/v2.37.0/cli-hasura-windows-amd64.exe" -OutFile "..\hasura\hasura.exe"}"
    HASURA_METADATA = @powershell -Command "cd ../hasura; .\hasura.exe metadata apply"
    HASURA_MIGRATE = @powershell -Command "cd ../hasura; .\hasura.exe migrate apply"
    HASURA_SEED = @powershell -Command "cd ../hasura; .\hasura.exe seed apply"
    REMOVE_HASURA = @powershell -Command "Remove-Item -Path "..\hasura\hasura.exe" "
else
    # Linux/Unix
    DOWNLOAD_HASURA = wget "https://github.com/hasura/graphql-engine/releases/download/v2.37.0/cli-hasura-linux-amd64" -O "hasura"
    HASURA_METADATA = cd ../hasura; ./hasura.exe metadata apply
    HASURA_MIGRATE = cd ../hasura; ./hasura.exe migrate apply
    HASURA_SEED = cd ../hasura; ./hasura.exe seed apply
    REMOVE_HASURA = rm -f ../hasura/hasura
endif

.PHONY: all

all:
	helm repo add hasura https://hasura.github.io/helm-charts
	helm repo update
	helm upgrade --install hasura -f ./helm/values.yaml hasura/graphql-engine
	$(DOWNLOAD_HASURA)
	$(HASURA_METADATA)
	$(HASURA_MIGRATE)
	$(HASURA_SEED)
	kubectl apply -f hasura-jwt

ingress:
	helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --create-namespace

clean:
	helm uninstall hasura
	kubectl delete -f hasura-jwt
	$(REMOVE_HASURA)