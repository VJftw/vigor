.PHONY: test

test:
	find . -maxdepth 1 -mindepth 1 -type d ! -name "*-TODO-*" | sort -u | \
		while IFS= read -r dir; do \
			$(MAKE) -C $$dir test || exit 1; \
		done
