#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <ctype.h>
#include <stdlib.h>
#include <stdbool.h>
#include <limits.h>

/*
 * I'm not happy with this brute force solution
 * down below I document my in-progress non brute force solution
 *
 * I want to track ranges of seeds rather than individual seeds
 * with the idea that when we get a mapping, and it happens to be inside a seed range
 * we split off into two new ranges
 *
 * so we'll still need some dynamic memory allocation
 * but not the ~15G it takes with the brute force solution
 */

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 212

// these are the max number of seeds we'll care about... anymore and we'll break
#define MAX_SEEDS 2000000000

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
	uint *seeds[2];
	seeds[0] = malloc(MAX_SEEDS * sizeof(uint));
	seeds[1] = malloc(MAX_SEEDS * sizeof(uint));

	uint *current = seeds[0];
	uint *next = seeds[1];

	int num_seeds = 0;

	uint start, length, bytesRead;;
	for (char *partOfString = buf + strlen("seeds: "); sscanf(partOfString, "%u %u%n", &start, &length, &bytesRead) == 2; partOfString += bytesRead + 1) {
		for (int i = start; i < start + length; i++) {
			assert (num_seeds < MAX_SEEDS);
			current[num_seeds++] = i;
		}
	}

	/*
	printf("Number of seeds: ud\n", num_seeds);
	for (int i = 0; i < num_seeds; i++) {
		printf("%u ", current[i]);
	}
	printf("\n");
	*/

	memcpy(next, current, sizeof(uint) * num_seeds);

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		if (strlen(buf) == 1) {
			// ignore the newline-only lines
			/*
			for (int i = 0; i < num_seeds; i++) {
				printf("%u ", current[i]);
			}
			printf("\n");
			*/
			continue;
		} else if (strchr(buf, ':') != NULL) {
			// printf("Found mapping header: %s", buf);
			memcpy(current, next, sizeof(uint) * num_seeds);
			continue;
		}
		// we have another mapping
		uint destination, source, length;
		assert (sscanf(buf, "%u %u %u", &destination, &source, &length) == 3);

		// now, check each seed if it shoudl be mapped
		for (int i = 0; i < num_seeds; i++) {
			int offset = destination - source;
			if (current[i] >= source && current[i] < source + length) {
				// map it!
				next[i] = current[i] + offset;
			}
		}
	}

	free(current);

	// cool, now find the smallest number in next
	result = UINT_MAX;
	for (int i = 0; i < num_seeds; i++) {
		if (next[i] < result) {
			result = next[i];
		}
	}

	free(next);

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%u\n", result);
}

/*
 *
 #include <stdio.h>
 #include <string.h>
 #include <assert.h>
 #include <ctype.h>
 #include <stdlib.h>
 #include <stdbool.h>
 #include <limits.h>
 
 // `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
 // one +1 to handle the \n
 // another +1 to handle being a NULL terminated string, rather than needing to track size
 #define MAX_LINE_LENGTH 212
 
 // these are the max number of seeds we'll care about... anymore and we'll break
 #define MAX_SEED_RANGES 100
 
 typedef struct {
 	uint start;
 	uint length;
 } seedRange;
 
 // if there is overlap, sets the overlap to `ol` and returns ol, otherwise returns NULL
 *seedRange overlap(seedRange a, seedRange b, *seedRange ol) seedRange {
 
 
 
 }
 
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
 	 \/
 
 	// first, we want to pull the first line containing all the seeds
 	assert(fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);
 	// assert we have a full line and aren't missing any input
 	assert (strchr(buf, '\n') != NULL);
 
 
 	seedRange current[MAX_SEED_RANGES];
 	seedRange next[MAX_SEED_RANGES];
 	uint numSeedRanges = 0;
 
 	// for the first line, get numbers in pairs
 	uint start, length, bytesRead;;
 	for (char *partOfString = buf + strlen("seeds: "); sscanf(partOfString, "%u %u%n", &start, &length, &bytesRead) == 2; partOfString += bytesRead + 1) {
 		next[numSeedRanges].start = start;
 		next[numSeedRanges++].length = length;
 	}
 
 	memcpy(current, next, sizeof(next));
 
 	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
 		assert (strchr(buf, '\n') != NULL);
 
 		if (strlen(buf) == 1) {
 			// empty line, skip
 			continue;
 		} else if (strchr(buf, ':') != NULL) {
 			// mapping header, skip
 			continue;
 		}
 
 		// we have another mapping
 		uint destination, source, length;
 		assert (sscanf(buf, "%u %u %u", &destination, &source, &length) == 3);
 
 		for (int i = 0; i < numSeedRanges; i++) {
 			printf("overlap: %b\n", overlaps(current[i], (seedRange){source, length}));
 		}
 
 		memcpy(current, next, sizeof(next));
 	}
 
 	// assert that we didn't stop procesing early due to an error
 	assert (!ferror(stdin));
 	assert (feof(stdin));
 
 	printf("%u\n", result);
 }
*/
