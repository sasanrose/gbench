fmt:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		echo "Fix $${file};" \
		gofmt -w -s "$${file}"; \
	done

fix-misspell:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		misspell -w "$${file}"; \
	done

check-cyclo:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		gocyclo -over 15 "$${file}"; \
	done
