# base.mk - Common Makefile for LazyCAT Apps projects
# This file should be included in your project's Makefile
#
# Required tools:
#   lzc-cli         - LazyCAT CLI tool
#                     Install: npm install -g @lazycatcloud/lzc-cli
#                     Auto-completion: lzc-cli completion >> ~/.zshrc
#
# Required variables (define in your Makefile):
#   PROJECT_TYPE    - lpk-only or docker-lpk
#
# Optional variables:
#   PROJECT_NAME    - Project name (default: current directory name)
#   VERSION         - Project version (default: git tag or commit)
#   REGISTRY        - Docker registry (for docker-lpk projects)
#   IMAGE_NAME      - Docker image name (default: PROJECT_NAME)
#   APP_ID_PREFIX   - Application ID prefix (default: cloud.lazycat.app.liu.)
#   APP_NAME        - Application name (default: current directory name)
#   APP_ID          - Full Application ID (default: APP_ID_PREFIX + APP_NAME)

# Application ID configuration
ifndef APP_ID_PREFIX
    APP_ID_PREFIX := cloud.lazycat.app.liu.
endif

ifndef PROJECT_NAME
    PROJECT_NAME := $(notdir $(CURDIR))
endif

ifndef APP_NAME
    APP_NAME := $(shell basename $(CURDIR))
endif

ifndef APP_ID
    APP_ID := $(APP_ID_PREFIX)$(APP_NAME)
endif

# LazyCAT Box configuration
LAZYCAT_BOX_FALLBACK ?= 0
ifndef LAZYCAT_BOX_NAME
    LAZYCAT_BOX_NAME := $(shell command -v lzc-cli >/dev/null 2>&1 && lzc-cli box default 2>/dev/null)
    ifeq ($(strip $(LAZYCAT_BOX_NAME)),)
        LAZYCAT_BOX_NAME := default
        LAZYCAT_BOX_FALLBACK := 1
    endif
endif

# Version detection
ifndef VERSION
    VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
endif

# Docker-related variables
ifndef REGISTRY
    ifeq ($(LAZYCAT_BOX_FALLBACK),1)
        REGISTRY :=
    else
        REGISTRY := docker-registry-ui.$(LAZYCAT_BOX_NAME).heiyu.space
    endif
endif

ifeq ($(LAZYCAT_BOX_FALLBACK),1)
$(warning LazyCAT box name was not auto-detected; install lzc-cli or set LAZYCAT_BOX_NAME/REGISTRY to avoid using the fallback settings.)
endif

ifdef REGISTRY
    IMAGE_PREFIX := $(REGISTRY)/
else
    IMAGE_PREFIX :=
endif

ifndef IMAGE_NAME
    IMAGE_NAME := $(PROJECT_NAME)
endif

FULL_IMAGE_NAME := $(IMAGE_PREFIX)$(IMAGE_NAME):$(VERSION)

