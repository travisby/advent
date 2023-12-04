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

// we need to store that we won games N+1, N+2, N+3, N+4, N+5 etc
// we PROBABLY could do some kind of linked list where the size is
// the max number of matches a card could have
// but that's a lot of work... even if we DID want to reduce memory usage
// maybe a followup
// this will be simple, and a bit more elegant and easier to grok
// knowing that NUM_GAMES just says the number of each game we have
#define NUM_GAMES 209


int main() {
	/*
	 * For part 2, the problem sounds trickier than it is
	 * 
	 * For each "match" a card has, it wins a NEW COPY of card n+1
	 * (or if it has two matches, it wins a copy of card n+1, and n+2)
	 * (or if it has three matches, it wins a copy of card n+1, n+2, and n+3, etc.)
	 *
	 * But winnings are only EVER forward-looking
	 * We never have to go backwards... so we don't need to recalculate explosion
	 * for like what game N-2 would have triggered and add new carsd for that
	 *
	 * Instead, as we go along we simply record how many copies of each card we have in `games[gameID-1]`
	 * when game N wins 3 matches, we simply increment games[gameID], games[gameID+1], games[gameID+2]
	 * FOR EACH games[gameID-1] we have (aka, multiply by games[gameID-1]
	 *
	 * we'll also separately track the total number of cards we have in `result` to make the final calculation simpler
	 *
	 * We could make this potentially use less memory by using a linked list and throwing out the games we've already processed
	 * but that's still a dynamic amount of memory, where here at least we statically know the memory won't be above NUM_GAMES*x
	 * it'll be static, even if "large" for smaller games
	 * */
	int result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	// games is an array of [gameID-1] to the number of copies of that game we have
	int games[NUM_GAMES];
	for (int i = 0; i < NUM_GAMES; i++) {
		games[i] = 0;
	}


	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		int gameID;
		assert (sscanf(buf, "Card %d:", &gameID) == 1);
		// we actually start all of our game counts at 0 (because NUM_GAMES could be larger than our input)
		// and it's possible that gameID-1 already added to the number of games we're doing for gameID
		// so to account for the original we want to increment, not explicitly set to 1
		games[gameID-1]++;

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

		int matches = 0;
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

				matches++;

				break;
			}

			line += charsread;
		}

		// so now, if we have 1 match, that means gameID+1 gets one extra card (for EVERY gameID card we have)
		// if we have 2 matches, that means gameID+2 gets one extra card (for EVERY gameID card we have)
		// etc.
		for (int i = 1; i <= matches; i++) {
			games[gameID-1+i]+= games[gameID-1];
			result += games[gameID-1];
		}

	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
