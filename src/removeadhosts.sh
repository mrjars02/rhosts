#!/bin/bash
sed -n '1,/# start of removeadhosts/p' /etc/hosts > /tmp/removeadhosts
cat /tmp/removeadhosts | tee /etc/hosts >/dev/null
rm /tmp/removeadhosts
echo "Appended old hosts"
if grep -qinE '# start of removeadhosts' /etc/hosts
then
        echo "This has been run before"
else
        echo "First Run"
        echo '# start of removeadhosts' >> /etc/hosts
fi

echo 'Downloading ad list'

if [ -e /etc/removeadhosts/adlistings.txt ]
then
        cat /etc/removeadhosts/adlistings.txt | \
        while read SITE; do
            curl $SITE 2>/dev/null | tee -a /etc/hosts >/dev/null
            echo "0.0.0.0 $CMD" >> /etc/hosts
        done
fi

echo 'Adding custom items from /etc/removeadhosts'
if [ -e /etc/removeadhosts/adlist.txt ]
then
        cat /etc/removeadhosts/adlist.txt | \
        while read CMD; do
            echo "0.0.0.0 $CMD" >> /etc/hosts
        done
fi
