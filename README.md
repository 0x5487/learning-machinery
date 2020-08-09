# learning-machinery

### Issues

1. `worker` 被 shutdown 的時候會把手上來在處裡中的任務完成 (gracefully shutdown)，但如果像停電這種場景服務直接 crash 的話，那這個任務就消失了，不會重新執行

   refs: https://github.com/RichardKnop/machinery/issues/586

### Questions

1. 如果 task 沒有註冊到worker 要如何發現？