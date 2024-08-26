# 
Try to only pass copies of data on channels. This implies that you shouldn’t pass direct pointers on channels in most cases. Passing pointers can result in multi- ple goroutines sharing memory, which can create race conditions. If you have to pass pointer references, use data structures in an immutable fashion—create them once, and don’t update them. Alternatively, pass a reference via a chan- nel, and then never use it again from the sender.

#
As much as possible, try not to mix message passing patterns with memory sharing. Using memory sharing together with message passing might create confusion as to the approach adopted in the solution.
