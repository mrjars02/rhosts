PREFIX=/usr/local
EXEC_PREFIX=$(DESTDIR)$(PREFIX)
BINDIR=$(EXEC_PREFIX)/bin
DATAROOTDIR=$(DESTDIR)$(PREFIX)/share
DATADIR=$(DATAROOTDIR)
$MANDIR=$DATAROOTDIR/man
#$INFODIR=$DATAROOTDIR/info
#$DOCDIR=$DATAROOTDIR/doc
VERSION=`cat version`
PROJROOT=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TARBALLPREFIX=rhosts-$(VERSION)
TARBALLNAME=$(TARBALLPREFIX).tar.gz
GOBUILDFLAGS=
GITOFF=0

build:
	if [ ! -d $(PROJROOT)/build ]; then \
		mkdir -p $(PROJROOT)build/share/rhosts/systemd $(PROJROOT)build/bin \
	;fi

	echo "package main\nvar version string=\"$(VERSION)\"" > $(PROJROOT)src/version.go

	cd $(PROJROOT)src && go build -o $(PROJROOT)build/bin/ $(GOBUILDFLAGS) ./
	cp -r $(PROJROOT)src/systemd $(PROJROOT)/build/share/rhosts/
build-win:
	if [ ! -d $(PROJROOT)/build ]; then \
		mkdir -p $(PROJROOT)build \
	;fi
	cd $(PROJROOT)src && GOOS=windows go build -o $(PROJROOT)build/ $(GOBUILDFLAGS) ./
install: build
	install -D $(PROJROOT)build/bin/rhosts $(BINDIR)/
	cp -r  $(PROJROOT)build/share/rhosts $(DATADIR)
	ln -s $(DATADIR)rhosts/systemd/rhosts.service /usr/lib/systemd/system/rhosts.service
	ln -s $(DATADIR)rhosts/systemd/rhosts.path /usr/lib/systemd/system/rhosts.path
	ln -s $(DATADIR)rhosts/systemd/rhosts.timer /usr/lib/systemd/system/rhosts.timer
uninstall:
	if [ -f $(BINDIR)/rhosts ]; then \
	rm $(BINDIR)/rhosts \
	;fi
	if [ -d $(DATADIR)/rhosts ]; then \
	rm -r $(DATADIR)/rhosts \
	;fi
	if [ -h /usr/lib/systemd/system/rhosts.service ]; then \
		rm /usr/lib/systemd/system/rhosts.service \
	;fi
	if [ -h /usr/lib/systemd/system/rhosts.path ]; then \
		rm /usr/lib/systemd/system/rhosts.path \
	;fi
	if [ -h /usr/lib/systemd/system/rhosts.timer ]; then \
		rm /usr/lib/systemd/system/rhosts.timer \
	;fi
clean:
	if [ -d $(PROJROOT)build ]; then \
		rm -r $(PROJROOT)build \
	;fi
	if [ -f $(PROJROOT)src/version.go ]; then \
		rm -r $(PROJROOT)src/version.go \
	;fi
	if [ -f $(PROJROOT)$(TARBALLNAME) ]; then \
		rm $(PROJROOT)$(TARBALLNAME) \
	;fi
dist: clean
	if [ -d $(PROJROOT).git ] && [ $(GITOFF) = 0 ]; then \
		git archive --format=tar.gz -o $(PROJROOT)$(TARBALLNAME) --prefix=$(TARBALLPREFIX)/ `git branch --show-current` \
	;else \
		mkdir $(PROJROOT)$(TARBALLPREFIX) && \
		find $(PROJECTROOT)* -maxdepth 0 -name $(TARBALLPREFIX) -prune -o -exec cp -r {} $(PROJROOT)$(TARBALLPREFIX)/ \; && \
		tar -czf $(PROJROOT)$(TARBALLNAME) -C $(PROJROOT) $(TARBALLPREFIX) && \
		rm -r $(PROJROOT)$(TARBALLPREFIX) \
	;fi
