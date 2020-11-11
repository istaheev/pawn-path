# Description

The algorithm uses depth-first search with the Warnsdorff's rule applied. It
starts from the start position and tries to build the longest path possible. If
the path doesn't cover all the cells on the board then the algorithm backtracks
several moves back to the position where alternative moves are available and
then tries to build a new path from that point and so on. When there are several
moves available from some cell the algorithm prioritizes moves with less next
moves available (the Warnsdorff's rule).

# Running

The application has no external dependencies, so it is just enough to clone it
and run

    go run main.go

The board size and the starting point are set in the beginning of `main`.

# Implementation details

The initial approach was just plain depth-first search which showed itself to
take a lot of time on big boards (bigger than 6x6). Then I read about
the Warnsdorff's rule and implemented it which gave a huge speed up for big
boards (~400ms for 10x10 boards on my machine). However, there are cases which
still take a lot of time, for example when we start from (1,0). I don't know if
there is a path exists at all in that case (so the algorithm has to go through
all possible combinations of moves) or it is just a in some way special starting
point.

What could be done to improve the worst case? Theoretically we can parallelize
calculations (e.g. calculate paths on different parts of the board and then
stitch them). We can also minimize search space by removing moves which build
symmetrical paths.