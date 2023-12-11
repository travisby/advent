#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <limits.h>
#include <math.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input| sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 142
#define MAX_LINES 142

enum tileType {
	NS = '|',
	EW = '-',
	NE = 'L',
	NW = 'J',
	SW = '7',
	SE = 'F',
	NO = '.',
	ST = 'S',
};

typedef struct tile {
	int x;
	int y;
	char tileType;
} tile;
typedef struct tilemap {
	tile tiles[MAX_LINES][MAX_LINE_LENGTH];
	// assumes all lines are the same length
	// by having "." (zero) padding on the shorter lines
	int width;
	int height;
	tile *start;
} tilemap;



// it's a valid neighbor IFF:
// 0. neither is a NO tile
// 1. current connects to it (e.g. NS means it can connect to N or S)
// 2. they are adjacent for what t1 can connect to
// 3. it can accept a connection from current (e.g. NS can't connect to EW)
bool validConnection(tile t1, tile t2) {
	/*
	NS = '|',
	EW = '-',
	NE = 'L',
	NW = 'J',
	SW = '7',
	SE = 'F',
	NO = '.',
	ST = 'S', */
	if (t1.tileType == NO || t2.tileType == NO) {
		return false;
	}

	// t1 is N of t2
	if (t1.x == t2.x && t1.y +1 == t2.y) {
		return (
			t1.tileType == ST ||
			t1.tileType == NS ||
			t1.tileType == SW ||
			t1.tileType == SE
			
		) && (
			t2.tileType == ST ||
			t2.tileType == NS ||
			t2.tileType == NE ||
			t2.tileType == NW
		);
	// t1 is S of t2
	} else if (t1.x == t2.x && t1.y - 1 == t2.y) {
		return (
			t1.tileType == ST ||
			t1.tileType == NS ||
			t1.tileType == NE ||
			t1.tileType == NW
		) && (
			t2.tileType == ST ||
			t2.tileType == NS ||
			t2.tileType == SW ||
			t2.tileType == SE
		);
	// t1 is W of t2
	} else if (t1.x + 1 == t2.x && t1.y == t2.y) {
		return (
			t1.tileType == ST ||
			t1.tileType == EW ||
			t1.tileType == NE ||
			t1.tileType == SE
		) && (
			t2.tileType == ST ||
			t2.tileType == EW ||
			t2.tileType == NW ||
			t2.tileType == SW
		);
	// t1 is E of t2
	} else if (t1.x - 1 == t2.x && t1.y == t2.y) {
		return (
			t1.tileType == ST ||
			t1.tileType == EW ||
			t1.tileType == NW ||
			t1.tileType == SW
		) && (
			t2.tileType == ST ||
			t2.tileType == EW ||
			t2.tileType == NE ||
			t2.tileType == SE
		);
	}


	return false;
}

int getNeighbors(tilemap *tiles, tile current, tile *neighbors) {
	int numNeighbors = 0;
	// a tile can have up to 2 neighbors

	if (current.x > 0) {
		tile neighbor = tiles->tiles[current.y][current.x - 1];
		if (validConnection(current, neighbor)) {
			neighbors[numNeighbors++] = neighbor;
		}
	}
	assert (numNeighbors <= 2);

	if (current.x < tiles->width - 1) {
		tile neighbor = tiles->tiles[current.y][current.x + 1];
		if (validConnection(current, neighbor)) {
			neighbors[numNeighbors++] = neighbor;
		}
	}
	assert (numNeighbors <= 2);

	if (current.y > 0) {
		tile neighbor = tiles->tiles[current.y - 1][current.x];
		if (validConnection(current, neighbor)) {
			neighbors[numNeighbors++] = neighbor;
		}
	}
	assert (numNeighbors <= 2);

	if (current.y < tiles->height - 1) {
		tile neighbor = tiles->tiles[current.y + 1][current.x];
		if (validConnection(current, neighbor)) {
			neighbors[numNeighbors++] = neighbor;
		}
	}
	assert (numNeighbors <= 2);

	// (including S, which will have exactly two pipes connecting to it, and which is assumed to connect back to those two pipes).
	if (current.tileType == ST) {
		assert (numNeighbors == 2);
	}

	return numNeighbors;
}


