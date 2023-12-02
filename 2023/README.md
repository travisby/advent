# AoC 2023

This year my intention is to get as far as I can using plain ol' C.  I like to play with microcontrollers on the side and this is an area where I'm hurting.

You'll notice some of the solutions look a bit obtuse -- I _am_ trying to keep memory allocation static, and small to simulate an MCU environment.

I'll generally prefer a solution that involves multiple passes over data, provided it leads to an O(1) memory solution.  I'll try to avoid other functions in the stdlib that might `malloc` as well.

Notable exceptions for malloc:
* `fgets` and `printf` seem to perform allocations.  Since this game is all about taking an input and printing a result, I think it's fair to use these.  (We could get rid of `fgets` by adding the input as static data in the source, but that doesn't feel worth it).
