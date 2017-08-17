#!/bin/bash

cd $(dirname $0)

GO_BIN=$(which go)
XSLT_BIN=$(which xsltproc)

BIN="xmldumper"
$GO_BIN build -o $BIN
chmod +x $BIN

FILE_INPUT=${1:-"sample.xml"}
if [ ! -f "${FILE_INPUT}" ]; then
	echo "ERROR: ${FILE_INPUT} doesn't not exists"
	exit 1
fi

PHONE_NUMBER=${2:-"0600000000"}

echo "Dump all SMS related to ${PHONE_NUMBER}..."
./$BIN -input ${FILE_INPUT} -output ${FILE_INPUT}.filtered -phonenumber ${PHONE_NUMBER}

echo "Regen HTML from XML files processing..."
$XSLT_BIN sms.xsl ${FILE_INPUT}.filtered > ${FILE_INPUT}.filtered.html

echo "Generated file: ${FILE_INPUT}.filtered.html"
echo "Hint: Open file inside a browser to print it as PDF"
