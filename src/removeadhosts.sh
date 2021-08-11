#!/bin/bash

# Create temporary files to hold the old configs
sed -n '1,/# start of removeadhosts/p' /etc/hosts > /tmp/removeadhosts-head
sed -n '/# start of removeadhosts/,/# Custom ad list/p' /etc/hosts | sed -e '/# Custom ad list/d' -e '1d' > /tmp/removeadhosts-curl

# Copy back all the custom entries
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


# Download entries from the listings list
echo 'Downloading ad list'
if [ -e /etc/removeadhosts/adlistings.txt ]
then
        cat /etc/removeadhosts/adlistings.txt | \
        while read SITE; do
            ESC_SITE=$(printf '%s\n' "$SITE" | sed -e 's/[\/&]/\\&/g')
            echo "# removeadhosts site - $SITE" >> /tmp/removeadhosts-curlbuff
            RC=0 ; curl $SITE >> /tmp/removeadhosts-curlbuff 2>/dev/null || RC=$?
            echo "# removeadhosts site - $SITE - end" >> /tmp/removeadhosts-curlbuff
            if [ $(cat /tmp/removeadhosts-curlbuff | wc -l) -lt 3 ] || [ ! $RC -eq 0 ]
            then
                    if [ $(sed -n "/removeadhosts site - $ESC_SITE/,/removeadhosts site - $ESC_SITE - end/p" /tmp/removeadhosts-curl | wc -l) -gt 2 ]
                    then
                            echo "Keeping old version of $SITE"
                            if [ $(cat /tmp/removeadhosts-curlbuff | wc -l) -eq 2 ]
                            then
                                    echo "Nothing was downloaded"
                            else
                                    echo "New version is $(cat /tmp/removeadhosts-curlbuff | wc -l) lines long"
                            fi
                            sed -n "/removeadhosts site - $ESC_SITE/,/removeadhosts site - $EXC_SITE - end/p" /tmp/removeadhosts-curl | tee -a /etc/hosts > /dev/null
                    else
                            echo "Unable to add $SITE"
                            
                    fi
            else
                    echo "Updating $(cat /tmp/removeadhosts-curlbuff | wc -l) lines from $SITE"
                    cat /tmp/removeadhosts-curlbuff | tee -a /etc/hosts >/dev/null
            fi
        done
fi
rm /tmp/removeadhosts-curl
rm /tmp/removeadhosts-curlbuff


# Add entries from adlist
echo 'Adding custom items from /etc/removeadhosts'
if [ -e /etc/removeadhosts/adlist.txt ]
then
        echo "# Custom ad list" >> /etc/hosts
        cat /etc/removeadhosts/adlist.txt | \
        while read CMD; do
            echo "0.0.0.0 $CMD" >> /etc/hosts
        done
fi
