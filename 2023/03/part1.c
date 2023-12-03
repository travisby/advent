#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <math.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input| sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 142

// grep -E -o '[^0-9.]' input | sort -u | paste -sd ''
#define SYMBOLS "#%&*+-/=@$"

#define NUMBERS "0123456789"

// at a minimum, we need to look at the previous line, the current line, and the next line
// so we need 3 buffers
#define MAX_LOOKBACK 3


bool isSymbol(char c) {
	return strchr(SYMBOLS, c) != NULL;
}

bool symbolInRange(char *begin,  char *end) {
	for (char *i = begin; i < end; i++) {
		if (isSymbol(*i)) {
			return true;
		}
	}
	return false;
}

int digits(int number) {
	return (int) (log10(number) + 1);
}

int main() {
	int result = 0;

	char bufs[MAX_LOOKBACK][MAX_LINE_LENGTH];

	char *previous = NULL;
	char *current = NULL;
	char *next = NULL;

	/* Our goal here is to find (and sum) all of the numbers who have symbols adjacent to them
	 * this can transcend its current line (meaning we need to look at the previous and next lines)
	 *
	 * so we stor 3 lines in buffers and point to them with (previous, current, next)
	 * in our loop we always operate on current -- sometimes we can look at previous, and sometimes we can look at next
	 * but everytime we can assume we can look at current
	 *
	 * we'll look for the beginning of a number in the current line
	 * then we'll determine how many characters that number is (e.g. 3 is 1 character, 123 is 3 characters)
	 * then we'll look at the characters to the left and right of that number on EACH line we have stored
	 * so, e.g. ".123." means we want to look at the [0]th and [4]th characters of each line
	 * if any of those result in SYMBOLS, we found what we needed and can add to the sum
	 */

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
		// we're going to want to do something for each number we encounter on this line
		for (char *nextNumber = strpbrk(current, NUMBERS); nextNumber != NULL; nextNumber = strpbrk(nextNumber + 1, NUMBERS)) {
			// cool, we encountered a number.  Let's make sure we get all the digits
			int value;
			assert(sscanf(nextNumber, "%d", &value) == 1);

			// to understand what indices are adjacent to this number, we need to know how many digits he is
			int length = digits(value);


			// NOW we want to operate on indices -- offsets from the line
			// because we want to use this in previous, current, AND next
			// so pointers would be insufficient
			int leftMostAdjacentIndex = nextNumber - current - 1;
			int rightMostAdjacentIndex = leftMostAdjacentIndex + length + 1;

			// ensure we stay in bounds of the line
			leftMostAdjacentIndex = leftMostAdjacentIndex < 0 ? 0 : leftMostAdjacentIndex;
			rightMostAdjacentIndex = rightMostAdjacentIndex > strlen(current) ? strlen(current) : rightMostAdjacentIndex;

			// Now is the fun part -- for previous, current, and next.. are there any symbols in the range [leftMostAdjacentIndex, rightMostAdjacentIndex]?
			if (
				// symbolInRange is [begin, end) so we add + 1 to rightMostAdjacentIndex
				(previous != NULL && symbolInRange(previous + leftMostAdjacentIndex, previous + rightMostAdjacentIndex+1)) ||
				// we already know current isn't NULL, but symmetry looks good
				(current != NULL && symbolInRange(current + leftMostAdjacentIndex, current + rightMostAdjacentIndex+1)) ||
				(next != NULL && symbolInRange(next + leftMostAdjacentIndex, next + rightMostAdjacentIndex+1))

			) {
				// we have an adjacency match!
				// perform the whole thing we're here to do... add up all the numbers with adjacent symbols
				result += value;
			}

			// we want to make sure with the number "123", we don't start processing "2" next round
			nextNumber += length;
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