// assumes the buf is NOT empty, and will clear it itself
// assumes visitBuf and distances has at least tiles->width * tiles->height elements
// XXX: Improvement here could be to reduce memory usage of visitBuf by / 8
//      by using bits instead of bools
// thanks to Dikjkstra for this
void shortestPath(tilemap *tiles, bool *visitBuf, int *distances) {
	tile *current = tiles->start;

	// Mark all nodes unvisited.
	for (int x = 0; x < tiles->width; x++) {
		for (int y = 0; y < tiles->height; y++) {
			visitBuf[y * tiles->width + x] = false;
		}
	}
	// Assign to every node a tentative distance value
	for (int x = 0; x < tiles->width; x++) {
		for (int y = 0; y < tiles->height; y++) {
			distances[y * tiles->width + x] = INT_MAX;
		}
	}
	// set it to zero for our initial node
	distances[current->y * tiles->width + current->x] = 0;

	while (current != NULL) {
		// for the current node, consier all of its unvisited neighbors

		// we know neighbors where 0 < neighbors <= 2
		// because one end of the "pipe" can be open
		tile neighbors[2];
		int numNeighbors = getNeighbors(tiles, *current, neighbors);

		for (int i = 0; i < numNeighbors; i++) {
			// consider all of its unvisited neighbors
			if (!visitBuf[neighbors[i].y * tiles->width + neighbors[i].x]) {
				// calculate their tentative distances through the current node
				int distance = distances[current->y * tiles->width + current->x] + 1;
				// Compare the newly calculated tentative distance to the one currently assigned to the neighbor
				// and assign it the smaller one
				distances[neighbors[i].y * tiles->width + neighbors[i].x] = fmin(distance, distances[neighbors[i].y * tiles->width + neighbors[i].x]);
			}
		}


		// Mark the current node as visited
		visitBuf[current->y * tiles->width + current->x] = true;


		// if the smallest tentative distance among the nodes in the unvisited set is infinity, then stop
		current = NULL;
		// otherwise, select the unvisited node that is marked with the smallest tentative distance
		int smallestDistance = INT_MAX;
		for (int i = 0; i < tiles->width*tiles->height; i++) {
			if (!visitBuf[i] && distances[i] < smallestDistance) {
				smallestDistance = distances[i];
				current = &tiles->tiles[i / tiles->width][i % tiles->width];
			}
		}

	}

}


int main() {
	int result = 0;

	tilemap tiles;
	tiles.height = 0;

	const char *validChars = "|-LJ7F.S";;

	int i = 0;
	char c;
	while (c = fgetc(stdin)) {
		if (feof(stdin)) {
			break;
		}
		switch (c) {
			case '\n':
				tiles.width = i;
				i = 0;
				assert (++tiles.height <= MAX_LINES);
				break;
			case ST:
				tiles.start = &tiles.tiles[tiles.height][i];
				// fallthrough
			case NS:
			case EW:
			case NE:
			case NW:
			case SW:
			case SE:
			case NO:
				tiles.tiles[tiles.height][i].x = i;
				tiles.tiles[tiles.height][i].y = tiles.height;
				tiles.tiles[tiles.height][i].tileType = c;
				i++;
				break;
			default:
				assert(false);
		}
	}
	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));


	bool visitBuf[MAX_LINES * MAX_LINE_LENGTH];
	int distances[MAX_LINES * MAX_LINE_LENGTH];

	shortestPath(&tiles, visitBuf, distances);

#ifdef PART1
	for (int i = 0; i < tiles.width * tiles.height; i++) {
		if (distances[i] != INT_MAX) {
			result = fmax(result, distances[i]);
		}
	}
#elifdef PART2
	// raycasting algorithm  thanks for the suggestion
	// https://stackoverflow.com/a/19447213
	for (int y = 0; y < tiles.height; y++) {
		int numCrossed = 0;
		for (int x = 0; x < tiles.width; x++) {
			// even-odd rule counts the number of crossings
			// crossings only happen at visits
			if (visitBuf[y * tiles.width + x]) {

				numCrossed++;
				// are we a horizontal line?
				// the more complicated part of the even-odd rule is to find where the horizontal line goes up or down
				// eventually on the left/right side
				// if they go the same direction (e.g. F-7 or L-J) then it counts as 2 crossings
				if (x < tiles.width - 1 && validConnection(tiles.tiles[y][x], tiles.tiles[y][x+1])) {
					// left goes up if [y-1][x] is a valid connection
					// but that's only valid if y-1 is in bounds too!
					bool leftGoesUp = y > 0 && validConnection(tiles.tiles[y][x], tiles.tiles[y-1][x]);

					// keep consuming horizontal connections until we can't anymore
					while (x < tiles.width && validConnection(tiles.tiles[y][x], tiles.tiles[y][x+1])) {
						x++;
					}

					// we've now consumed all the way to the right side of the horizontal line
					// which we does that face?
					bool rightGoesUp = y > 0 && validConnection(tiles.tiles[y][x], tiles.tiles[y-1][x]);
					// if they BOTH face the same direction (whether up or down), it counts as a double-crossing
					// we already added one, so add one more
					if (leftGoesUp == rightGoesUp) {
						numCrossed++;
					}
				}
			} else if (numCrossed % 2 == 1) {
				result++;
			}
		}
	}
	

#else
#error "PART1 or PART2 must be defined"
#endif
	printf("%d\n", result);
}
