package can2mqtt_tuc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	CAN "github.com/brendoncarroll/go-can"
)

// convert2CAN does the following:
// 1. receive topic and payload
// 2. use topic to examine corresponding cconvertmode and CAN-ID
// 3. execute conversion
// 4. build CANFrame
// 5. returning the CANFrame
func convert2CAN(topic, payload string) CAN.CANFrame {
	convertMethod := getConvTopic(topic)
	var Id uint32 = uint32(getId(topic))
	var data [8]byte
	var len uint32
	if convertMethod == "none" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode none (reverse of %s)\n", convertMethod)
		}
		data, len = ascii2bytes(payload)
	} else if convertMethod == "uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint8 (reverse of %s)\n", convertMethod)
		}
		data[0] = ascii2uint8(payload)
		len = 8
	} else if convertMethod == "uint162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint16(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2uint16(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		len = 8
	} else if convertMethod == "uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint32(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2uint32(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		len = 8
	} else if convertMethod == "int322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2int32(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2int32(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		len = 8
	} else if convertMethod == "uint642ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint64(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2uint64(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		data[4] = tmp[4]
		data[5] = tmp[5]
		data[6] = tmp[6]
		data[7] = tmp[7]
		len = 8

	} else if convertMethod == "setup2floats" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode setup2floats %s)\n", convertMethod)
		}
		tmp := setup2floats(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		data[4] = tmp[4]
		data[5] = tmp[5]
		data[6] = tmp[6]
		data[7] = tmp[7]
		len = 8

	} else if convertMethod == "int32int16" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode int32int16 %s)\n", convertMethod)
		}
		nums := strings.Split(payload, " ")
		tmp := ascii2int32(nums[0])
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		tmp = ascii2int16(nums[1])
		data[4] = tmp[0]
		data[5] = tmp[1]
		len = 8

	} else if convertMethod == "setup2motor" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode motor einstellung %s)\n", convertMethod)
		}
		nums := strings.Split(payload, " ")
		tmp := ascii2int16(nums[0])
		data[0] = tmp[0]
		data[1] = tmp[1]
		tmp = ascii2int16(nums[1])
		data[2] = tmp[0]
		data[3] = tmp[1]
		tmp = ascii2int16(nums[2])
		data[4] = tmp[0]
		data[5] = tmp[1]
		len = 8

	} else if convertMethod == "2uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii22uint32(reverse of %s)\n", convertMethod)
		}
		nums := strings.Split(payload, " ")
		tmp := ascii2uint32(nums[0])
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		tmp = ascii2uint32(nums[1])
		data[4] = tmp[0]
		data[5] = tmp[1]
		data[6] = tmp[2]
		data[7] = tmp[3]
		len = 8

	} else if convertMethod == "float2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii22uint32(reverse of %s)\n", convertMethod)
		}
		//nums := strings.Split(payload, " ")
		tmp := ascii2dfloat(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]

		len = 8

	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback none\n", convertMethod)
		}
		data, len = ascii2bytes(payload)
	}
	mycf := CAN.CANFrame{ID: Id, Len: len, Data: data}
	return mycf
}

// convert2MQTT does the following
// 1. receive ID and payload
// 2. lookup the correct convertmode
// 3. executing conversion
// 4. building a string
// 5. return
func convert2MQTT(id int, length int, payload [8]byte) []string {
	convertMethod := getConvId(id)
	retstr := []string{}

	if convertMethod == "none" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode none\n")
		}
		retstr = append(retstr, bytes2ascii(uint32(length), payload))
		return retstr
	} else if convertMethod == "uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint82ascii\n")
		}
		retstr = append(retstr, uint82ascii(payload[0]))
		return retstr
	} else if convertMethod == "uint162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint162ascii\n")
		}
		retstr = append(retstr, uint162ascii(payload[0:2]))
		return retstr
	} else if convertMethod == "uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint322ascii\n")
		}
		retstr = append(retstr, uint322ascii(payload[0:4]))
		return retstr
	} else if convertMethod == "int322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: db int32  \n")
		}
		retstr = append(retstr, int322ascii(payload[0:4]))
		return retstr
	} else if convertMethod == "2int322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: db 2 int32  \n")
		}
		retstr = append(retstr, int322ascii(payload[0:4]))
		retstr = append(retstr, int322ascii(payload[4:8]))
		return retstr
	} else if convertMethod == "float2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: db float  \n")
		}
		retstr = append(retstr, dfloat2ascii(payload[0:4]))
		return retstr
	} else if convertMethod == "2float2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: db two float  \n")
		}
		retstr = append(retstr, dfloat2ascii(payload[0:4]))
		retstr = append(retstr, dfloat2ascii(payload[4:8]))
		return retstr
	} else if convertMethod == "uint642ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint642ascii\n")
		}
		retstr = append(retstr, uint642ascii(payload[0:8]))
		return retstr
	} else if convertMethod == "2uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 2uint322ascii\n")
		}
		retstr = append(retstr, uint322ascii(payload[0:4])+" "+uint322ascii(payload[4:8]))
		return retstr
	} else if convertMethod == "clock2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: Clock pulse\n")
		}
		retstr = append(retstr, uint322ascii(payload[0:4])+":"+uint322ascii(payload[4:8]))
		return retstr
	} else if convertMethod == "motor2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: motor msg pulk \n")
		}
		retstr = append(retstr, int162ascii(payload[0:2]))
		retstr = append(retstr, int162ascii(payload[2:4]))
		retstr = append(retstr, uint82ascii(payload[5]))
		return retstr
	} else if convertMethod == "setup2motor" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode setup2motor\n")
		}
		retstr = append(retstr, int162ascii(payload[0:2])+" "+int162ascii(payload[2:4])+" "+int162ascii(payload[4:6]))
		return retstr
	} else if convertMethod == "pixelbin2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode pixelbin2ascii\n")
		}
		retstr = append(retstr, uint82ascii(payload[0])+" "+bytecolor2colorcode(payload[1:4]))
		return retstr
	} else if convertMethod == "bytecolor2colorcode" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode bytecolor2colorcode\n")
		}
		retstr = append(retstr, bytecolor2colorcode(payload[0:2]))
		return retstr
	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback none\n", convertMethod)
		}
		retstr = append(retstr, bytes2ascii(uint32(length), payload))
		return retstr
	}
}

