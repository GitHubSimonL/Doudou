package network

import (
	"strconv"
)

const (
	h_END_MASK = 0x4000
	h_ICS_MASK = 0x2000
	h_ICI_MASK = 0x1000
	h_LEN_MASK = 0x0FFF
)

func IPFieldMarch(str1, str2 []string) (ok bool) {
	if len(str1) != len(str2) {
		return
	}

	var turnInt func(s string) (result int64, err error)

	if len(str1) == 4 {
		turnInt = func(s string) (result int64, err error) {
			result, err = strconv.ParseInt(s, 10, 8)
			return
		}
	} else {
		turnInt = func(s string) (result int64, err error) {
			result, err = strconv.ParseInt(s, 16, 16)
			return
		}
	}

	for i := 0; i < len(str1); i++ {
		if str1[i] == "*" {
			continue
		}

		n1, err := turnInt(str1[i])

		if err != nil {
			return
		}

		n2, err := turnInt(str2[i])

		if err != nil {
			return
		}

		if n1 != n2 {
			return
		}
	}

	return true
}

func IPv6Full(fields []string) (result []string) {
	if len(fields) <= 0 {
		return
	}

	fullCount := 8 - len(fields)
	for _, field := range fields {
		if len(field) > 0 {
			result = append(result, field)
			continue
		}

		result = append(result, "0")

		for ; fullCount > 0; fullCount-- {
			result = append(result, "0")
		}
	}

	return
}

// func EncodeMsgHeader(isEnd, includeSeed, includeIP bool, msgLen int16) (data []byte, ok bool) {
// 	if msgLen > h_LEN_MASK {
// 		return
// 	}
//
// 	var header uint16
//
// 	if isEnd {
// 		header |= h_END_MASK
// 	}
//
// 	if includeSeed {
// 		header |= h_ICS_MASK
// 	}
//
// 	if includeIP {
// 		header |= h_ICI_MASK
// 	}
//
// 	header |= uint16(msgLen)
//
// 	return []byte{byte(header >> 8 & 0x00FF), byte(header & 0x00FF)}, true
// }
//
// func DecodeMsgHeader(data []byte) (isEnd, includeSeed, includeIP bool, msgLen int16, ok bool) {
// 	if len(data) < 2 {
// 		return
// 	}
//
// 	header := uint16(data[0])<<8 | uint16(data[1])
//
// 	if header&h_END_MASK > 0 {
// 		isEnd = true
// 	}
//
// 	if header&h_ICS_MASK > 0 {
// 		includeSeed = true
// 	}
//
// 	if header&h_ICI_MASK > 0 {
// 		includeIP = true
// 	}
//
// 	msgLen = int16(header) & h_LEN_MASK
// 	ok = true
//
// 	return
// }

// func PackageSplit(totalLen int) (lList []int) {
// 	fpkg := totalLen / h_LEN_MASK
// 	epkg := totalLen % h_LEN_MASK
//
// 	for i := 0; i <= fpkg; i++ {
// 		lList = append(lList, i*h_LEN_MASK)
// 	}
//
// 	if fpkg <= 0 || epkg > 0 {
// 		lList = append(lList, totalLen)
// 	}
//
// 	return
// }
