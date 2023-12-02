#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <math.h>

/* NOTE: because we use `math.h` we need to compile with `-lm` to link the math library
 * e.g. `gcc part2.c -lm`
 */

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
	 * Here we want the sum of the "power" of each game
	 * The power of a game is the multiplication of each of the minimum cubes needed to make the game possible
	 *
	 * A game is made up of multiple rounds separated by `; `
	 * A round is made up of multiple block specifications separated by `, `, with the number of blocks per color before the color name, e.g. `3 red, 2 blue, 1 green`
	 *
	 * We tackle this by first looping over each line (a game)
	 *   creating variables to store the maximum red, blue, green cubes used in the game
	 *   then loop over each round
	 *     then looping over each color spec
	 *       and matching each spec to a color and a number
	 *         and setting the r, b, or g variable if the number is greater than the current max
	 * and finally calculating + adding the power to the result
	 */

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		int red = 0;
		int blue = 0;
		int green = 0;

		// so for the game content, we'll instead just go until after the ": "
		char *game = strchr(buf, ':') + 2;


		// and now, each round is separated by a `; `
		char *round = game;
		while (round != NULL) {
			// the round ends either at the end of te string or at the next `;`
			char *endOfRound = strchr(round, ';');

			/* a round is made up of multiple color specifications
			 * a color specification is a number of cubes followed by a color
			 * a color specification is separated from the next by a `, `
			 * so we loop until we've looked at each color spec
			 * and potentially set the max for each color
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
						red = fmax( red, count);
						break;
					case 'b':
						blue = fmax(blue, count);
						break;
					case 'g':
						green = fmax(green, count);
						break;
					default:
						assert (false);
				}

				char *nextColorSpec = strchr(colorSpec, ',');
				colorSpec = nextColorSpec != NULL ? nextColorSpec + 2 : NULL;
			}

			// continue the loop at the next round
			// +2 covers getting past the `; `
			round = endOfRound != NULL ? endOfRound + 2 : NULL;
		}

		result += red * blue * green;

	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
