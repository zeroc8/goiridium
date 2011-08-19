package utils

import (
	"fmt"
	"bytes")

const(
	MO_HDR = 0x01
	MO_PLD = 0x02
	MO_LOC = 0x03
	MT_HDR = 0x41
	MT_PLD = 0x42
	MT_CNF = 0x44
)
const (
	MTFL_FLUSH_MT_QUEUE = 1
	MTFL_SEND_RA = 2
	MTFL_FORCE_RA = 4
	MTFL_UPDATE_SSD_LOC = 8
	MTFL_HIGH_PRIORITY = 16
)

func Encode(imei,msg,counter string, mtflags uint16) *bytes.Buffer {		
	
	mtbuf := bytes.NewBuffer([]byte {})	
	mthdr := bytes.NewBuffer([]byte {})
	mtpld := bytes.NewBuffer([]byte {})
	
	// header
	mthdr.WriteByte(MT_HDR)
	mthdr.WriteByte(0)
	mthdr.WriteByte(21)
	mthdr.WriteString(counter)
	mthdr.WriteString(imei)
	mthdr.WriteByte(uint8(mtflags >> 8))
	mthdr.WriteByte(uint8(mtflags))
		
	// mtpayload	
	mtpld.WriteByte(MT_PLD)
	var pld_len uint16 = uint16(len(msg))	
	mtpld.WriteByte(uint8(pld_len >> 8))
	mtpld.WriteByte(uint8(pld_len))
		
	mtpld.WriteString(msg)
		
	var total_len uint16 = uint16(3+mthdr.Len()+mtpld.Len())
	
	mtbuf.WriteByte(1) // protocol revision number
	mtbuf.WriteByte(uint8(total_len) >> 8)
	mtbuf.WriteByte(uint8(total_len))
	
	mtbuf.Write(mthdr.Bytes())	
	mtbuf.Write(mtpld.Bytes())
		
	return mtbuf;	
	
}

func Decode(buf [] uint8) (imei, payload string) {	

	total_length := int(buf[1] <<8 | buf[2])
	ieis := buf[3:3+total_length]
	
	offset := 0  
	total:=0

	for ;; {			  
		
		iei_id := ieis[offset]	  		
		
		// total length of iei
		ieilen := int(ieis[offset+1] << 8 | ieis[offset+2])+3		
		iei := ieis[offset:offset+ieilen]		
		total+=ieilen
				  
		switch(iei_id) {
		case MO_HDR: 			
			imei = string(iei[7:22])
		case MO_PLD:			
			payload = string(iei[3:ieilen])			
		case MO_LOC:			
			fmt.Println(iei)		
		case MT_HDR:			
			imei = string(iei[7:22])
		case MT_PLD:
			payload = string(iei[3:ieilen])			
		case MT_CNF:
			fmt.Println("MT Confirmation message IEI found");
			fmt.Println(iei)
		}

		offset += ieilen				
		if(total==total_length) {
			break;
		}
	}
	return

}
