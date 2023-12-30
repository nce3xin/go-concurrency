# Stream

把 Channel 当作流式管道使用的方式，也就是把 Channel 看作流（Stream），提供跳过几个元素，或者是只取其中的几个元素等方法。

常见的流的方法：

1. takeN：只取流中的前 n 个数据；
2. takeFn：筛选流中的数据，只保留满足条件的数据；
3. takeWhile：只取前面满足条件的数据，一旦不满足条件，就不再取；
4. skipN：跳过流中前几个数据；
5. skipFn：跳过满足条件的数据；
6. skipWhile：跳过前面满足条件的数据，一旦不满足条件，当前这个元素和以后的元素都会输出给 Channel 的 receiver。
