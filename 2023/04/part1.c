#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <ctype.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 118

// later we will use sprintf to convert a number to a string
// with a leading space (to avoid matching "142" when we're looking for "42")
// so we need the largest number of digits + 2 characters to store that string
#define NUMBER_AS_STRING_SIZE 4


int main() {
	int result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	 * We have a few card games in the format:
	 * Card {id}: {winning number} {winning number}... | {number} {number}...
	 * Each winning number match on the RHS of | is worth one point
	 * and we want the total for all games
	 *
	 * for each game:
	 * we take a lower-memory approach and for each winning number search the string for more of that winning number
	 * that means we'll scan each line num_winning_number times
	 *
	 * the alternative is we could create an array of all the winning numbers and read the rest of the line and for each number see if it fits
	 * but that would be more memory intensive -- and we're aiming for low memory usage
	 *
	 * taking this approach, our biggest concern is making sure if our winning number is "42" that we don't match "142" or "1423"
	 * so we'll look for a leading space (everything has a leading space), and THEN check if there's another digit after it
	 *
	 * finally, the "result" is going to be the sum of all game winnings.  A game winning is matchPoints^2
	 */

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		// first, spit out the `Card N:`

		char *line = strchr(buf, ':') + 1;
		assert (line != NULL);

		// the winning numbers stop just before the pipe character
		char *winningNumbersStopsAt = strchr(line, '|') - 2;
		assert (winningNumbersStopsAt != NULL);

		// and the game numbers begin just after the pipe character
		// based on the alg we specified earlier to detect "whole" numbers
		// rather than just parts of numbers, we also want to preserve the space before the first game number
		// so we use +2 to get the pipe
		char *gameNumbersStartAt = winningNumbersStopsAt + 2;

		int matchPoints = 0;
		// for each winning number...
		while (line < winningNumbersStopsAt) {
			// we want to find numbers, but we only want to operate on strings here
			// so this is only to get scanf a target after it finds a number for us
			// NOTE: we took the scanf %d approach here because we didn't want to deal with
			// how scanf doesn't return a handle to the string it found for %s, but COPIES
			// into a new buffer and we wanted to avoid that second buffer
			// but we ended up having to use one ANYWAY to support our "leading space" approach
			// we could end up refactoring this now.
			int winningNumber;

			// we want to track how far we read into the string, so we can advance the pointer
			// so we can then scan for the _next_ number
			int charsread = 0;

			// the return value of sscanf is how many items were successfully parsed
			// that does not include %n because that's not something parsed from the string
			assert(sscanf(line, " %d%n", &winningNumber, &charsread) == 1);


			// grrrr, I didn't want to do this
			// but to search for the string version of a number 
			// we need to convert it to a string, and that uses sprintf
			// which needs an additional buffer :/
			char winningNumberAsString[NUMBER_AS_STRING_SIZE];

			// ALL numbers will have a space before them
			// but only SOME numbers will have a space after
			// so we ensure this isn't us looking for "1" in "142" by printing with a space now
			// and LATER we'll have to do our own check for the potential space after
			int snprintf_result = snprintf(winningNumberAsString, NUMBER_AS_STRING_SIZE, " %d", winningNumber);
			// > Notice that only when this returned value is non-negative and less than n, the string has been completely written.
			// if this assert fails, we need a larger NUMBER_AS_STRING_SIZE
			assert (snprintf_result > 0 && snprintf_result < NUMBER_AS_STRING_SIZE);

			for (char *possibleWinningNumberAt = strstr(gameNumbersStartAt, winningNumberAsString); possibleWinningNumberAt != NULL; possibleWinningNumberAt = strstr(possibleWinningNumberAt+1, winningNumberAsString)) {
				// our `winningNumberAsAString` ensures that there's digit BEFORE the number
				// but we need to validate that ourselves too
				// since C strings are NULL terminated, AND we know all of our strings contain a `\n` (note the assert at the beginning of the fgets loop)
				// we can SAFELY check the character after the end of our string and see if it's also a digit
				if (isdigit(*(possibleWinningNumberAt + strlen(winningNumberAsString)))) {
					// if it is, we're only a puny substring of the number we're interested in, and should continue looking
					continue;
				}

				// the points for matchPoints is {0, 1, 2, 4....}
				if (!matchPoints) {
					matchPoints = 1;
				} else {
					matchPoints<<=1;
				}

				break;
			}

			line += charsread;
		}
		result += matchPoints;
		printf(gameNumbersStartAt);
	}

	//

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
