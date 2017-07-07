# static files directory
STATIC_DIR = src/gui/static

# electron files directory
ELECTRON_DIR = electron

.PHONY: build clean

# build electron apps, the builds are located in electron/release folder.
build: 
	cd $(STATIC_DIR) && gulp dist
	cd $(ELECTRON_DIR) && ./build.sh
	@echo release files are in the folder of electron/release

# clean dist files and delete all builds in electron/release 
clean: 
	cd $(STATIC_DIR) && gulp clean
	rm $(ELECTRON_DIR)/release/*