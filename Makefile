GREENFG = \e[32m
BLUEFG  = \e[34m
REDFG   = \e[31m
DEFFG   = \e[39m

CHECKTAG     = [$(BLUEFG)CHECK$(DEFFG)]
CHECKSTART   = $(CHECKTAG) $(BLUEFG)START$(DEFFG)
CHECKSUCCESS = $(CHECKTAG) $(GREENFG)SUCCESS$(DEFFG)
CHECKFAILED  = $(CHECKTAG) $(REDFG)FAILED$(DEFFG)

VETTAG     = [$(BLUEFG)VET$(DEFFG)]
VETSTART   = $(VETTAG) $(BLUEFG)START$(DEFFG)
VETSUCCESS = $(VETTAG) $(GREENFG)SUCCESS$(DEFFG)
VETFAILED  = $(VETTAG) $(REDFG)FAILED$(DEFFG)

LINTTAG     = [$(BLUEFG)LINT$(DEFFG)]
LINTSTART   = $(LINTTAG) $(BLUEFG)START$(DEFFG)
LINTSUCCESS = $(LINTTAG) $(GREENFG)SUCCESS$(DEFFG)
LINTFAILED  = $(LINTTAG) $(REDFG)FAILED$(DEFFG)

BUILDTAG     = [$(BLUEFG)BUILDPKG$(DEFFG)]
BUILDSTART   = $(BUILDTAG) $(BLUEFG)START$(DEFFG)
BUILDSUCCESS = $(BUILDTAG) $(GREENFG)SUCCESS$(DEFFG)
BUILDFAILED  = $(BUILDTAG) $(REDFG)FAILED$(DEFFG)

DEMOTAG     = [$(BLUEFG)DEMOBUILD$(DEFFG)]
DEMOSTART   = $(DEMOTAG) $(BLUEFG)START$(DEFFG)
DEMOSUCCESS = $(DEMOTAG) $(GREENFG)SUCCESS$(DEFFG)
DEMOFAILED  = $(DEMOTAG) $(REDFG)FAILED$(DEFFG)

CLEANTAG     = [$(BLUEFG)CLEAN$(DEFFG)]
CLEANSUCCESS = $(CLEANTAG) $(GREENFG)SUCCESS$(DEFFG)
CLEANFAILED  = $(CLEANTAG) $(REDFG)FAILED$(DEFFG)


default: build vet lint check demo

vet:
	@echo -e "$(VETSTART)"
	@GOOS=js GOARCH=wasm go vet \
	 && (echo -e "$(VETSUCCESS)") || (echo -e "$(VETFAILED)" && false)

lint:
	@echo -e "$(LINTSTART)"
	@GOOS=js GOARCH=wasm golint \
	 && (echo -e "$(LINTSUCCESS)") || (echo -e "$(LINTFAILED)" && false)

check:
	@echo -e "$(CHECKSTART)"
	@GOOS=js GOARCH=wasm go test -exec="${GOROOT}/misc/wasm/go_js_wasm_exec" \
	 && (echo -e "$(CHECKSUCCESS)") || (echo -e "$(CHECKFAILED)" && false)

build:
	@echo -e "$(BUILDSTART)"
	@GOOS=js GOARCH=wasm go build && (echo -e "$(BUILDSUCCESS)") || (echo -e "$(BUILDFAILED)" && false)

demo/main.wasm: demo/main.go *.go
	@echo -e "$(DEMOSTART)"
	@GOOS=js GOARCH=wasm go build -o demo/main.wasm demo/main.go \
	      && (echo -e "$(DEMOSUCCESS)") || (echo -e "$(DEMOFAILED)" && false)

demo: demo/main.wasm

clean:
	@rm -f demo/main.wasm && (echo -e "$(CLEANSUCCESS)") || (echo -e "$(CLEANFAILED)" && false)
