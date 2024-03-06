import os
import re
import random
bots = os.listdir("/data")
with open("/app/prox.txt") as f:
    prox_data = f.read()

print("Updating proxies...")
for bot in bots:
    print(bot)
    counter = 0
    if os.name == "nt":
        with open(f"/data/{bot}/proxy_list.txt", "w") as f:
            f.write(prox_data)
        with open(f"/data/{bot}/proxy_list.txt", "r", encoding="utf-8") as f:
            proxy_list = f.read().split('\n')

        files = os.listdir(f"/data/{bot}/")
        for file in files:
            if file.endswith(".bat") and "start_" in file:
                with open(f"/data/{bot}/"+file, 'r', encoding="utf-8") as f:
                    data = f.read()
                proxies = re.findall(r"\"(?:.+)\" \"(?:.+)\" \"(?:.+)\" \"(.*)\" \"(.*)\"\n", data)[0]
                for p in proxies:
                    proxy = random.choice(proxy_list)
                    data = data.replace(p, proxy)
                with open(f"/data/{bot}/"+file, 'w', encoding="utf-8") as f:
                    f.write(data)
                counter += 1
        print(counter)

print("Done")