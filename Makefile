time := $(shell date +'%Y%m%d-%H%M%S')

publish-apply: publish apply

publish:
	fandogh image publish --version $(time)

apply:
	fandogh service apply -f fandogh.yml -p TAG=$(time) -d