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

#define NUM_HANDS 1000

typedef struct hand {
	// XXX: we actually don't need to store cards
	// but the refactor is laborious
	int cards[HAND_SIZE];
	int bid;

	// keeping track of score is an optimization
	// we could keep re-re-re-calculating it every time we need it
	int64_t score;
} hand;
// strhand converts a string representation of a hand into a hand type
// on an invalid hand it will return NULL
// otherwise it returns h
hand *strhand(char *str, hand *h, hand *scratch);

// handstr convers a hand to a str
// it either returns str or NULL if the buffer is too small 
#ifdef DEBUG
char *handstr(hand h, char *str, size_t len);
#endif

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
enum handType handhandType(hand h, hand *scratch);
// it takes in a scratch hand to avoid modifying the original, we do not assume that scratch already has a dupe of hand's data

// handint converts a hand to an integer representation
// this is used to sort hands
int64_t handint(hand h, hand *scratch);

int compare_hands(const void *p, const void *q);

int main() {
	#if !defined(PART1) && !defined(PART2)
		#error "You must define either PART1 or PART2"
	#endif

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

		hand scratch;
		strhand(buf, &hands[numHands++], &scratch);
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	int64_t result = 0;
	qsort(hands, numHands, sizeof(hand), compare_hands);
	for (int i = 0; i < numHands; i++) {
		#ifdef DEBUG
		printf("%s %d * %d\n", handstr(hands[i], buf, MAX_LINE_LENGTH), hands[i].bid, i+1);
		#endif
		result += hands[i].bid * (i+1);
	}

	printf("%d\n", result);
}

hand *strhand(char *str, hand *h, hand *scratch) {
	if (strlen(str) <= HAND_SIZE) {
		return NULL;
	}

	for (int i = 0; i < HAND_SIZE; i++) {
		switch(str[i]) {
			case 'A':
				h->cards[i] = 14;
				break;
			case 'K':
				h->cards[i] = 13;
				break;
			case 'Q':
				h->cards[i] = 12;
				break;
			case 'J':
				#ifdef PART1
				h->cards[i] = 11;
				#elifdef PART2
				// now it's a joker,
				// individually weak
				// but in a group powerful
				h->cards[i] = 1;
				#endif
				break;
			case 'T':
				h->cards[i] = 10;
				break;
			default:
				if (!isdigit(str[i])) {
					return NULL;
				}
				h->cards[i] = str[i] - '0';
				break;
		}
	}

	h->bid = atoi(str + HAND_SIZE + 1);
	h->score = handint(*h, scratch);

	return h;
}

#ifdef DEBUG
char *handstr(hand h, char *str, size_t len) {
	if (len < HAND_SIZE + 1) {
		return NULL;
	}

	for (int i = 0; i < HAND_SIZE; i++) {
		switch(h.cards[i]) {
			case 14:
				str[i] = 'A';
				break;
			case 13:
				str[i] = 'K';
				break;
			case 12:
				str[i] = 'Q';
				break;
			case 1:
			case 11:
				str[i] = 'J';
				break;
			case 10:
				str[i] = 'T';
				break;
			default:
				str[i] = h.cards[i] + '0';
				break;
		}
	}

	str[HAND_SIZE] = '\0';

	return str;
}
#endif

int compare_ints(const void *p, const void *q) {
	// thanks wikipedia
	int x = *(const int *)p;
	int y = *(const int *)q;

	return (x > y) - (x < y);
}

int compare_hands(const void *p, const void *q) {
	int64_t x = ((hand*)p)->score;
	int64_t y = ((hand*)q)->score;


	return (x > y) - (x < y);
}

// we take in a scratch hand because we don't want to modify the original
// but we want to be able to sort it, so we need a copy
// we will copy the hand into scratch, then sort scratch
enum handType handhandType(hand h, hand *scratch) {
	memcpy(scratch->cards, h.cards, sizeof(h.cards));

	// it's easy to count pairs if we sort the hand, pairs will appear at n and n-1
	qsort(scratch, HAND_SIZE, sizeof(int), compare_ints);

	// it's impossible to have 3 kinds of pairs with 5 cards
	// at most we can have 2
	// we care how many of those cards there are in a pair
	int numberOfCardsInAPair[2] = {0, 0};;
	int numPairs = 0;

	// XXX: numberOfCardsInAPair doesn't include the original card
	// math is hard

	int numJokers = 0;

	for (int i = 1; i < HAND_SIZE; i++) {
		if (scratch->cards[i-1] == 1) {
			numJokers += 1;
			continue;
		}

		// if n and n-1 match, we have a pair and we increment the count of that pair
		if (scratch->cards[i] == scratch->cards[i - 1]) {
			numberOfCardsInAPair[numPairs] += 1;
		// if we HAD a pair before but no longer have a pair, bump numPairs to get ready for the next pair
		} else if (numberOfCardsInAPair[numPairs] > 0){
			numPairs += 1;
		}
	}

	// was our last number a pair, and we need to bump numPairs?
	if (scratch->cards[HAND_SIZE] == scratch->cards[HAND_SIZE - 1] && numberOfCardsInAPair[numPairs] > 0) {
		numPairs += 1;
	}
	if (scratch->cards[HAND_SIZE] == 1) {
		numJokers += 1;
	}

	if (numberOfCardsInAPair[0] >= numberOfCardsInAPair[1]) {
		numberOfCardsInAPair[0] += numJokers;
	} else {
		numberOfCardsInAPair[1] += numJokers;
	}

	if (numJokers > 0) {
		// we have a joker, so we have at least a pair
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
// the 10^(2*6) place is the hand type (four of a kind, full house, etc)
// the 10^5 place is the first card
// the 10^4 place is the second card
// etc
// which should lead to correct ordering
//
/*
 *
import math

x = [3, 2, 10, 3, 13]
y = 1

result = 0 
result += y * math.pow(10, 10)

for i in range(len(x)):
    result += x[i] * math.pow(10, (len(x)-i-1)*2)
print(result)
*/
int64_t handint(hand h, hand *scratch) {
	int64_t result = handhandType(h, scratch) * pow(10, 10);
	for (int i = 0; i < HAND_SIZE; i++) {
		result += h.cards[i] * pow(10, (HAND_SIZE - i - 1) * 2);
	}
	return result;
}
