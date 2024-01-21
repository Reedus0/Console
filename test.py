import os
import random
import time

def main():
    for i in range (100):
        time.sleep(1)
        for i in range(1, 10, 1):
            os.system(f"python3 client.py botnumber" + str(random.randint(1, 300)))

if __name__ == "__main__":
    main()