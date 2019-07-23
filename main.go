package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/argandas/serial"
)

const (
	rtuMinSize = 4
	rtuMaxSize = 256
)

var (
	atmhz  *serial.SerialPort
	strhex string
	data   [256]byte
	line   uint8
)

func main() {

	atmhz = serial.New()
	err := atmhz.Open("/dev/ttyUSB0", 9600, 5*time.Second)
	if err != nil {
		log.Println("PORT BUSY")
	} else {
		log.Println("SUCCESS OPEN PORT")
	}

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			strhex = ""

			query := []uint8{255, 001, 134, 000, 000, 000, 000, 000, 121}
			atmhz.Write(query)
			time.Sleep(50 * time.Millisecond)

			for x := 0; x < 10; x++ {
				line, _ = atmhz.Read()
				//fmt.Printf("%x  ", line)

				str := fmt.Sprintf("%x", line)
				strhex += str + ","
			}
			//fmt.Println()
			log.Printf("strhex:%s", strhex)
			if strings.Contains(strhex, "86") && strings.Contains(strhex, ",0") {
				mainHEX := strhex[strings.Index(strhex, "86")+3 : strings.Index(strhex, ",0")]
				log.Printf("mainHEX:%s", mainHEX)

				THEX := mainHEX[strings.LastIndex(mainHEX, ",")+1 : len(mainHEX)]
				log.Printf("THEX:%s", THEX)

				HiHEX := mainHEX[0:strings.Index(mainHEX, ",")]
				log.Printf("HiHEX:%s", HiHEX)

				LoHEX := mainHEX[strings.Index(mainHEX, ",")+1 : strings.LastIndex(mainHEX, ",")]
				log.Printf("LoHEX:%s", LoHEX)

				Suhu, err := strconv.ParseInt(hexaNumberToInteger(THEX), 16, 64)
				if err == nil {
					Suhu -= 40
					log.Printf("Suhu:%v C", Suhu)
				}

				DecHi, err := strconv.ParseInt(hexaNumberToInteger(HiHEX), 16, 64)
				if err == nil {
					log.Printf("DecHi:%v", DecHi)
				}

				DecLi, err := strconv.ParseInt(hexaNumberToInteger(LoHEX), 16, 64)
				if err == nil {
					log.Printf("DecLi:%v", DecLi)
				}

				ConCO2 := (DecHi * 256) + DecLi
				log.Printf("ConCO2:%v ppm", ConCO2)
			} else {
				log.Println("can't find hex")
			}
			fmt.Println()
		}
	}
}
func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

/*
2019/07/23 15:11:12 strhex:ff,86,5,5f,44,0,0,0,d2,0,
2019/07/23 15:11:12 mainHEX:5,5f,44

*/
