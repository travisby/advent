#include <stdio.h>
#include <stdbool.h>
#include <assert.h>
#include <stdlib.h>

int main() {
	int result = 0;

	int c;
	while (c != EOF) {
		char number[3] = {'\0', '\0'};
		do {
			c = getchar();
			if (c >= 48 && c < 58) {
				if (number[0] == '\0') {
					number[0] = c;
					// the last number we encountered could have been the first too!
					number[1] = c;
				} else {
					number[1] = c;
				}
			}
		} while (c != '\n' && c != EOF);
		
		// did we fill in _anything_, or do we think we just encountered \n for the last line of the file?
		if (number[0] != '\0' && number[1] != '\0') {
			result += atoi(number);
		}
	}

	printf("%d\n", result);
}
