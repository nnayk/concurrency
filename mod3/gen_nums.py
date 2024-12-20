file="foo.txt"

with open(file,"w+") as f:
    data = ""
    for i in range(1,10001):
        data += str(i)+" "
    f.write(data)
