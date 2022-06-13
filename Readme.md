# rhosts

This reroutes urls to 0.0.0.0 and ::1 in order to block them from being reached. This is useful for blocking different types of content.   

## How to use

Open the config file:    

        Linux: /etc/rhosts/rhosts.cfg

        Windows: \ProgramData\rhosts\rhosts.cfg


There are 3 types of entries: download, site, and whitelist. Downloads are downloaded and stripped of comments and bad entries if possible before being added to a list of sites. Whitelisted urls are removed from the list of sites. From there all the urls are added to the hosts file for both IPv4 and IPv6. You can also add comments by prepending with a '#'.    

Example:    

        # This is a static entry
        site=www.site.xyz
        # This is a download entry
        download=w3.site.xyz/location/to/config.txt
		# This is a whitelist entry
		whitelist=www.site.xyz

## How to Install

Build Dependencies:

- make
- golang

Linux/Systemd:

		make install

Windows:

		make build-win

The windows build can be performed on a Linux system and moved to Windows.
