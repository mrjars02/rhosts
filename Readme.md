# rhosts

This reroutes sites to 0.0.0.0 in order to block them from being reached by adding them automatically to the hosts file.   

## Requirements to install and run

### Linux

- gcc

- make

- libcurl4-gnutls-dev

## Install

### Linux

        sudo make install

## Update

### Linux

        sudo make reinstall

## How to use

Open the config file:    

        Linux: /etc/rhosts/rhosts.cfg

        Windows: \ProgramData\rhosts\rhosts.cfg


There are 2 types of entries: download and site. Downloads are currently downloaded straight into the config without error checking, sites are added with an IP address prepended. You can also add comments by prepending with a '#'.    

Example:    

        # This is a static entry
        site=www.site.xyz
        # This is a download entry
        download=w3.site.xyz/location/to/config.txt
