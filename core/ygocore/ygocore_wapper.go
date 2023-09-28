package ygocore

/*
#include "ocgcore.h"
*/
import "C"
import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

// 初始化一个随机因子
func init() {

	rand.Seed(time.Now().UnixNano())
}

const scriptMaxSize = 0x20000

//export goScriptReader
func goScriptReader(scriptName *C.char, slen *C.int) *C.uchar {
	// 将C字符串转换为Go字符串

	*slen = 0
	s := C.GoString(scriptName)
	fmt.Println(s)
	// 调用适当的函数读取脚本内容
	data, _ := os.ReadFile(C.GoString(scriptName))
	if len(data) == 0 {
		// 处理错误
		return (*C.uchar)(nil)
	}

	// 将数据长度设置到slen指针
	*slen = C.int(len(data))
	// 创建C字节数组并将数据复制到其中
	return (*C.uchar)(C.CBytes(data))

}

//export goMessageHandler
func goMessageHandler(data unsafe.Pointer, size C.uint32_t) {
	fmt.Println("goMessageHandler")
	// 处理消息
}

//export goCardReader
func goCardReader(cardID C.uint32_t, data *C.card_data) C.uint32_t {

	//TODO 这里进行了内存拷贝需要重新操作下
	var (
		dataC CardDataC
	)
	if getDataForCore(uint32(cardID), &dataC) {
		data.code = C.uint32_t(dataC.code)
		data.alias = C.uint32_t(dataC.alias)
		data.setcode = C.uint64_t(dataC.setcode)
		data._type = C.uint32_t(dataC.typ)
		data.level = C.uint32_t(dataC.level)
		data.attribute = C.uint32_t(dataC.attribute)
		data.race = C.uint32_t(dataC.race)
		data.attack = C.int32_t(dataC.attack)
		data.defense = C.int32_t(dataC.defense)
		data.lscale = C.uint32_t(dataC.lscale)
		data.rscale = C.uint32_t(dataC.rscale)
		data.link_marker = C.uint32_t(dataC.link_marker)

	} else {
		data.code = C.uint32_t(0)
		data.alias = C.uint32_t(0)
		data.setcode = C.uint64_t(0)
		data._type = C.uint32_t(0)
		data.level = C.uint32_t(0)
		data.attribute = C.uint32_t(0)
		data.race = C.uint32_t(0)
		data.attack = C.int32_t(0)
		data.defense = C.int32_t(0)
		data.lscale = C.uint32_t(0)
		data.rscale = C.uint32_t(0)
		data.link_marker = C.uint32_t(0)
	}
	fmt.Printf("%+v\n", dataC)
	return 0

}
func InitCore() {
	C.set_script_reader(C.script_reader(C.goScriptReader))
	C.set_message_handler(C.message_handler(C.goMessageHandler))
	C.set_card_reader(C.card_reader(C.goCardReader))
	t := CreateGame()
	fmt.Println(t)
}
func CreateGame() int64 {
	seed := rand.Int31()
	pDuel := C.create_duel(C.int(seed))
	return int64(pDuel)

}
