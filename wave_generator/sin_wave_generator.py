import random
import math

import matplotlib.pyplot as plt

total = 2176

n=300
x= 1.02
random.seed(1)
threshold= 5
list =[]
counter = 0




while counter <= total:
    a = counter / 272
    if a % 2 == 0:
        a = a + 1
    y = ((272**2) - (counter-(272 *  a))**2)**0.5

    if ( random.randint(0,1)):
      val = (y*100/(pow(x,n-threshold)))
    else :
      val = (y*100/(pow(x,n+threshold)))

    list.append(val)
    counter = counter+1

print(list[:2048])
print (list[2048:2176])

plt.plot(list[:2048])
plt.ylabel('some numbers')
plt.show()
