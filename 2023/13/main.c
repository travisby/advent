#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <math.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 19

#define MAX_LINES_IN_PATTERN MAX_LINE_LENGTH - 2

#if !defined(PART1) && !defined(PART2)
#error "Must define PART1 or PART2"
#endif

void copyColumns(char *dest, char (*pattern)[MAX_LINES_IN_PATTERN], int pattern_rows, int column) {
		for (int i = 0; i < pattern_rows; i++) {
			dest[i] = pattern[i][column];
		}
}

int strnAlmostComp(char *a, char *b, int len) {
	int result = 0;
	for (int i = 0; i < len; i++) {
		if (a[i] != b[i]) {
			result++;
		}
	}
	return result;
}

int summary(char (*pattern)[MAX_LINES_IN_PATTERN], int pattern_rows, int pattern_columns) {
	// vertical
	char buf1[MAX_LINES_IN_PATTERN];
	char buf2[MAX_LINES_IN_PATTERN];
	for (int j = 0; j < pattern_columns - 1; j++) {
		copyColumns(buf1, pattern, pattern_rows, j);
		copyColumns(buf2, pattern, pattern_rows, j+1);

#ifdef PART2
		bool foundSmudge = false;
#endif

		bool matched = true;
		for (int i = 0; (j - i) >= 0 && pattern_columns > j+i+1; i++ ) {
			copyColumns(buf1, pattern, pattern_rows, j-i);
			copyColumns(buf2, pattern, pattern_rows, i+j+1);
#ifdef PART1
			if (strncmp(buf1, buf2, pattern_rows) != 0) {
				matched = false;
				break;
			}
#elifdef PART2
			int almostComp = strnAlmostComp(buf1, buf2, pattern_rows);
			if (almostComp > 1 ) {
				matched = false;
				break;
			} else if (almostComp == 1 && !foundSmudge) {
				foundSmudge = true;
			}
#endif
		}
#ifdef PART1
		if (matched) {
#elifdef PART2
		if (matched && foundSmudge) {
#endif
			return j + 1;
		}

	}

	// horizontal
	for (int i = 0; i < pattern_rows - 1; i++) {
		bool matched = true;
#ifdef PART2
		bool foundSmudge = false;
#endif
		for (int j = 0; (i - j) >= 0 && pattern_rows > i+j+1; j++ ) {
#ifdef PART1
			if (strncmp(pattern[i-j], pattern[i+j+1], pattern_columns) != 0) {
				matched = false;
				break;
			}
#elifdef PART2
			int almostComp = strnAlmostComp(pattern[i-j], pattern[i+j+1], pattern_columns);
			if (almostComp > 1 ) {
				matched = false;
				break;
			} else if (almostComp == 1 && !foundSmudge) {
				foundSmudge = true;
			}
#endif
		}
#ifdef PART1
		if (matched) {
#elifdef PART2
		if (matched && foundSmudge) {
#endif
			return 100 * (i+1);
		}
	}
	assert(false);
}

int main() {
	int result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	  *
	 */

	char pattern[MAX_LINES_IN_PATTERN][MAX_LINE_LENGTH-2];
	int pattern_rows = 0;
	int pattern_columns = 0;
	int i = 0;
	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		i++;
		// assert we have a full line
		char *newlineAt = strchr(buf, '\n');
		assert (newlineAt != NULL);

		if (strlen(buf) == 1) {
			result += summary(pattern, pattern_rows, pattern_columns);
			pattern_rows = 0;
			pattern_columns = 0;
		} else {
			// don't include '\n' or '\0'
			strncpy(pattern[pattern_rows++], buf, strlen(buf) - 1);
			pattern_columns = strlen(buf) - 1;
			assert(pattern_rows <= MAX_LINES_IN_PATTERN);
		}
	}

	// pattern _also_ finished
	result += summary(pattern, pattern_rows, pattern_columns);
	pattern_rows = 0;
	pattern_columns = 0;

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
