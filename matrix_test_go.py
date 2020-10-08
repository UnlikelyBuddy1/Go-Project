import numpy as np
from numpy import savetxt
size=1024
matrixA=np.random.rand(size,size)
matrixB=np.random.rand(size,size)
matrixes=[matrixA, matrixB]
names=["matrixA.txt", "matrixB.txt"]
index=0
for matrix in matrixes:
    with open(names[index], "w") as a_file:
            np.savetxt(a_file, matrix)
    index+=1
