#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <ctype.h>
#include <stdlib.h>
#include <math.h>
#include <stdbool.h>

// `expr $(while read line; do echo -n $line | wc -c; done < input_part2 | sort -n | tail -n1) + 2`
// one +1 to handle the \n
// another +1 to handle being a NULL terminated string, rather than needing to track size
#define MAX_LINE_LENGTH 38

#ifdef PART1
#define MAX_RACES 10
#define RACE_TYPE int
#define RACE_FLOAT double
#else
#define MAX_RACES 1
#define RACE_TYPE long long
#define RACE_FLOAT long double
#endif

typedef struct {
	RACE_TYPE time;
	RACE_TYPE distance;
} race;


int main() {
	#if !defined(PART1) && !defined(PART2)
		#error "You must define either PART1 or PART2"
	#endif

	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	 * Here we're going to read all of the input and actually store it for later consumption
	 * because each "round" is split between L1 and L2
	 * we'll keep track of every race's time and distance
	 *
	 * then we'll use some quadratic formulas to determine at which point we cross into "record breaking territory"
	 * and then we'll count the number of integers between those two points to multiply against our result to get our final answer
	 */
	race races[MAX_RACES];
	int numRaces = 0;

	// the number of lines is well-defined: the first line is the time of races
	// the second line is the distance to beat
	// although the number of races is not well defined
	
	char *saveptr;

	
	// first, get all of the times
	assert (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);
	// assert we have a full line
	assert (strchr(buf, '\n') != NULL);
	#ifdef PART1
	for (char *token = strtok_r(buf, " ", &saveptr); token != NULL; token = strtok_r(NULL, " ", &saveptr)) {
		// throw out leads like "Time:" or "Distance:"
		// if it's not a number, we don't want it!
		if (!isdigit(token[0])) {
			continue;
		}

		races[numRaces++].time = atoi(token);
	}
        #elifdef PART2
	char temp[MAX_LINE_LENGTH];
	int tempCount  = 0;
	for (int i = 0; buf[i] != '\0'; i++) {
		if (isdigit(buf[i])) {
			temp[tempCount++] = buf[i];
		}

	}
	// finish making it a string
	temp[tempCount++] = '\0';
	races[numRaces++].time = atoll(temp);
	#endif


	// now, get all of the distances
	assert (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);
	// assert we have a full line
	assert (strchr(buf, '\n') != NULL);
	numRaces = 0; // we need to restart numRaces; this assumes the number of times and distances were the same
	#ifdef PART1
	for (char *token = strtok_r(buf, " ", &saveptr); token != NULL; token = strtok_r(NULL, " ", &saveptr)) {
		// throw out leads like "Time:" or "Distance:"
		// if it's not a number, we don't want it!
		if (!isdigit(token[0])) {
			continue;
		}

		races[numRaces++].distance = atoi(token);
	}
        #elifdef PART2
	tempCount  = 0;
	for (int i = 0; buf[i] != '\0'; i++) {
		if (isdigit(buf[i])) {
			temp[tempCount++] = buf[i];
		}

	}
	// finish making it a string
	temp[tempCount++] = '\0';
	races[numRaces++].distance = atoll(temp);
	#endif


	// ensure there's no more input
	assert (fgets(buf, MAX_LINE_LENGTH, stdin) == NULL);
	assert (!ferror(stdin));
	assert (feof(stdin));

	// our result ends up being a multiplication rather than an addition
	// so we'll start at 1
	RACE_TYPE result = 1;
	for (int i = 0; i < numRaces; i++) {
		// these races are a parabola of the form:
		// distance[total race time, hold time] = -hold^2 + x*hold
		// if we want to find all the times > record holder, we need to solve for:
		// distance = -hold^2 + x*hold ---> 0 = -hold^2 + x*hold - record distance
		// and pick the (exclusive?) range of all integers inside there
		// to find the roots of a parabola, we can use x = (-b +/- sqrt(b^2 - 4ac)) / 2a
		// where a = -1, b = total time, c = -record distance
		// (basically, we convert our refcord ddistance = ax^2 + bx + c to the form ax^2 + bx + c = 0, so we start solving the standard "x intercepts of a parabola")

		int a = -1;
		RACE_TYPE b = races[i].time;
		RACE_TYPE c = -1 * races[i].distance;

		RACE_FLOAT x0 = (-b + sqrt(b*b - 4*a*c)) / (2*a);
		RACE_FLOAT x1 = (-b - sqrt(b*b - 4*a*c)) / (2*a);

		// calculating the integers between the two intercepts was tricky
		// and resulted in "try floor/ceil +/-1 until the numbers look right"
		// I can't tell ya why this is right, only that it is
		result *= ceil(x1) - floor(x0) -1;
	}

	#ifdef PART1
	printf("%d\n", result);
	#elifdef PART2
	printf("%lld\n", result);
	#endif
}
