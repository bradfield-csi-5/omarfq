
#include "calc.h"
#include <stdio.h>
#define MAXVAL 100 /* maximum depth of val stack */

int sp = 0;         /* next free stack position */
double val[MAXVAL]; /* value stack */

/* push:  push f onto value stack */
void push(double f) {
  printf("Before push, sp is: %d\n", sp);
  if (sp < MAXVAL) {
    val[sp++] = f;
    printf("%d\n", sp);
    for (int i = 0; i < MAXVAL; i++) {
      printf("%f, ", val[i]);
    }
  } else
    printf("error: stack full, can't push %g\n", f);
  printf("After push, sp is: %d\n", sp);
}
/* pop:  pop and return top value from stack */
double pop(void) {
  if (sp > 0)
    return val[--sp];
  else {
    printf("error: stack empty\n");
    return 0.0;
  }
}
