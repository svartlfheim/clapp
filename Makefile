GO_IMAGE=golang:1.15
WD = $(shell pwd)
DOCKER_RUN=docker run -it -v "$(WD):/srv" -w /srv $(GO_IMAGE)

# This is a tweak of the following suggestions:
# https://gist.github.com/prwhite/8168133#gistcomment-1420062
help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-30s %s\n" "target" "help" ; \
	printf "%-30s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
			IFS=$$':' ; \
			help_split=($$help_line) ; \
			help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			printf '\033[36m'; \
			printf "%-30s %s" $$help_command ; \
			printf '\033[0m'; \
			printf "%s\n" $$help_info; \
	done

.PHONY: test
test: ## Run the tests for the package
	$(DOCKER_RUN) go test -cover ./...

.PHONY: tidy
tidy: ## Tidy the dependencies for the package
	$(DOCKER_RUN) go mod tidy

.PHONY: runtime
runtime: ## Start a golang runtime in docker with this code mounted
	$(DOCKER_RUN) bash