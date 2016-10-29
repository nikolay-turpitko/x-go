.PHONY: \
	all \
	clean \
	test \
	build

# This file is used to automate test and check build of multiple subprojects
# with Travis CI.

all: \
	clean \
	test \
	build
	@$(MAKE) clean --no-print-directory
	@echo "All Done."

clean:
	@echo "Clean"
	@glide nv | xargs go clean -i -r
	@echo "Clean Done."

test:
	@echo "Test"
	@glide nv | xargs go test
	@-glide nv | xargs go vet
	@-glide nv | xargs -L 1 golint
	@echo "Test Done."

build:
	@echo "Build"
	@glide nv | xargs go build
	@echo "Build Done."
	@echo "Note: this build only checks if it can be done, it does not preserve output files."
