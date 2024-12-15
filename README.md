# Gempool

A memory pool implementation in Golang whose main goal is to be as fast as possible.
This means that the code may bypass Golang safety constraints with regards to memory,
goroutines, scheduling, etc.

## Rules
- Every decision in the system will be taken after meticulous benchmarking
- Performance will be more important than code safety and standards
- Multiple CPU architectures and OSes don't need to be supported

## Possible areas of improvement
- Manual memory allocation
- CPU affinity
    - Using CPU Set affinity
    - Using taskset and running different golang processes communicating over mmap
- io_uring vs epoll