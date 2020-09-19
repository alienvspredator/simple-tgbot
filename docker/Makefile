.PHONY: build-%
build-%:
	@docker build \
		--tag $(*):latest \
		--file $(*).dockerfile \
		../