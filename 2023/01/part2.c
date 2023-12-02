#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <assert.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 61

// this is the number of single digit numbers....
// this is because getting the number of elements in an array is not actually simple
// in C because `sizeof` gives you the byte size of the array, not the number of elements
// so let's just create a defines for it
#define NUM_NUMBERS 10

// how many possible ways do we have to say one number?
// this is because getting the number of elements in an array is not actually simple
// in C because `sizeof` gives you the byte size of the array, not the number of elements
// so let's just create a defines for it
#define NUMBER_OF_STRINGS_FOR_A_NUMBER 2
char* STRINGS_OF_NUMBER[10][NUMBER_OF_STRINGS_FOR_A_NUMBER] = {
	{"0", "zero"},
	{"1", "one"},
	{"2", "two"},
	{"3", "three"},
	{"4", "four"},
	{"5", "five"},
	{"6", "six"},
	{"7", "seven"},
	{"8", "eight"},
	{"9", "nine"},
};

// an encounter is what and where we found a number
struct encounter {
	// our MAX_LINE_LENGTH is below 128, so char is safe here
	char index;

	// we use a character to help build up strings
	// e.g. we want "1" and "0" for two encounters to make the number 10 easier
	char value;
};

int main() {
	int result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	 * we want to find the first and last number in each line
	 * convert that into a base10 number (e.g. `1` and `2` is `12`)
	 * and sum them all
	 * where "number" is fuzzy -- meaning either `1` or `one`
	 * as an edge case, `4twoone` needs to turn into `421` -> `41`
	 * we can't just replace all the strings with their numbers, as some words
	 * might serve double duty
	 *
	 * So our approach here will be to sum the results of:
	 * for each line
	 *   create a "first" and "last" encounter variable
	 *     where we can store the first number we've found (and where we found it, i.e. its index)
	 *     set first to a very high index, so any real match will take over
	 *     and set last to a very low index, so any real match will take over
	 *   for [0..9]
	 *     look for both the integer and string version of the number (e.g. "0" and "zero") in the line
	 *     if we find a match check if its  either sooner than first (index wise) or later than last, and if so replace the encounter
	 *     NEXT, we can't just stop here.. its possible that there's two ocurrences of this number, and `strstr` our substring match function only returns the first occurrence
	 *       so repeat the above, but starting AFTER $index (use line[$index+1] as the new "line")
	 *       repeat until there's no more matches
	 *     we now have two encounters, simply put them together as a string and conver that to an int with atoi
	 */

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		assert (strchr(buf, '\n') != NULL);

		// make first & last outside of the expect bounds
		// so we'll always get "real" values for the first encountered numbers
		struct encounter first = {MAX_LINE_LENGTH, 0};
		struct encounter last = {0, 0};


		// for [0..9]
		for (int i = 0; i < NUM_NUMBERS; i++) {
			// for "0" and "zero"
			for (int j = 0; j < NUMBER_OF_STRINGS_FOR_A_NUMBER; j++) {
				// for each find of the number in the line
				for (char* substr = strstr(buf, STRINGS_OF_NUMBER[i][j]); substr != NULL; substr = strstr(substr+1, STRINGS_OF_NUMBER[i][j])) {
					// we do pointer math to determine the index of the match based on buf
					int index = substr - buf;
					if (index <= first.index) {
						first.index = index;
						// [0] holds the character representation of the number
						// and we only want a character now, so [0] that!
						first.value = STRINGS_OF_NUMBER[i][0][0];
					}
					if (index >= last.index) {
						last.index = index;
						// [0] holds the character representation of the number
						// and we only want a character now, so [0] that!
						last.value = STRINGS_OF_NUMBER[i][0][0];
					}
				}
			}
		}

		// convert our encounters into a string and then atoi() it
		// don't forget the NULL terminator for strings!
		result += atoi((char[]){first.value, last.value, '\0'});
	}

	assert (!ferror(stdin));
	
	printf("%d\n", result);
}
