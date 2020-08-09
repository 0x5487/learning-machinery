# learning-machinery

### Issues

1. `worker` 被 shutdown 的時候會把手上來在處裡中的任務完成 (gracefully shutdown)，但如果像停電這種場景服務直接 crash 的話，那這個任務就消失了，不會重新執行 (redis broker)

   refs: https://github.com/RichardKnop/machinery/issues/586

1. workflow 的狀態無法被觀察到目前執行到哪個階段, 也缺少一個視覺的工具才查詢

### Chain

1. 在傳遞參數到後面每一個 task, 其實很不直覺, 例如第一個 order 服務產生出order id, 要一路往後面的 task 傳,　非常麻煩


### Questions

1. 如果 task 沒有註冊到worker 要如何發現？