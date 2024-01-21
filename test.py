import os
import random

def main():
    for i in range(1, 10, 1):
        os.system(f"python3 client.py botnumber" + str(random.randint(1, 300)))
    return 0

if __name__ == "__main__":
    main()