# Reasons for concurrency and parallelism


To complete this exercise you will have to use git. Create one or several commits that adds answers to the following questions and push it to your groups repository to complete the task.

When answering the questions, remember to use all the resources at your disposal. Asking the internet isn't a form of "cheating", it's a way of learning.

 ### What is concurrency? What is parallelism? What's the difference?
 > Parallellism: Multiple tasks or parts of tasks are *physically* running at the same time, using multiple processing units. Parallellism is not possible in single-core CPU.
 Concurrency: Multiple tasks or parts of tasks *appears* to be running at the same time (not neccesarily on more than one processing unit). CPU-time-slicing is used: one part of a task runs, then goes to waiting state. While part one is in waiting state the CPU runs a part of the other task, and so on. 
 
 ### Why have machines become increasingly multicore in the past decade?
 > CPU clock speeds improved rapidly until around 2005, but since then it almost hasn't improved. Therefore, multicore computing (-->parallellism) is used to increase computer performance instead.
 
 ### What kinds of problems motivates the need for concurrent execution?
 (Or phrased differently: What problems do concurrency help in solving?)
 > In real-world systems many things may be happening at the same time, at random times or in random order. The software must react to these events in-real time. It is in many cases a lot easier to partition the system into concurrent software elements than to design a procedural program.
Multitasking can also speed up the execution e.g. by preventing that one activity that just waits for I/O blocks other activities.
 
 ### Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
 (Come back to this after you have worked on part 4 of this exercise)
 > As mentioned on the previous question, it is in most cases very hard to solve the mentioned problems with a procedural program, and solving it using concurrency CAN be easier. But concurrent programs are in general more difficult to understand and to write (?). Determining when and how program segments that may interact with each other should be executed is difficult
 
 ### What are the differences between processes, threads, green threads, and coroutines?
 > Process: an executing instance of an application, for example "run Microsoft Word". Used to accomplish bigger tasks.
 Thread: a part of the execution within the process. Used to accomplish small tasks. 
 ---> A process can contain multiple threads. When you start Word, the operating system creates a process and begins executing the primary thread of that process. One thread alone could be considered a lightweight process.
 
> Green thread: threads that are scheduled by a runtime library or virtual machine (VM) instead of by the underlying OS. Green threads *emulate* multithreaded environments without relying on any OS capabilities, so they work on OS that do not have thread support (they are managed in user space instead of kernel space).
 
> Coroutines: with threads, the OS switches between threads (concurrency) according to its scheduler, while the switching between coroutines is determined by the programmer and programming language. Tasks are cooperatively multitasked by pausing and resuming functions at set points, typically (but not necessarily) within a single thread. More like a procedural program
 
 ### Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 > pthread_create(): (native) threads (dependent on OS, needs POSIX-compliance)
 threading.Thread(): green threads (?)
 go: native threads
 
 ### How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 > GIL is a mutex that allows only one thread to hold the control of the Python interpreter/protects access to Python objects. It prevents multiple threads from executing at once (even with more than one CPU core). This lock is necessary mainly because CPython's memory management is based on reference counting --> not thread-safe. Therefore CPU-bound (as opposed to I/O-bound) processes are single-threaded anyways.
 
 ### With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 > The multiprocessing module: Using multiple processes instead of threads. Each Python process gets its own Python interpreter and memory space so the GIL won’t be a problem. 
 
 ### What does `func GOMAXPROCS(n int) int` change? 
 > GOMAXPROCS sets the maximum number of native threads that can execute user-level code simultaneously.
