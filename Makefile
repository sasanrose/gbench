GOCYCLO := $(shell type gocyclo 2>/dev/null)
INEFFASSIGN := $(shell type ineffassign 2>/dev/null)
MISSPELL := $(shell type misspell 2>/dev/null)

fmt:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		echo "Fix $${file}"; \
		gofmt -s -w "$${file}"; \
	done

fix-misspell:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		misspell -w "$${file}"; \
	done

check-cyclo:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		gocyclo -over 15 "$${file}"; \
	done

check-golint:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		golint "$${file}"; \
	done

setup-dev:
ifndef GOCYCLO
	go get github.com/fzipp/gocyclo
endif
ifndef INEFFASSIGN
	go get github.com/gordonklaus/ineffassign
endif
ifndef MISSPELL
	go get github.com/client9/misspell/cmd/misspell
endif
	if [ ! -e bin/git-hooks ]; then \
		wget https://raw.githubusercontent.com/sasanrose/git-hooks/master/git-hooks -O bin/git-hooks && chmod u+x bin/git-hooks && bin/git-hooks --install bin; \
	else \
		bin/git-hooks --uninstall && bin/git-hooks --install bin; \
	fi;
