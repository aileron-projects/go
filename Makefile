# Help is in util.mk
.DEFAULT_GOAL:=help

ifneq (,$(wildcard .env.mk))
  include .env.mk
endif

# Override settings.
PRETTIER_TARGET := "**/*.{md,yaml,yml,toml,js,jsx,ts,html,css}"

include _scripts/makefiles/adoc.mk
include _scripts/makefiles/cspell.mk
include _scripts/makefiles/drawio.mk
include _scripts/makefiles/go-build.mk
include _scripts/makefiles/go-licenses.mk
include _scripts/makefiles/go-test.mk
include _scripts/makefiles/go.mk
include _scripts/makefiles/goda.mk
include _scripts/makefiles/golangci-lint.mk
include _scripts/makefiles/govulncheck.mk
include _scripts/makefiles/graphviz.mk
include _scripts/makefiles/markdownlint.mk
include _scripts/makefiles/mermaid.mk
include _scripts/makefiles/nfpm.mk
include _scripts/makefiles/plantuml.mk
include _scripts/makefiles/prettier.mk
include _scripts/makefiles/scanoss.mk
include _scripts/makefiles/shellcheck.mk
include _scripts/makefiles/shfmt.mk
include _scripts/makefiles/trivy.mk
include _scripts/makefiles/util.mk

LOCAL_CHECKS += go-licenses-run
LOCAL_CHECKS += golangci-lint-run
LOCAL_CHECKS += markdownlint-run
LOCAL_CHECKS += prettier-run

.PHONY: local-check
local-check: $(LOCAL_CHECKS)

.PHONY: local-format
local-format:
	$(MAKE) go-fmt ARGS="-w"
	$(MAKE) prettier-run ARGS="--write"
