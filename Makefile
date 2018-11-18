GREENFG = \e[32m
BLUEFG  = \e[34m
REDFG   = \e[31m
DEFFG   = \e[39m

CHECKTAG     = [$(BLUEFG)CHECK$(DEFFG)]
CHECKSTART   = $(CHECKTAG) $(BLUEFG)START$(DEFFG)
CHECKSUCCESS = $(CHECKTAG) $(GREENFG)SUCCESS$(DEFFG)
CHECKFAILED  = $(CHECKTAG) $(REDFG)FAILED$(DEFFG)

BUILDTAG     = [$(BLUEFG)BUILD$(DEFFG)]
BUILDSTART   = $(BUILDTAG) $(BLUEFG)START$(DEFFG)
BUILDSUCCESS = $(BUILDTAG) $(GREENFG)SUCCESS$(DEFFG)
BUILDFAILED  = $(BUILDTAG) $(REDFG)FAILED$(DEFFG)

default: build check

check:
	@echo -e "$(CHECKSTART)"
	@GOOS=js GOARCH=wasm go test -exec="${HOME}/devel/go/misc/wasm/go_js_wasm_exec" \
	 && (echo -e "$(CHECKSUCCESS)") || (echo -e "$(CHECKFAILED)" && false)

build:
	@echo -e "$(BUILDSTART)"
	@GOOS=js GOARCH=wasm go build && (echo -e "$(BUILDSUCCESS)") || (echo -e "$(BUILDFAILED)" && false)
