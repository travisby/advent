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
			
			// so first, it's a bit gross but we're learning C here so this is the best I know!
			// let's gather all possible numbers into an array... we have up to 8 possibilities for adjacency
			char *numberStrings[8];
			// and we'll want to track where to insert into next
			int numberOfStringNextIndex = 0;

			// for each row, we want to look at the index to the left of the gear, the right, and the same index as the gear (for prev&next, not necessary for current)
			int leftIndex = nextGear - current - 1;
			int middleIndex = nextGear - current;
			int rightIndex = nextGear - current + 1;

			// this is also gross... hopefully a refactor could simplify this
			// but for right now, for each of those 8 possible matches:
			// see if the row exists, see if the index is valid, and see if the character is a digit
			// if so, add a string pointer to our array that covers WHERE THE NUMBER BEGINS IN THE STRING
			// that will allow us to dedupe later to make sure we're only counting distinct numbers
			if (previous != NULL && leftIndex >= 0 && isdigit(previous[leftIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(previous, leftIndex);
			}
			if (previous != NULL && isdigit(previous[middleIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(previous, middleIndex);
			}
			if (previous != NULL && rightIndex <= strlen(previous) && isdigit(previous[rightIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(previous, rightIndex);
			}

			if (current != NULL && leftIndex >= 0 && isdigit(current[leftIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(current, leftIndex);
			}
			if (current != NULL && isdigit(current[middleIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(current, middleIndex);
			}
			if (current != NULL && rightIndex <= strlen(current) && isdigit(current[rightIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(current, rightIndex);
			}

			if (next != NULL && leftIndex >= 0 && isdigit(next[leftIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(next, leftIndex);
			}
			if (next != NULL && isdigit(next[middleIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(next, middleIndex);
			}
			if (next != NULL && rightIndex <= strlen(next) && isdigit(next[rightIndex])) {
				numberStrings[numberOfStringNextIndex++] = findBeginningOfNumber(next, rightIndex);
			}

			// okay, now we have an array of all of the strings of numbers
			// now we need to dedupe -- it's possible for the string "123\n.*." and 123 would appear in here three separate times
			// it should only count once
			int unique = 0;
			int uniqueNumbers[2] = {0, 0};

			// we're going to see if this is a gear by checking if there's exactly two unique numbers
			// so this code here is to gather EITHER the count of unique numbers, OR to say "it's greater than 2, bomb out"
			for (int i = 0; i < numberOfStringNextIndex; i++) {
				if (unique == 0) {
					// fallthrough, to our base condition
					// where we'll add to the array and increment unique count
				} else if (numberStrings[i] == numberStrings[i-1]) {
					// this isn't unique, so we can skip it
					continue;
				} else if (unique > 1) {
					// we're already larger than uniqueNumbers can be
					// and we know we're not a gear
					// so just mark unique as large and quit early
					unique++;
					break;
				}
				uniqueNumbers[unique++] = atoi(numberStrings[i]);
			}

			// it is only a gear if there are exactly two numbers
			if (unique == 2) {
				// we've now proven that after accounting for dupes, there are only two numbers
				// we got a gear, hurray!
				// so add the gear ratio (multiple of the two numbers) to our result
				result += uniqueNumbers[0] * uniqueNumbers[1];
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
