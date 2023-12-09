#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <stdbool.h>
#include <stdlib.h>
#include <math.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input| sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 122

// `expr $(while read line; do echo $line | tr -cd ' ' | wc -c; done < input | sort -n | tail -n1) +1`
#define MAX_HISTORY 21

// finds the next value in the sequence
int predict(int *values, int numHistory);

// finds the previous value in the sequence
int rpredict(int *values, int numHistory);

#ifdef PART1
#define PREDICT predict
#elifdef PART2
#define PREDICT rpredict
#else
#error "Must define PART1 or PART2"
#endif

int main() {
	/*
	 * This was a fun problem to solve
	 *
	 * We recursively solve for a derivative of the sequence
	 * (we actually use malloc/free in the recursive function, so watch out for memory explosion :scream:)
	 *
	 * And if we're solving for part1 we're looking for the n+1 value
	 * and if we're solving for part2 we're looking for the 0th (well, -1st? value)
	 *
	 * the n+1 value is the nth value + the last derivative
	 * the -1st value is the 0th value - the first derivative
	 */
	int result = 0;

	char buf[MAX_LINE_LENGTH];

	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		int values[MAX_HISTORY];

		int history = 0;
		char *saveptr;
		for (char *token = strtok_r(buf, " ", &saveptr); token != NULL; token = strtok_r(NULL, " ", &saveptr)) {
			values[history++] = atoi(token);
		}

		result += PREDICT(values, history);
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	printf("%d\n", result);
}

int predict(int *values, int numHistory) {
	bool allZero = true;
	for (int i = 0; i < numHistory; i++) {
		if (values[i] != 0) {
			allZero = false;
			break;
		}
	}
	if (allZero) {
		return 0;
	}

	// DANGER
	int *nextHistory = malloc(sizeof(int) * (numHistory - 1));
	for (int i = 0; i < numHistory - 1; i++) {
		nextHistory[i] = values[i + 1] - values[i];
	}

	int result = values[numHistory-1] + predict(nextHistory, numHistory-1);
	free(nextHistory);

	return result;
}

int rpredict(int *values, int numHistory) {
	bool allZero = true;
	for (int i = 0; i < numHistory; i++) {
		if (values[i] != 0) {
			allZero = false;
			break;
		}
	}
	if (allZero) {
		return 0;
	}

	// DANGER
	int *nextHistory = malloc(sizeof(int) * (numHistory - 1));
	for (int i = 0; i < numHistory - 1; i++) {
		nextHistory[i] = values[i + 1] - values[i];
	}

	int result = values[0] - rpredict(nextHistory, numHistory-1);

	free(nextHistory);

	return result;
}
