package utils

import (
	"regexp"
)

func FormatMSISDNKE(msisdn string) string {

	formatedNumber := ""

	if len(msisdn) < 9 {

		return formatedNumber
	}

	var re = regexp.MustCompile(`[\n\r\s\D]+`)
	msisdn = re.ReplaceAllString(msisdn, "")

	if len(msisdn) < 9 {

		return formatedNumber

	}

	input1 := msisdn[0:1]
	input2 := msisdn[0:2]
	input4 := msisdn[0:4]
	input5 := msisdn[0:5]

	length := len(msisdn)

	if input1 == "7" && length == 9 {

		formatedNumber = SubstrReplace(msisdn, "2547", 0, 1)

	} else if input1 == "1" && length == 9 {

		formatedNumber = SubstrReplace(msisdn, "2541", 0, 1)

	} else if input2 == "07" && length == 10 {

		formatedNumber = SubstrReplace(msisdn, "254", 0, 1)

	} else if input2 == "01" && length == 10 {

		formatedNumber = SubstrReplace(msisdn, "254", 0, 1)

	} else if input2 == "7" {
		formatedNumber = SubstrReplace(msisdn, "2547", 0, 1)

	} else if input2 == "1" {
		formatedNumber = SubstrReplace(msisdn, "2547", 0, 1)

	} else if input4 == "2547" && length == 12 {

		formatedNumber = msisdn

	} else if input4 == "2541" && length == 12 {

		formatedNumber = msisdn

	} else if input5 == "+2547" && length == 13 {

		formatedNumber = SubstrReplace(msisdn, "", 0, 1)

	} else if input5 == "+2541" && length == 13 {

		formatedNumber = SubstrReplace(msisdn, "", 0, 1)

	}

	return formatedNumber
}

func SubstrReplace(haystack string, needle string, start int, end int) string {

	if start == 0 {

		parts := haystack[end:]
		return needle + parts
	}

	part1 := haystack[0:start]
	part2 := haystack[end:]

	return part1 + needle + part2
}
