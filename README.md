#### MAP实现key不存在时 get操作等待 直到key存在或者超时

#### 用法：
    //创建waitingMap
    var wMap = waitingmap.NewMap()
    //赋值
    wMap.Wt(key, value)
    //获取值
    getVal := wMap.Rd(key, timeout)
