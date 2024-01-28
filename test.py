import os
import random
from threading import Thread
import client


def main():
    workers = 100
    for _ in range(workers):
        # trunk-ignore(bandit/B311)
        Thread(target=client.test, args=("botnumber" + str(random.randint(1, 300)),)).start()

if __name__ == "__main__":
    main()