SWAG_OUT   = cmd/server/docs
SWAG_MAIN  = cmd/server/main.go

.PHONY: swag-install docker-install swagger

swag:swagger
	@echo "Проверяем наличие swag..."
	@if command -v swag >/dev/null 2>&1; then \
	  echo "swag уже установлен: $$(swag --version)"; \
	else \
	  echo "swag не найден. Ставим через go install..."; \
	  go install github.com/swaggo/swag/cmd/swag@latest; \
	  echo "Установили swag: $$(swag --version)"; \
	fi

docker-install:
	@echo "Проверяем наличие Docker..."
	@if command -v docker >/dev/null 2>&1; then \
	  echo "Docker уже установлен: $$(docker --version)"; \
	else \
	  echo "Docker не найден. Устанавливаем…"; \
	  if [ -r /etc/os-release ]; then . /etc/os-release; else echo "Не удалось определить дистрибутив"; exit 1; fi; \
	  case "$$ID" in \
	    ubuntu|debian) sudo apt update && sudo apt install -y docker.io ;; \
	    centos|rhel)   sudo yum install -y docker ;; \
	    fedora)        sudo dnf install -y docker ;; \
	    arch)          sudo pacman -Sy --noconfirm docker ;; \
	    *) echo "Автоустановка не поддерживается для дистрибутива $$ID"; exit 1 ;; \
	  esac; \
	  sudo systemctl enable --now docker; \
	  echo "Docker установлен: $$(docker --version)"; \
	fi

swagger: swag-install  
	@echo "Генерируем Swagger локально..."
	swag init -g $(SWAG_MAIN) -o $(SWAG_OUT)
	@echo "Корректируем docs.go — удаляем LeftDelim и RightDelim..."
	sed -i '/LeftDelim:/d; /RightDelim:/d' $(SWAG_OUT)/docs.go
	@echo "Swagger-сборка завершена."
	@$(MAKE) docker-install
