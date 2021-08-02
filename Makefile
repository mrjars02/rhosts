dir: 
	if [ ! -d /usr/local/share/removeadhosts ];then mkdir -p /usr/local/share/removeadhosts;fi
	if [ ! -d /etc/removeadhosts ];then mkdir /etc/removeadhosts;fi
install: dir
	touch /etc/removeadhosts/ads.txt
	cp src/* /usr/local/share/removeadhosts/
	chown -R root:root /usr/local/share/removeadhosts
	chmod +x /usr/local/share/removeadhosts/removeadhosts.sh
	cp /usr/local/share/removeadhosts/removeadhosts.service /etc/systemd/system/
	cp /usr/local/share/removeadhosts/removeadhosts.timer /etc/systemd/system/
activate: install
	systemctl enable removeadhosts.timer
	systemctl start removeadhosts.timer
deactivate: 
	systemctl disable removeadhosts.timer
	systemctl stop removeadhosts.timer
remove: deactivate
	rm -f /etc/systemd/system/removeadhosts*
	rm -fr /usr/local/share/removeadhosts/
purge: remove
	rm -fr /etc/removeadhosts

reinstall: remove activate
