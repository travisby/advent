#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 142
#define MAX_COLUMS MAX_LINE_LENGTH - 2
#define MAX_ROWS 140

#ifdef PART1
#define EMPTY_MODIFIER 1
#elifdef PART2
#define EMPTY_MODIFIER 1000000-1
#else
#error "Must define PART1 or PART2"
#endif


int main() {
	int64_t result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	  * TODO
	 */

	char image[MAX_ROWS*MAX_COLUMS];
	int imageSize = 0;
	int columns = 0;
	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		char *newline = strchr(buf, '\n');
		assert (newline != NULL);
		columns = newline - buf;
		assert (strncpy(image + imageSize, buf, columns) == image + imageSize);
		imageSize += columns;
	}
	int rows = imageSize / columns;


	char *galaxies[MAX_ROWS*MAX_COLUMS];
	int numGalaxies = 0;
	for (int i = 0; i < imageSize; i++) {
		if (image[i] == '#') {
			galaxies[numGalaxies++] = image + i;
		}
	}

	int emptyRows[MAX_ROWS];
	int numEmptyRows = 0;
	for (int i = 0; i < imageSize; i+= columns) {
		if (memchr(&image[i], '#', columns * sizeof(image[0])) == NULL) {
			emptyRows[numEmptyRows++] = i / columns;
		}
	}

	int emptyColumns[MAX_COLUMS];
	int numEmptyColumns = 0;
	for (int i = 0; i < columns; i++) {
		bool found = false;
		for (int j = 0; j < rows; j++) {
			if (image[j * columns + i] == '#') {
				found = true;
				break;
			}
		}
		if (!found) {
			emptyColumns[numEmptyColumns++] = i;
		}
	}

	for (int i = 0; i < numGalaxies; i++) {
		int startRow = (galaxies[i] - image) / columns;
		int startColumn = (galaxies[i] - image) % columns;
		for (int j = i + 1; j < numGalaxies; j++) {
			int64_t miniresult = 0;
			int endRow = (galaxies[j] - image) / columns;
			int endColumn = (galaxies[j] - image) % columns;
			
			miniresult += abs(endRow - startRow) + abs(endColumn - startColumn);

			// how many empty rows and columns are between these two galaxies?
			// each of those counts as two steps
			for (int k = 0; k < numEmptyRows; k++) {
				if ((emptyRows[k] > startRow && emptyRows[k] < endRow) || ((emptyRows[k] > endRow && emptyRows[k] < startRow))) {
					miniresult += EMPTY_MODIFIER;
				}
			}
			for (int k = 0; k < numEmptyColumns; k++) {
				if ((emptyColumns[k] > startColumn && emptyColumns[k] < endColumn) || (emptyColumns[k] > endColumn && emptyColumns[k] < startColumn)) {
					miniresult += EMPTY_MODIFIER;
				}
			}
			result += miniresult;
		}
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%lld\n", result);
}
