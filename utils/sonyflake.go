package utils

import (
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
)

var Sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		return viper.GetUint16("machine.id"), nil
	}
	Sf = sonyflake.NewSonyflake(st)
	if Sf == nil {
		panic("sonyflake not created")
	}
}
