package nes

import "log"

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorQuad       = 2
)

type Cartridge struct {
	PRG      []byte // PRG-ROM banks
	CHR      []byte // CHR-ROM banks
	SRAM     []byte // Save RAM
	Mapper   byte   // mapper type
	Mirror   byte   // mirroring mode
	Battery  byte   // battery present
	prgBanks int
	chrBanks int
	prgBank1 int
	prgBank2 int
	chrBank  int
}

func NewCartridge(prg, chr []byte, mapper, mirror, battery byte) *Cartridge {
	prgBanks := len(prg) / 0x4000
	chrBanks := len(chr) / 0x2000
	prgBank1 := 0
	prgBank2 := prgBanks - 1
	chrBank := 0
	sram := make([]byte, 0x2000)
	cartridge := Cartridge{
		prg, chr, sram, mapper, mirror, battery,
		prgBanks, chrBanks, prgBank1, prgBank2, chrBank}
	return &cartridge
}

func (c *Cartridge) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		index := c.chrBank*0x2000 + int(address)
		return c.CHR[index]
	case address >= 0xC000:
		index := c.prgBank2*0x4000 + int(address-0xC000)
		return c.PRG[index]
	case address >= 0x8000:
		index := c.prgBank1*0x4000 + int(address-0x8000)
		return c.PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return c.SRAM[index]
	default:
		log.Fatalf("unhandled cartridge read at address: 0x%04X", address)
	}
	return 0
}

func (c *Cartridge) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		index := c.chrBank*0x2000 + int(address)
		c.CHR[index] = value
	case address >= 0x8000:
		c.prgBank1 = int(value)
	case address >= 0x6000:
		index := int(address) - 0x6000
		c.SRAM[index] = value
	default:
		log.Fatalf("unhandled cartridge write at address: 0x%04X", address)
	}
}
