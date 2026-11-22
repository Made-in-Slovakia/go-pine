package pine

import (
	"testing"
)

func TestDeserialize(t *testing.T) {
	var tests = []struct {
		name   string
		opCode opCode
		input  []byte
		want   any
	}{
		// NOTE: These test are not testing specific serialization,
		// they are focused on propher returned type from deserialization
		// or on fact that for some opCode is response not serialized and
		// default empty string is returned.
		{"msgRead8 opCode", msgRead8, []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}, uint8(80)},
		{"msgRead16 opCode", msgRead16, []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}, uint16(17232)},
		{"msgRead32 opCode", msgRead32, []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}, uint32(1481851728)},
		{"msgRead64 opCode", msgRead64, []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}, uint64(3636129149750559568)},
		{"msgWrite8 opCode", msgWrite8, []byte{80, 67, 83, 88, 50}, ""},
		{"msgWrite16 opCode", msgWrite16, []byte{80, 67, 83, 88, 50}, ""},
		{"msgWrite32 opCode", msgWrite32, []byte{80, 67, 83, 88, 50}, ""},
		{"msgWrite64 opCode", msgWrite64, []byte{80, 67, 83, 88, 50}, ""},
		{"msgVersion opCode", msgVersion, []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"msgSaveState opCode", msgSaveState, []byte{80, 67, 83, 88, 50}, ""},
		{"msgLoadState opCode", msgLoadState, []byte{80, 67, 83, 88, 50}, ""},
		{"msgTitle opCode", msgTitle, []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"msgId opCode", msgId, []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"msgUuid opCode", msgUuid, []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"msgGameVersion opCode", msgGameVersion, []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"msgStatus opCode", msgStatus, []byte{80, 67, 83, 88, 50}, uint32(1481851728)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := deserialize(tt.opCode, tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
			if err != nil {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

func TestDeserializeErrors(t *testing.T) {
	var tests = []struct {
		name   string
		opCode opCode
		input  []byte
	}{
		{"msgRead8 opCode with nil input", msgRead8, nil},
		{"msgRead8 opCode with 0 bytes input", msgRead8, []byte{}},
		{"msgRead16 opCode with 1 byte input", msgRead16, []byte{80}},
		{"msgRead32 opCode with 3 bytes input", msgRead32, []byte{80, 67, 83}},
		{"msgRead64 opCode with 7 bytes input", msgRead64, []byte{80, 67, 83, 88, 50, 32, 118}},
		{"msgUnimplemented opCode", msgUnimplemented, []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := deserialize(tt.opCode, tt.input)
			if err == nil {
				t.Errorf("got %v, want error", ans)
			}
		})
	}
}

func TestToUint8(t *testing.T) {
	input := []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}
	var want uint8 = 80
	ans, _ := toUint8(input)
	if ans != want {
		t.Errorf("got %v, want %v", ans, want)
	}
}

func TestToUint16(t *testing.T) {
	input := []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}
	var want uint16 = 17232
	ans, _ := toUint16(input)
	if ans != want {
		t.Errorf("got %v, want %v", ans, want)
	}
}

func TestToUint32(t *testing.T) {
	input := []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}
	var want uint32 = 1481851728
	ans, _ := toUint32(input)
	if ans != want {
		t.Errorf("got %v, want %v", ans, want)
	}
}

func TestToUint64(t *testing.T) {
	input := []byte{80, 67, 83, 88, 50, 32, 118, 50, 46, 54, 46, 51}
	var want uint64 = 3636129149750559568
	ans, _ := toUint64(input)
	if ans != want {
		t.Errorf("got %v, want %v", ans, want)
	}
}

func TestToString(t *testing.T) {
	var tests = []struct {
		name  string
		input []byte
		want  string
	}{
		{"empty string should be returned for nil", nil, ""},
		{"deserialized string should be returned", []byte{80, 67, 83, 88, 50}, "PCSX2"},
		{"null-terminated string should be deserialized", []byte{80, 67, 83, 88, 50, 0}, "PCSX2"},
		{"null-terminated string should be deserialized and garbage should be removed", []byte{80, 67, 83, 88, 50, 0, 80, 67, 83, 88, 50}, "PCSX2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := toString(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
