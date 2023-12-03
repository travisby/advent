#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <ctype.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input| sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 142

#define GEAR '*'

// at a minimum, we need to look at the previous line, the current line, and the next line
// so we need 3 buffers
#define MAX_LOOKBACK 3

// findBeginningOfNumber takes a pointer to a string, along with an index of a "known digit"
// and keeps searching backwards until it either eaches the beginning of the string OR a non-digit character
char* findBeginningOfNumber(char *beginningOfString, int knownDigitIndex) {
	while (isdigit(beginningOfString[knownDigitIndex]) && knownDigitIndex >= 0) {
		knownDigitIndex--;
	}
	return beginningOfString + knownDigitIndex + 1;
}

int main() {
	int result = 0;

	char bufs[MAX_LOOKBACK][MAX_LINE_LENGTH];

	char *previous = NULL;
	char *current = NULL;
	char *next = NULL;

	// we're going to cycle through our 3 buffers
	// so we can always have the buffer used in next the round previous in current
	// and the buffer used two rounds before in previous
	int nextbuffer = 0;
	current = fgets(bufs[nextbuffer++], MAX_LINE_LENGTH, stdin);
	// instead of using a for loop where we keep setting next
	// what we really care about is `current` being something we can look at
	// so go until current is null
	while (current != NULL) {
		// we want to make sure we're ready to use the next buffer, so we'll increment it for later
		next = fgets(bufs[nextbuffer++ % MAX_LOOKBACK], MAX_LINE_LENGTH, stdin);
		// we're going to want to do something for each gear we encounter on this line
		for (char *nextGear = strchr(current, GEAR); nextGear != NULL; nextGear = strchr(nextGear + 1, GEAR)) {
			// for a number to be part of the gear ratio, it needs to be adjacent to the gear
			// that gives us possible digit locations: [previous-1, previous, previous+1, current-1, current+1, next-1, next, next+1]
			// then we need to make sure there are exactly two distinct numbers (so if previous-1, previous, and previous+1 == "123", that shouldn't count as three numbers, but one)
			
			int uniqueNumberOfStringsIndex = 0;
			char *uniqueNumberStrings[2] = {NULL, NULL};

			// we're going to iterate over each line
			// it feels a little sad/dirty that we're creating a new array here for this
			// but hoping that the compiler would just unroll this loop (and allocation)
			char *lines[3] = {previous, current, next};
			for (int i = 0; i < 3; i++) {
				char *line = lines[i];

				// ensure prev, next exist
				// if not, skip this line
				if (line == NULL) {
					continue;
				}

				// now we want to look at the line from one character to the left of the gear's index
				int leftIndex = nextGear - current - 1;
				// (if it exists, otherwise we'll just start at the gear index.. which we now know is 0)
				leftIndex = leftIndex >= 0 ? leftIndex : 0;

				// and one character to the right of the gear's index
				int rightIndex = nextGear - current + 1;
				// (if it exists, otherwise we'll just use the gear's index, whic we now know is the end of the line)
				// considering lines are NULL-terminated, this might be safe regardless, but defensive coding and all that...
				rightIndex = rightIndex <= strlen(line) ? rightIndex : strlen(line);

				// so for every adjacent character on this line, see if it's a number
				// if it is... see if it's unique
				// and finally, if it's unique, make sure we're not over our limit of two unique numbers
				for (int j = leftIndex; j <= rightIndex && uniqueNumberOfStringsIndex <= 2; j++) {
					if (!isdigit(line[j])) {
						continue;
					}

					// for deduping purposes, we look at the number as a string based on its starting position
					// (e.g. if line[j] was '2' apart of '123', we'd want `num` here to point at the '1'
					// that way, if the next loop we see the '3' we can have the same `num`,
					// and know it's not unique and skip double-counting (or triple-counting)
					char *num = findBeginningOfNumber(line, j);

					if (uniqueNumberOfStringsIndex > 0 && num == uniqueNumberStrings[uniqueNumberOfStringsIndex - 1]) {
						// we've already seen this number, so skip it
						continue;
					}

					uniqueNumberStrings[uniqueNumberOfStringsIndex++] = num;
				}
			}

			// it is only a gear if there are exactly two numbers
			if (uniqueNumberOfStringsIndex == 2) {
				// we've now proven that after accounting for dupes, there are only two numbers
				// we got a gear, hurray!
				// so add the gear ratio (multiple of the two numbers) to our result
				result += atoi(uniqueNumberStrings[0]) * atoi(uniqueNumberStrings[1]);
			}

		}

		// now that we're done processing this line, we want to move our pointers forward
		previous = current;
		current = next;
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
