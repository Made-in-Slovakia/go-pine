package pine

type OpCode byte

// PINE operation codes
const (
	MsgRead8         OpCode = 0x00 // Read   8 bit value from memory
	MsgRead16        OpCode = 0x01 // Read  16 bit value from memory
	MsgRead32        OpCode = 0x02 // Read  32 bit value from memory
	MsgRead64        OpCode = 0x03 // Read  64 bit value from memory
	MsgWrite8        OpCode = 0x04 // Write  8 bit value from memory
	MsgWrite16       OpCode = 0x05 // Write 16 bit value from memory
	MsgWrite32       OpCode = 0x06 // Write 32 bit value from memory
	MsgWrite64       OpCode = 0x07 // Write 64 bit value from memory
	MsgVersion       OpCode = 0x08 // Returns the emulator version, example "PCSX2 v2.6.3"
	MsgSaveState     OpCode = 0x09 // Saves a savestate
	MsgLoadState     OpCode = 0x0A // Loads a savestate
	MsgTitle         OpCode = 0x0B // Returns the game title, user defined value
	MsgId            OpCode = 0x0C // Returns the game ID, a.k.a Serial, example "SCES-51607"
	MsgUuid          OpCode = 0x0D // Returns the game UUID, a.k.a. CVC, example "2F486E6F"
	MsgGameVersion   OpCode = 0x0E // Returns the game verion, example "1.00"
	MsgStatus        OpCode = 0x0F // Returns the emulator status, possible values: {0:Running,1:Paused,2:Shutdown}
	MsgUnimplemented OpCode = 0xFF // Unimplemented IPC message
)
