#!/bin/bash
sed -n '1,/# start of removeadhosts/p' /etc/hosts > /tmp/removeadhosts-head
sed -n '/# start of removeadhosts/,/# Custom ad list/p' /etc/hosts | sed -e '/# Custom ad list/d' -e '1d' > /tmp/removeadhosts-curl
cat /tmp/removeadhosts-head | tee /etc/hosts >/dev/null
rm /tmp/removeadhosts-head
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

if [ $(sed -n '/# start of removeadhosts/,/# Custom ad list/p' /etc/hosts | sed -e '/# Custom ad list/d' -e '1d' | wc -l) -lt 2 ]
then
        echo "No hosts were downloaded, reusing the old ones"
        cat /tmp/removeadhosts-curl | tee -a /etc/hosts >/dev/null
fi
rm /tmp/removeadhosts-curl

echo 'Adding custom items from /etc/removeadhosts'
if [ -e /etc/removeadhosts/adlist.txt ]
then
        echo "# Custom ad list" >> /etc/hosts
        cat /etc/removeadhosts/adlist.txt | \
        while read CMD; do
            echo "0.0.0.0 $CMD" >> /etc/hosts
        done
fi
