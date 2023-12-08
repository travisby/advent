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

void printNetwork(Network network) {
	for (int i = 0; i < network.nodeCount; i++) {
		Node node = network.nodes[i];
		printf("%s: %s %s\n", node.name, node.left->name, node.right->name);
	}
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

	printNetwork(network);
	
	Node *current = findNode(&network, "AAA");
	assert (current != NULL);

	int steps = 0;
	int instructionsLength = strlen(instructions);
	while (strcmp(current->name, "ZZZ") != 0) {
		if (instructions[steps % instructionsLength] == 'L') {
			current = current->left;
		} else if (instructions[steps % instructionsLength] == 'R') {
			current = current->right;
		}

		steps++;
		assert (current != NULL);
	}

	// assert that we didn't stop procesing early due to an error
	assert (!ferror(stdin));
	assert (feof(stdin));

	int result = 0;
	printf("%d\n", steps);
}