//######################################################################
//#				NONE				       #
//######################################################################

func bytes2ascii(length uint32, payload [8]byte) string {
	return string(payload[:length])
}

func ascii2bytes(payload string) ([8]byte, uint32) {
	var returner [8]byte
	var i uint32 = 0
	for ; int(i) < len(payload) && i < 8; i++ {
		returner[i] = payload[i]
	}
	return returner, i
}

// ######################################################################
// #			UINT82ASCII				       #
// ######################################################################
// uint82ascii takes exactly one byte and returns a string with a
// numeric decimal interpretation of the found data
func uint82ascii(payload byte) string {
	return strconv.FormatInt(int64(payload), 10)
}

func ascii2uint8(payload string) byte {
	return ascii2uint16(payload)[0]
}

// ######################################################################
// #			UINT162ASCII				       #
// ######################################################################
// uint162ascii takes 2 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint162ascii(payload []byte) string {
	if len(payload) != 2 {
		return "Err in CAN-Frame, data must be 2 bytes."
	}
	data := binary.BigEndian.Uint16(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint16(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint16(tmp)
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #			UINT322ASCII				       #
// ######################################################################
// uint322ascii takes 4 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint322ascii(payload []byte) string {
	if len(payload) != 4 {
		return "Err in CAN-Frame, data must be 4 bytes."
	}
	data := binary.LittleEndian.Uint32(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint32(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint32(tmp)
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #			INT322ASCII				       #
// ######################################################################
// int322ascii takes 4 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func int322ascii(payload []byte) string {
	if len(payload) != 4 {
		return "Err in CAN-Frame, data must be 4 bytes."
	}
	data := binary.BigEndian.Uint32(payload)
	data2 := int32(data)
	return strconv.FormatInt(int64(data2), 10)
}

func ascii2int32(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint32(tmp)
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #			UINT642ASCII				       #
// ######################################################################
// uint642ascii takes 8 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint642ascii(payload []byte) string {
	if len(payload) != 8 {
		return "Err in CAN-Frame, data must be 8 bytes."
	}
	data := binary.LittleEndian.Uint64(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint64(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint64(tmp)
	a := make([]byte, 8)
	binary.LittleEndian.PutUint64(a, number)
	return a
}

//########################################################################
//######################################################################
//#			double Float2ASCII				       #
//######################################################################
// Drillbotics double float msg -> two 32bit floats
// decimal interpretation of the found data as ascii-string

func dfloat2ascii(payload []byte) string {
	if len(payload) != 4 {
		return "Err in CAN-Frame, data must be 8 bytes."
	}
	data := binary.LittleEndian.Uint32(payload)
	float := math.Float32frombits(data)

	return strconv.FormatFloat(float64(float), 'f', 5, 32)
}

func ascii2dfloat(payload string) []byte {
	tmp, _ := strconv.ParseFloat(payload, 32)
	tmp2 := float32(tmp)
	number := math.Float32bits(tmp2)
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #             bytecolor2colorcode
// ######################################################################
// bytecolor2colorcode is a convertmode that converts between the binary
// 3 byte representation of a color and a string representation of a color
// as we know it (for example in html #00ff00 is green)
func bytecolor2colorcode(payload []byte) string {
	colorstring := hex.EncodeToString(payload)
	return "#" + colorstring
}

func colorcode2bytecolor(payload string) []byte {
	var a []byte
	var err error
	a, err = hex.DecodeString(strings.Replace(payload, "#", "", -1))
	if err != nil {
		return []byte{0, 0, 0}
	}
	return a
}

//########################################################################

// drilllbotics setup 2 floats

func setup2floats(payload string) []byte {
	nums := strings.Split(payload, " ")
	tmp_0, _ := strconv.ParseFloat(nums[0], 32)
	tmp_0f := float32(tmp_0)

	tmp_1, _ := strconv.ParseFloat(nums[1], 32)
	tmp_1f := float32(tmp_1)

	number_0 := math.Float32bits(tmp_0f)
	number_1 := math.Float32bits(tmp_1f)

	a_0 := make([]byte, 4)
	a_1 := make([]byte, 4)
	a := make([]byte, 8)
	binary.LittleEndian.PutUint32(a_0, number_0)
	binary.LittleEndian.PutUint32(a_1, number_1)

	a[0] = a_0[0]
	a[1] = a_0[1]
	a[2] = a_0[2]
	a[3] = a_0[3]
	a[4] = a_1[0]
	a[5] = a_1[1]
	a[6] = a_1[2]
	a[7] = a_1[3]

	return a
}

// this really annoyws me ... Little endian
func ascii2int16(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint16(tmp)
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, number)
	return a
}

func int162ascii(payload []byte) string {
	if len(payload) != 2 {
		return "Err in CAN-Frame, data must be 2 bytes."
	}
	data := binary.LittleEndian.Uint16(payload)
	data2 := int16(data)
	return strconv.FormatInt(int64(data2), 10)
}