# Colors for output
COLOR_RESET   := \033[0m
COLOR_INFO    := \033[34m
COLOR_SUCCESS := \033[32m
COLOR_WARNING := \033[33m
COLOR_ERROR   := \033[31m

define print_info
    @echo "$(COLOR_INFO)[INFO]$(COLOR_RESET) $(1)"
endef

define print_success
    @echo "$(COLOR_SUCCESS)[SUCCESS]$(COLOR_RESET) $(1)"
endef

define print_warning
    @echo "$(COLOR_WARNING)[WARNING]$(COLOR_RESET) $(1)"
endef

define print_error
    @echo "$(COLOR_ERROR)[ERROR]$(COLOR_RESET) $(1)"
endef

# Default target
.DEFAULT_GOAL := help

##@ General

.PHONY: help-default
help-default: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: info-default
info-default: ## Show project information
	$(call print_info,Project: $(PROJECT_NAME))
	$(call print_info,Type: $(PROJECT_TYPE))
	$(call print_info,Version: $(VERSION))
	$(call print_info,App ID: $(APP_ID))
ifeq ($(PROJECT_TYPE),docker-lpk)
	$(call print_info,Image: $(FULL_IMAGE_NAME))
endif
ifeq ($(LAZYCAT_BOX_FALLBACK),1)
	$(call print_warning,LazyCAT box name not detected. Install lzc-cli or set LAZYCAT_BOX_NAME/REGISTRY for Docker workflows.)
endif

##@ Building

.PHONY: clean-default
clean-default: ## Clean build artifacts
	$(call print_info,Cleaning build artifacts...)
	rm -rf bin/ dist/ build/ *.out coverage.html htmlcov/ *.lpk
	$(call print_success,Cleaned)

##@ Docker (docker-lpk projects only)

ifeq ($(PROJECT_TYPE),docker-lpk)

.PHONY: docker-build-default
docker-build-default: ## Build Docker image
	$(call print_info,Building Docker image: $(FULL_IMAGE_NAME))
	docker build -t $(FULL_IMAGE_NAME) .
	$(call print_success,Docker image built: $(FULL_IMAGE_NAME))

.PHONY: docker-push-default
docker-push-default: docker-build-default ## Push Docker image to registry
	$(call print_info,Pushing Docker image: $(FULL_IMAGE_NAME))
	docker push $(FULL_IMAGE_NAME)
	$(call print_success,Docker image pushed: $(FULL_IMAGE_NAME))

.PHONY: docker-run-default
docker-run-default: ## Run Docker container locally
	$(call print_info,Running Docker container...)
	docker run --rm -it $(FULL_IMAGE_NAME)

endif

##@ Release

.PHONY: lpk-default
lpk-default: ## Package LPK
	$(call print_info,Building LPK package...)
	@command -v lzc-cli >/dev/null 2>&1 || ($(call print_error,lzc-cli not found. Please install LazyCAT CLI) && exit 1)
	lzc-cli project build
	$(call print_success,LPK package built successfully)

.PHONY: deploy-default
deploy-default: lpk-default ## Build and install LPK package
	$(call print_info,Installing LPK package...)
	@LPK_FILE=$$(ls -t *.lpk 2>/dev/null | head -n 1); \
	if [ -z "$$LPK_FILE" ]; then \
		$(call print_error,No LPK file found); \
		exit 1; \
	fi; \
	echo "Installing $$LPK_FILE..."; \
	lzc-cli app install "$$LPK_FILE"
	$(call print_success,Installation completed)

.PHONY: uninstall-default
uninstall-default: ## Uninstall the LPK package
	$(call print_info,Uninstalling $(APP_ID)...)
	lzc-cli app uninstall $(APP_ID)
	$(call print_success,Uninstallation completed)

.PHONY: list-packages-default
list-packages-default: ## List all LPK packages in current directory
	@echo "Available LPK packages:"
	@ls -lht *.lpk 2>/dev/null || echo "No LPK packages found"

.PHONY: all-default
all-default: help-default ## Default target: help

.PHONY: release-default
release-default: ## Create a release
ifeq ($(PROJECT_TYPE),docker-lpk)
	@$(MAKE) docker-push-default
endif
	@$(MAKE) lpk-default
	$(call print_success,Release $(VERSION) completed)

##@ Utilities

.PHONY: install-lzc-cli-default
install-lzc-cli-default: ## Install lzc-cli tool
	$(call print_info,Installing lzc-cli...)
	@command -v npm >/dev/null 2>&1 || ($(call print_error,npm not found. Please install Node.js first) && exit 1)
	npm install -g @lazycatcloud/lzc-cli
	$(call print_success,lzc-cli installed successfully)
	$(call print_info,To enable auto-completion, run: lzc-cli completion >> ~/.zshrc)

.PHONY: version-default
version-default: ## Show version
	@echo $(VERSION)

.PHONY: check-tools-default
check-tools-default: ## Check if required tools are installed
	$(call print_info,Checking required tools...)
	@command -v lzc-cli >/dev/null 2>&1 || ($(call print_error,lzc-cli not found) && exit 1)
ifeq ($(PROJECT_TYPE),docker-lpk)
	@command -v docker >/dev/null 2>&1 || ($(call print_error,docker not found) && exit 1)
endif
	$(call print_success,All required tools are installed)

# Pattern rule to allow overriding any -default target
# Usage: Define a target with the same name (without -default) in your Makefile to override
%: %-default
	@ true
