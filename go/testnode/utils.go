// https://www.socketloop.com/tutorials/golang-convert-integer-to-binary-octal-hexadecimal-and-back-to-integer
package main
import (  
    "strconv"
    "fmt"
    "math/rand"
)

// Convert Str(0x.....) -> int64
func hex2intr(hexStr string) int64 {
     d, _ := strconv.ParseInt("0x76c0", 0, 64)
    return d
}

// Convert int -> Str(0x.....) -> 
func intr2hex (i int) string {
     d:= fmt.Sprintf("0x%02x", i)
    return d
}

func bin(i int, prefix bool) string {
     i64 := int64(i)
  if prefix {
     return "0b" + strconv.FormatInt(i64, 2)  // base 2 for binary
  } else {
     return strconv.FormatInt(i64, 2)         // base 2 for binary
  }
}

func bin2int(binStr string) int {
     // base 2 for binary
     result, _ := strconv.ParseInt(binStr, 2, 64)
     return int(result)
}

func oct(i int, prefix bool) string {
    i64 := int64(i)

    if prefix {
       return "0o" + strconv.FormatInt(i64, 8) // base 8 for octal
    } else {
      return strconv.FormatInt(i64, 8) // base 8 for octal
    }
}

// Oct to integer
func oct2int(octStr string) int {
    // base 8 for octal
    result, _ := strconv.ParseInt(octStr, 8, 64)
    return int(result)
}

// Hex
func hex(i int, prefix bool) string {
    i64 := int64(i)

    if prefix {
       return "0x" + strconv.FormatInt(i64, 16) // base 16 for hexadecimal
    } else {
      return strconv.FormatInt(i64, 16) // base 16 for hexadecimal
    }
}

// Hex to integer
func hex2int(hexStr string) int {
          // base 16 for hexadecimal
          result, _ := strconv.ParseInt(hexStr, 16, 64)
          return int(result)
}

// Hex test
func binhextest(){
          num := 123456789
          fmt.Println("Integer : ", num)
          fmt.Println("Binary : ",  bin(num, false))
          fmt.Println("Octal : ",   oct(num, true))
          fmt.Println("Hex : ",     hex(num, true))

          // bin2int function does not handle the prefix
          // so set second parameter to false
          // otherwise you will get funny result

          fmt.Println("Binary to Integer : ",      bin2int(bin(num, false)))
          fmt.Println("Octal to Integer : ",       oct2int(oct(num, false)))
          fmt.Println("Hexadecimal to Integer : ", hex2int(hex(num, false)))
}

// "github.com/ethereum/go-ethereum/common/hexutil"
// func hextostring() {
//     addrBytes := []byte{20, 123, 142, 185, 127, 210, 71, 208, 108, 64, 6, 210, 105, 201, 12, 25, 8, 251, 93, 84}
//     fmt.Println(hexutil.Encode(addrBytes)) // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54

//     addrHex, _ := hexutil.Decode("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
//     fmt.Println(addrHex) // [20 123 142 185 127 210 71 208 108 64 6 210 105 201 12 25 8 251 93 84]
// }


/***************************************************************
 * Convert float64 to string
 ***************************************************************/
func FltoStr(Ints float64) string {
    str := fmt.Sprintf("%v", Ints)
    return str
}


// Convert string to int
func Str2int(Num string)int{
     s, _ := strconv.Atoi(Num)
     return s
}
// Возвращение random чисел в диапазоне
func Rendom(Par int) int {
     return rand.Intn(Par)
}