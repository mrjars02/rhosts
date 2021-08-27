dir: 
	if [ ! -d /usr/local/share/rhosts ];then mkdir -p /usr/local/share/rhosts;else echo "/usr/local/share/rhosts already exists";fi
	if [ ! -d /etc/rhosts ];then mkdir /etc/rhosts;else echo "/etc/rhosts already exists";fi
	if [ ! -d build ];then mkdir build;else echo "build already exists";fi
build: dir
	gcc src/rhosts.c src/download.c -lcurl -o build/rhosts
clean:
	- rm -rf build
install: build
	if [ ! -e /etc/rhosts/rhosts.cfg ];then touch /etc/rhosts/rhosts.cfg;else echo "/etc/rhosts/rhosts.cfg exists";fi
	cp -r src/systemd /usr/local/share/rhosts/
	cp build/rhosts /usr/local/bin/
	chown -R root:root /usr/local/bin/rhosts
	chmod +x /usr/local/bin/rhosts
	if [ ! -e /etc/systemd/system/rhosts.service ];then \
		ln -s /usr/local/share/rhosts/systemd/rhosts.service /etc/systemd/system/;else \
		echo "/etc/systemd/system/rhosts.service already exist";fi
	if [ ! -e /etc/systemd/system/rhosts.timer ];then \
		ln -s /usr/local/share/rhosts/systemd/rhosts.timer /etc/systemd/system/;else \
		echo "/etc/systemd/system/rhosts.timer already exist";fi
	if [ ! -e /etc/systemd/system/rhosts.path ];then \
		ln -s /usr/local/share/rhosts/systemd/rhosts.path /etc/systemd/system/;else \
		echo "/etc/systemd/system/rhosts.path already exist";fi
	systemctl daemon-reload
activate: install
	systemctl enable rhosts.timer
	systemctl start rhosts.timer
	systemctl enable rhosts.path
	systemctl start rhosts.path
deactivate: 
	- systemctl disable rhosts.timer
	- systemctl stop rhosts.timer
	- systemctl disable rhosts.path
	- systemctl stop rhosts.path
remove: deactivate
	rm -f /etc/systemd/system/rhosts*
	rm -fr /usr/local/share/rhosts/
	rm -f /usr/local/bin/rhosts
purge: remove
	rm -fr /etc/rhosts

reinstall: remove activate
