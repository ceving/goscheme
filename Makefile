PROGRAM=goscheme

ifndef GOROOT
$(error GOROOT undefined)
endif

ifndef GOARCH
$(error GOARCH undefined)
endif

ifeq ($(GOARCH),386)
GOARCHID=8
endif

ifeq ($(GOARCH),amd64)
GOARCHID=6
endif

ifndef GOARCHID
$(error Can not identify architecture)
endif

all: $(PROGRAM)

clean: 
	rm -f $(PROGRAM) *.$(GOARCHID)

GC=$(GOROOT)/bin/$(GOARCHID)g
GL=$(GOROOT)/bin/$(GOARCHID)l

%.$(GOARCHID): %.go
	$(GC) $<

$(PROGRAM): main.$(GOARCHID)
	$(GL) -o $@ $<

main.$(GOARCHID): scheme.$(GOARCHID)
