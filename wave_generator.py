import random

total = 2176

n=300
x= 1.02
random.seed(1)
threshold= 5
list =[]
counter = 0

while counter <= total:


  for i in range(0,n):
    y=pow(x,i)
    if ( random.randint(0,1)):
      val = (y*100/(pow(x,n-threshold)))
    else :
      val = (y*100/(pow(x,n+threshold)))
      

    counter = counter+1
    list.append(val)


  for i in range(n,0,-1):
    y=pow(x,i)
    if ( random.randint(0,1)):
      val = (y*100/(pow(x,n-threshold)))
    else :
      val = (y*100/(pow(x,n+threshold)))


    list.append(val)
    counter = counter+1

  print(list[:2048])
  print (list[2048:2176])
