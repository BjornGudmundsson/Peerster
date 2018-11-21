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

//CreateBudgetList takes in two integers
//n and m and returns a list of uint64 that
//are an even way to split them up such that no index
//is more than 1 greater than any other index
func CreateBudgetList(n, m uint64) []uint64 {
	if n <= m {
		temp := make([]uint64, n)
		div := m / n
		mod := m % n
		for i := uint64(0); i < n; i++ {
			temp[i] = temp[i] + uint64(div)
		}
		for i := uint64(0); i < mod; i++ {
			temp[i] = temp[i] + uint64(1)
		}
		return temp
	}
	temp := make([]uint64, m)
	for i := uint64(0); i < m; i++ {
		temp[i] = 1
	}
	return temp
}
