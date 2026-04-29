package pine

type opCode byte

const (
	msgRead8         opCode = 0x00 // Read   8 bit value from memory
	msgRead16        opCode = 0x01 // Read  16 bit value from memory
	msgRead32        opCode = 0x02 // Read  32 bit value from memory
	msgRead64        opCode = 0x03 // Read  64 bit value from memory
	msgWrite8        opCode = 0x04 // Write  8 bit value from memory
	msgWrite16       opCode = 0x05 // Write 16 bit value from memory
	msgWrite32       opCode = 0x06 // Write 32 bit value from memory
	msgWrite64       opCode = 0x07 // Write 64 bit value from memory
	msgVersion       opCode = 0x08 // Returns the emulator version, example "PCSX2 v2.6.3"
	msgSaveState     opCode = 0x09 // Saves a savestate
	msgLoadState     opCode = 0x0A // Loads a savestate
	msgTitle         opCode = 0x0B // Returns the game title, user defined value
	msgId            opCode = 0x0C // Returns the game ID, a.k.a Serial, example "SCES-51607"
	msgUuid          opCode = 0x0D // Returns the game UUID, a.k.a. CVC, example "2F486E6F"
	msgGameVersion   opCode = 0x0E // Returns the game verion, example "1.00"
	msgStatus        opCode = 0x0F // Returns the emulator status, possible values: {0:Running,1:Paused,2:Shutdown}
	msgUnimplemented opCode = 0xFF // Unimplemented IPC message
)
