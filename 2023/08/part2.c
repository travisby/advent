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
#define MAX_LINE_LENGTH 271

#define MAX_NODES 1000

typedef struct Node {
	char name[4];
	struct Node *left;
	struct Node *right;
} Node;

typedef struct Network {
	int nodeCount;
	Node nodes[MAX_NODES];
} Network;

typedef struct visit {
	Node *node;
	int steps;
} visit;

typedef struct cycle {
	int start;
	int end;
	int z1;
	int z2;
} cycle;

bool cycleAtZ(cycle c, int64_t t) {
	return t >= c.start && (((t - c.start) % (c.end - c.start) == c.z1 - c.start) || ((t - c.start) % (c.end - c.start) == c.z2 - c.start));
}

Node *findNode(Network *network, char *name) {
	for (int i = 0; i < network->nodeCount; i++) {
		if (strcmp(network->nodes[i].name, name) == 0) {
			return &network->nodes[i];
		}
	}
	return NULL;
}


void addNode(Network *network, char *name, char *left, char *right) {
	Node *node = findNode(network, name);
	if (node == NULL) {
		node = &network->nodes[network->nodeCount++];
		strcpy(node->name, name);
	}

	Node *l = findNode(network, left);
	if (l == NULL) {
		l = &network->nodes[network->nodeCount++];
		strcpy(l->name, left);
	}
	node->left = l;

	Node *r = findNode(network, right);
	if (r == NULL) {
		r = &network->nodes[network->nodeCount++];
		strcpy(r->name, right);
	}
	node->right = r;
}

int main() {
	/*
	 * Today's input creates a map
	 * of [LOC] -> [nextLocationLeft, nextLocationRight]
	 * and a series of Left/Rights
	 *
	 * to accomplish this, we will need to store the whole input
	 * rather than processling line-by-line as we normally try to
	 */
	// buf holds (what should be) a whole line from stdin
	char buf[MAX_LINE_LENGTH];

	/* 
	 * first get the LR set of instructions
	 * and remove the newline
	 * so we can iterate over it easily
	 */
	assert (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);
	// assert we have a full line / aren't missing input
	char *newLine = strchr(buf, '\n');
	assert (newLine != NULL);

	char instructions[MAX_LINE_LENGTH];
	assert (strncpy(instructions, buf, newLine - buf + 1) != NULL);
	instructions[newLine - buf] = '\0';

	// next line should be just a newline... skip it
	assert (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL);
	assert (strcmp(buf, "\n") == 0);

	Network network = {0};
	while (fgets(buf, MAX_LINE_LENGTH, stdin) != NULL) {
		// assert we have a full line
		assert (strchr(buf, '\n') != NULL);

		char n[4], l[4], r[4];

		assert (sscanf(buf, "%3s = (%3s, %3s)", n, l, r) == 3);
		addNode(&network, n, l, r);
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	Node startNodes[MAX_NODES];
	int nodeCount = 0;
	for (int i = 0; i < network.nodeCount; i++) {
		Node node = network.nodes[i];
		if (node.name[2] == 'A') {
			startNodes[nodeCount++] = node;
		}
	}

	cycle cycles[100];
	int numCycles = 0;
	int instructionCount = strlen(instructions);
	for (int i = 0; i < nodeCount; i++) {
		int steps = 0;
		Node *node = &startNodes[i];
		visit visits[100000];
		int numVisits = 0;

		// XXX; I don't know how we could prove that there's not many zs
		// in each cycle, but it seems like there's only 1 or 2
		int zsAt[2] = {-1, -1};
		int numZs = 0;

		while (true) {
			if (node->name[2] == 'Z') {
				zsAt[numZs++] = steps;
			}

			bool cycl = false;
			for (int j = 0; j < numVisits; j++) {
				if (visits[j].node == node && visits[j].steps % instructionCount == steps % instructionCount) {
					printf("Beginning of cycle is at: %d, end is at: %d, the z is at %d and %d\n", visits[j].steps, steps, zsAt[0], zsAt[1]);
					cycl = true;

					cycles[numCycles++] = (cycle){visits[j].steps, steps, zsAt[0], zsAt[1]};
					break;
				}
			}


			if (cycl) {
				break;
			}
			visits[numVisits++] = (visit){node, steps};

			if (instructions[steps % instructionCount] == 'L') {
				node = node->left;
			} else {
				node = node->right;
			}

			steps++;
		}
	}

	int bigCycle = 3;
	int64_t i = cycles[bigCycle].z1;
	while (true) {
		bool found = true;
		for (int j = 0; j < numCycles; j++) {
			if (!cycleAtZ(cycles[j], i)) {
				found = false;
				break;
			}
		}

		if (found) {
			printf("%lld\n", i);
			return 0;
		}
		i+= cycles[bigCycle].end - cycles[bigCycle].start;
	}

	printf("%d\n", 0);
}
