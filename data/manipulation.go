package data

import (
	"strconv"
	"strings"
)

//SplitIP takes a string in the tradition
//IPv4 format and returns the byte representation
//of each part of the IP address string
func SplitIP(ip string) []byte {
	nums := strings.Split(ip, ".")
	arr := make([]byte, len(nums))
	for i, val := range nums {
		n, _ := strconv.Atoi(val)
		arr[i] = byte(n)
	}

	return arr
}

//FormatPeers takes in a string of IP-addresses seperated by commas
//and splits them up at the commas
func FormatPeers(peers string) []string {
	split := strings.Split(peers, ",")
	return split
}

//SliceToBoolMap takes an array of strings and
//returns a map where the keys are the array values
//and the values are bools
func SliceToBoolMap(arr []string) map[string]bool {
	m := make(map[string]bool)
	for _, val := range arr {
		m[val] = true
	}
	return m
}
