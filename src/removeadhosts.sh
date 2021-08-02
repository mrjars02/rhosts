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

curl https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts 2>/dev/null | tee -a /etc/hosts >/dev/null

echo 'Adding custom items from /etc/removeadhosts'
if [ -e /etc/removeadhosts/adlist.txt ]
then
        cat /etc/removeadhosts/adlist.txt | \
        while read CMD; do
            echo "0.0.0.0 $CMD" >> /etc/hosts
        done
fi
