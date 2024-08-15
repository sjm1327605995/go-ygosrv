### 调用ocgcore动态库

### 平台兼容

- windows
        
      需要把 libs目录下的 ocgcore_windows.dll文件放在项目根目录下，重命名为ocgcore.dll
- linux

       暂时不支持，代码不用怎么改动。
       ygocore 动态库编译没搭建成功.然和代码设置连接位置

### api

  ```
    func GetLogMessage(pduel uintptr) []byte
    
    func GetMessage(pduel uintptr)
    
    func Process(pduel uintptr) int32 
    
    func NewCard(pduel uintptr, code uint32, owner, playerid, location, sequence, position uint8)
    
    // QueryCard  buf 长度要大于 0x2000
    func QueryCard(pduel uintptr, playerid, location, sequence uint8, queryFlag int32, buf []byte, useCache int32) 
    
    func QueryFieldCount(pduel uintptr, playerid, location uint8) int32 
    
    func QueryFieldCard(pduel uintptr, playerid, location uint8, queryFlag int32, buf []byte, useCache int32) int32 
    
    func QueryFieldInfo(pduel uintptr, buf []byte) int32 
    
    func SetResponsei(pduel uintptr, value int32) 
    
    func SetResponseb(pduel uintptr, buf []byte) 
    
    func PreloadScript(pduel uintptr, script []byte) int32 

  ```