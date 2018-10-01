fmt:
	for file in $$(find -type f -name "*.go" | grep -v "vendor"); do \
		echo "Check $${file};" \
		gofmt -w -s "$${file}"; \
	done
