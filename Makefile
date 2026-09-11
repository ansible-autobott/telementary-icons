COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIRR:= ${CURDIR}
SHELL := /bin/bash

# Version for local builds. Releases in CI derive it from the pushed git tag
# (see .github/workflows/release.yml); nfpm strips any leading "v".
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo 0.0.0-dev)

OUTPUT ?= $(CURDIR)/output
TELA ?= $(CURDIR)/tela-icons
ICONS_DIR ?= /usr/share/icons
ADDR ?= :3033

default: help

#==========================================================================================
##@ Building
#==========================================================================================
build: theme ## build the debian package (regenerates the theme first; version from git describe, override with VERSION=x.y.z)
	@VERSION=$(VERSION) nfpm package -f nfpm.yaml -p deb

clean: ## clean build artifacts (.deb and the generated output/ theme)
	@rm -f *.deb
	@rm -rf $(OUTPUT)

#==========================================================================================
##@ Icons
#==========================================================================================
# icon theme toolbox (tools/). Variables consumed by the targets below:
#   OUTPUT    generated icon theme root served/regenerated    (default $(CURDIR)/output)
#   TELA      vendored Tela base dir (src/+links/) laid down by `theme` first
#             (default $(CURDIR)/tela-icons)
#   ICONS_DIR dir scanned by `view` for base themes — each offered as a selectable
#             reference and fallback            (default /usr/share/icons)
#   REF       extra base theme path for `view` (e.g. an uninstalled source),
#             offered in the dropdown alongside the discovered ones (optional)
#   ADDR      address for the `make view` dev server           (default :3033)
#   APPLY     set APPLY=1 to actually copy in `overlay` (default: dry-run)

theme: ## regenerate output/ from vendored Tela + overlays (wipe + tela base + overlay + gen-theme)
	@cd tools && go run . theme -tela "$(TELA)" -theme "$(OUTPUT)"

view: theme ## serve a browsable view of the theme; pick a base theme from /usr/share/icons in the UI (override OUTPUT=dir, ICONS_DIR=dir, REF=path, ADDR=:port)
	@cd tools && go run . view -theme "$(OUTPUT)" -icons-dir "$(ICONS_DIR)" $(if $(REF),-ref "$(REF)") -addr $(ADDR)

overlay: ## apply design overlays: copy each overlays/*/overlay.yaml source straight into output/ (usage: make overlay APPLY=1; omit APPLY for a dry-run)
	@cd tools && go run . overlay -theme "$(OUTPUT)" $(if $(APPLY),-apply)

gen-theme: ## regenerate the theme's index.theme, @2x/@3x scaled dirs, and alias symlinks from the on-disk tree (override OUTPUT=dir)
	@cd tools && go run . gen-theme -theme "$(OUTPUT)"

#==========================================================================================
##@ Release
#==========================================================================================
tag: check-git-clean check-branch ## create and push a git tag to trigger the release workflow (usage: make tag version="v1.2.3")
	@[ "${version}" ] || ( echo ">> version is not set, usage: make tag version=\"v1.2.3\""; exit 1 )
	@git tag -d $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push --delete origin $(version) || true
	@git push origin $(version)

check-branch:
	@current_branch=$$(git symbolic-ref --short HEAD) && \
	if [ "$$current_branch" != "main" ]; then \
		echo "Error: You are on branch '$$current_branch'. Please switch to 'main'."; \
		exit 1; \
	fi

check-git-clean:
	@git diff --quiet && git diff --cached --quiet || ( echo "Error: git working tree is not clean, commit or stash first."; exit 1 )

#==========================================================================================
##@ Help
#==========================================================================================
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
