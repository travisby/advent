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
#define MAX_LINE_LENGTH 12

#define HAND_SIZE 5

#define NUM_HANDS 20

typedef int hand[HAND_SIZE];
// strhand converts a string representation of a hand into a hand type
// on an invalid hand it will return NULL
// otherwise it returns h
int *strhand(char *str, hand h);

// handstr convers a hand to a str
// it either returns str or NULL if the buffer is too small 
char *handstr(hand h, char *str, size_t len);

enum handType {
	HIGH_CARD,
	ONE_PAIR,
	TWO_PAIR,
	THREE_OF_A_KIND,
	FULL_HOUSE,
	FOUR_OF_A_KIND,
	FIVE_OF_A_KIND,
};

// handhandType returns the type of hand
enum handType handhandType(hand h, hand scratch);
// it takes in a scratch hand to avoid modifying the original, we do not assume that scratch already has a dupe of hand's data

// handint converts a hand to an integer representation
// this is used to sort hands
int handint(hand h, hand scratch);

int main() {
	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/*
	 * This is our first problem where I think we need an O(n) memory solution
	 * based on the number of hands
	 * (technically, since we use staticly defined arrays it's all O(1), but if we wanted
	 *  to handle more hands we'd need to increase our array size)
	 */

	hand hands[NUM_HANDS];
	int numHands = 0;
	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		strhand(buf, hands[numHands++]);
	}

	for (int i = 0; i < numHands; i++) {
		handstr(hands[i], buf, MAX_LINE_LENGTH);
		hand scratch;
		printf("%s = %d\n", buf, handint(hands[i], scratch));
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	int result = 0;
	printf("%d\n", result);
}

int *strhand(char *str, hand h) {
	if (strlen(str) <= HAND_SIZE) {
		return NULL;
	}

	for (int i = 0; i < HAND_SIZE; i++) {
		switch(str[i]) {
			case 'A':
				h[i] = 14;
				break;
			case 'K':
				h[i] = 13;
				break;
			case 'Q':
				h[i] = 12;
				break;
			case 'J':
				h[i] = 11;
				break;
			case 'T':
				h[i] = 10;
				break;
			default:
				if (!isdigit(str[i])) {
					return NULL;
				}
				h[i] = str[i] - '0';
				break;
		}
	}

	return h;
}

char *handstr(hand h, char *str, size_t len) {
	if (len < HAND_SIZE + 1) {
		return NULL;
	}

	for (int i = 0; i < HAND_SIZE; i++) {
		switch(h[i]) {
			case 14:
				str[i] = 'A';
				break;
			case 13:
				str[i] = 'K';
				break;
			case 12:
				str[i] = 'Q';
				break;
			case 11:
				str[i] = 'J';
				break;
			case 10:
				str[i] = 'T';
				break;
			default:
				str[i] = h[i] + '0';
				break;
		}
	}

	str[HAND_SIZE] = '\0';

	return str;
}

int compare_ints(const void *p, const void *q) {
	// thanks wikipedia
	int x = *(const int *)p;
	int y = *(const int *)q;

	return (x > y) - (x < y);
}

// we take in a scratch hand because we don't want to modify the original
// but we want to be able to sort it, so we need a copy
// we will copy the hand into scratch, then sort scratch
enum handType handhandType(hand h, hand scratch) {
	memcpy(scratch, h, sizeof(hand));

	// it's easy to count pairs if we sort the hand, pairs will appear at n and n-1
	qsort(scratch, HAND_SIZE, sizeof(int), compare_ints);

	// it's impossible to have 3 kinds of pairs with 5 cards
	// at most we can have 2
	// we care how many of those cards there are in a pair
	int numberOfCardsInAPair[2] = {0, 0};;
	int numPairs = 0;

	// XXX: numberOfCardsInAPair doesn't include the original card
	// math is hard

	for (int i = 1; i < HAND_SIZE; i++) {
		// if n and n-1 match, we have a pair and we increment the count of that pair
		if (scratch[i] == scratch[i - 1]) {
			numberOfCardsInAPair[numPairs] += 1;
		// if we HAD a pair before but no longer have a pair, bump numPairs to get ready for the next pair
		} else if (numberOfCardsInAPair[numPairs] > 0){
			numPairs += 1;
		}
	}

	// was our last number a pair, and we need to bump numPairs?
	if (scratch[HAND_SIZE] == scratch[HAND_SIZE - 1] && numberOfCardsInAPair[numPairs] > 0) {
		numPairs += 1;
	}

	// XXX: numberOfCardsInAPair doesn't include the original card
	// math is hard

	if (numberOfCardsInAPair[0] == 4) {
		return FIVE_OF_A_KIND;
	} else if (numberOfCardsInAPair[0] == 3) {
		return FOUR_OF_A_KIND;
	} else if ((numberOfCardsInAPair[0] == 2 && numberOfCardsInAPair[1] == 1) || (numberOfCardsInAPair[0] == 1 && numberOfCardsInAPair[1] == 2)) {
		return FULL_HOUSE;
	} else if (numberOfCardsInAPair[0] == 2) {
		return THREE_OF_A_KIND;
	} else if (numberOfCardsInAPair[0] == 1 && numberOfCardsInAPair[1] == 1) {
		return TWO_PAIR;
	} else if (numberOfCardsInAPair[0] == 1) {
		return ONE_PAIR;
	}
	return HIGH_CARD;
}

// our calculation here is fun:
// the 10^6 place is the hand type (four of a kind, full house, etc)
// the 10^5 place is the first card
// the 10^4 place is the second card
// etc
// which should lead to correct ordering
int handint(hand h, hand scratch) {
	int result = handhandType(h, scratch) * pow(10, 6);
	// thinking through the loop was hard because the index is 0 but the place is 5
	// we can just write it all out...
	result += h[0] * pow(10, 5);
	result += h[1] * pow(10, 4);
	result += h[2] * pow(10, 3);
	result += h[3] * pow(10, 2);
	result += h[4] * pow(10, 1);

	return result;
}
