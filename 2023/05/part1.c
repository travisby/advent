#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <ctype.h>
#include <stdlib.h>
#include <stdbool.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 212

// these are the max number of seeds we'll care about... anymore and we'll break
#define MAX_SEEDS 20

int main() {
	uint result = 0;

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/* 
	 * Part 1's input is a lot more complicated than previous
	 * we're not dealing with "something new and independent each line"
	 * instead, there's a few different types of lines:
	 * * a seed initializer
	 * * empty lines that separate...
	 * * the mapping header
	 * * individual mappings
	 *
	 * Luckily, we don't need to "look back" over any of these mappings
	 * so we can still keep processing line-by-line without much additional storage
	 *
	 * We keep track of each seed's value in two arrays (current and next)
	 * for each individual mapping we see if any current[i] exists in the mapping range
	 * and if so, copy it to next[i]
	 *
	 * once we get to the end (the beginning of the _next_ mapping) we reset "current"
	 * to be what's in next, and are happy to reuse next for the next mapping
	 *
	 * we "copy" rather than swapping buffers back & forth because we rely on "current"
	 * having the latest data in it, so we don't need to trakc which seeds WEREN't mapped
	 * Because out of all of this.. not every seed is in a mapping, and if not next[i]=current[i]
	 * so memcpy gives us that for free space, if not free time
	 *
	 * at the end, we just find the smallest value in current for our answer
	 */

	// first, we want to pull the first line containing all the seeds
	assert(fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);

	// assert we have a full line and aren't missing any input
	assert (strchr(buf, '\n') != NULL);

	// double buffer for "current, next" seeds
	uint seeds[2][MAX_SEEDS];

	uint *current = seeds[0];
	uint *next = seeds[1];

	int num_seeds = 0;
	// strtok_r requires a saveptr
	char *saveptr;
	for (char *token = strtok_r(buf, " ", &saveptr); token != NULL; token = strtok_r(NULL, " ", &saveptr)) {
		// is this a seed?  since token != NULL we know it's safe to look at token[0]
		// if it's a number, we've got a seed starting at token!
		if (isdigit(token[0])) {
			// seeds only holds up to MAX_SEEDS seeds... so error if we go above that
			assert (num_seeds < MAX_SEEDS);
			next[num_seeds++] = strtoll(token, NULL, 10);
		}
	}


	memcpy(current, next, sizeof(seeds[0]));

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		if (strlen(buf) == 1) {
			// ignore the newline-only lines
			continue;
		} else if (strchr(buf, ':') != NULL) {
			memcpy(current, next, sizeof(seeds[0]));
			continue;
		}
		// we have another mapping
		uint destination, source, length;
		assert (sscanf(buf, "%d %d %d", &destination, &source, &length) == 3);

		// now, check each seed if it shoudl be mapped
		for (int i = 0; i < num_seeds; i++) {
			int offset = destination - source;
			if (current[i] >= source && current[i] < source + length) {
				// map it!
				next[i] = current[i] + offset;
			}
		}
	}

	// cool, now find the smallest number in next
	result = -1;
	for (int i = 0; i < num_seeds; i++) {
		if (next[i] < result) {
			result = next[i];
		}
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}
