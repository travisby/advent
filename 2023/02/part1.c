#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 162

#define RED 12
#define GREEN 13
#define BLUE 14


int main() {
	int result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	 * Here we want the sum of game IDs that are possible games
	 * A game is possible if none of its rounds use more than the expected number of cubes
	 * A game is made up of multiple rounds separated by `; `
	 * A round is made up of multiple block specifications separated by `, `, with the number of blocks per color before the color name, e.g. `3 red, 2 blue, 1 green`
	 *
	 * We tackle this by first looping over each line (a game)
	 *   then loop over each round
	 *     then looping over each color spec
	 *       and matching each spec to a color and a number
	 *         and setting a `possible` bool to false if the number is too high for that color, and breaking out early
	 * and finally adding the game ID to the result if it's possible
	 */

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		// scan the line for the game ID
		// we don't use this to also parse the game content
		// because that would require a second buffer
		// sscanf will COPY that second string
		// rather than just give us a pointer to it :(
		int id;
		assert (sscanf(buf, "Game %d: ", &id) == 1);

		// so for the game game, we'll instead just go until after the ": "
		char *game = strchr(buf, ':') + 2;


		// and now, each round is separated by a `; `
		char *round = game;
		bool possible = true;
		while (round != NULL) {
			// the round ends either at the end of te string or at the next `;`
			char *endOfRound = strchr(round, ';');

			/* a round is made up of multiple color specifications
			 * a color specification is a number of cubes followed by a color
			 * a color specification is separated from the next by a `, `
			 * so we loop until we've looked at each color spec
			 * and set possible if it's not possible
			 */
			char *colorSpec = round;
			while (colorSpec != NULL) {

				// we can use a neat little trick here
				// we don't need a whole buffer for the word `red`
				// and to store the whole word -- we can just capture the first character
				// because `r`, `b`, and `g` are all unique
				int count;
				char color;
				assert (sscanf(colorSpec, "%d %c", &count, &color) == 2);

				switch (color) {
					case 'r':
						if (count > RED) {
							possible = false;
							break;
						}
						break;
					case 'b':
						if (count > BLUE) {
							possible = false;
							break;
						}
						break;
					case 'g':
						if (count > GREEN) {
							possible = false;
							break;
						}
						break;
					default:
						assert (false);
				}

				if (!possible) {
					break;
				}

				char *nextColorSpec = strchr(colorSpec, ',');
				colorSpec = nextColorSpec != NULL ? nextColorSpec + 2 : NULL;
			}

			if (!possible) {
				break;
			}

			// continue the loop at the next round
			// +2 covers getting past the `; `
			round = endOfRound != NULL ? endOfRound + 2 : NULL;
		}

		if (possible) {
			result += id;
		}

	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
