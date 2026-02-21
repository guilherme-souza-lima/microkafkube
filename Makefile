.PHONY: m1 m2 m3 messaging help

# Roda o Microserviço 1 (Porta App: 6441 | Porta DB: 5441)
m1:
	@echo "Iniciando MICRO_UM na porta 6441..."
	cd m1 && go run cmd/main.go

# Roda o Microserviço 2 (Porta App: 6442 | Porta DB: 5442)
m2:
	@echo "Iniciando MICRO_DOIS na porta 6442..."
	cd m2 && go run cmd/main.go

# Roda o Microserviço 3 (Porta App: 6443 | Porta DB: 5443)
m3:
	@echo "Iniciando MICRO_TRES na porta 6443..."
	cd m3 && go run cmd/main.go

status:
	docker ps

help:
	@echo "Comandos disponíveis:"
	@echo "  make messaging  - Sobe Kafka e UI"
	@echo "  make m1         - Roda API do Microserviço 1"
	@echo "  make m2         - Roda API do Microserviço 2"
	@echo "  make m3         - Roda API do Microserviço 3"

network:
	docker network create estudo-shared-net
