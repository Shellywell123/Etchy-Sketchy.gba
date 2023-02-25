TARGET=gameboy-advance
ROM=Etchy-Sketchy.gba

.PHONY: build
build:
	cd src && tinygo build -size short -o ../${ROM} -target ${TARGET} etchy_sketchy.go && cd ..