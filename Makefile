update_submodule:
	git submodule init
	git submodule update --remote --merge

.PHONY: update_submodule