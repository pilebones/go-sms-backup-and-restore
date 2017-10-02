package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type XMLSMSes struct {
	XMLName    xml.Name `xml:"smses"`
	Count      int      `xml:"count,attr"`
	BackupSet  string   `xml:"backup_set,attr"`
	BackupDate string   `xml:"backup_date,attr"`
	SMSes      []XMLSMS `xml:"sms"`
}

type Status string

func (s Status) IsNull() bool {
	return s == "null"
}

func (s Status) Code() (int, error) {
	i, err := strconv.ParseInt(string(s), 10, 0)
	return int(i), err
}

type XMLSMS struct {
	XMLName       xml.Name `xml:"sms"`
	Protocol      string   `xml:"protocol,attr"`
	Address       string   `xml:"address,attr"`
	Date          int      `xml:"date,attr"`
	Type          int      `xml:"type,attr"`
	Subject       string   `xml:"subject,attr"`
	Toa           string   `xml:"toa,attr"`
	ScToa         string   `xml:"sc_toa,attr"`
	Body          string   `xml:"body,attr"`
	Number        string   `xml:"number,attr"`
	ServiceCenter string   `xml:"service_center,attr"`
	Read          bool     `xml:"read,attr"`
	Status        Status   `xml:"status,attr"`
	Locked        bool     `xml:"locked,attr"`
	DateSent      int      `xml:"date_sent,attr"`
	ReadableDate  string   `xml:"readable_date,attr"`
	ContactName   string   `xml:"contact_name,attr"`
}

func ReadSMSes(reader io.Reader) (xmlSMSes XMLSMSes, err error) {
	err = xml.NewDecoder(reader).Decode(&xmlSMSes)
	return
}

func main() {
	input := flag.String("input", "input.xml", "Input absolute file path")
	output := flag.String("output", "filtered.xml", "Output absolute file path")
	phoneNumber := flag.String("phonenumber", "0600000000", "Filter messages with this phone number")
	flag.Parse()

	if phoneNumber == nil {
		log.Println("No phone number provided, abort")
		os.Exit(1)
	}

	// Build the location of the sms.xml file
	// filepath.Abs appends the file name to the default working directly
	smsesFilePath, err := filepath.Abs(*input)
	if err != nil {
		fmt.Println("Unable to get absolute path of input file, err:", err)
		os.Exit(1)
	}

	// Open the smses.xml file
	file, err := os.Open(smsesFilePath)
	if err != nil {
		fmt.Println("Unable to open input file, err:", err)
		os.Exit(1)
	}

	defer file.Close()

	log.Println("Dump SMSes from ", *input, "to", *output)

	// Read the straps file
	xmlSMSes, err := ReadSMSes(file)
	if err != nil {
		fmt.Println("Unable to read input file, err:", err)
		os.Exit(1)
	}

	log.Println("SMS found :", len(xmlSMSes.SMSes))
	if len(xmlSMSes.SMSes) > 0 {

		smsesKeep := XMLSMSes{
			XMLName:    xmlSMSes.XMLName,
			BackupDate: xmlSMSes.BackupDate,
			BackupSet:  xmlSMSes.BackupSet,
			SMSes:      make([]XMLSMS, 0),
		}

		log.Println("Filter with", NormalizePhoneNumber(*phoneNumber), "as phone number")

		for _, sms := range xmlSMSes.SMSes {
			if NormalizePhoneNumber(sms.Address) == NormalizePhoneNumber(*phoneNumber) {
				smsesKeep.SMSes = append(smsesKeep.SMSes, sms)
			}
		}

		smsesKeep.Count = len(smsesKeep.SMSes)

		log.Println("SMS keep after filtering : ", smsesKeep.Count)

		data, err := xml.Marshal(smsesKeep)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(*output, data, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Println("Dump SMS to", *output, "finished")
	}
}

// return phone number from :
// - +33 X XX XX XX XX
// - +33XXXXXXXXX to :
// - 0XXXXXXXXX
// FIXME: I18nize this func, just yet devel for my needs
func NormalizePhoneNumber(num string) string {
	rv := strings.Replace(num, " ", "", -1)
	if strings.HasPrefix(rv, "+33") {
		rv = "0" + strings.TrimPrefix(rv, "+33")
	}
	return rv
}
