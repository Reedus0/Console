import base64
import sys
import requests
import random
import re
import hashlib
import os
import logger
prox_sess = requests.Session()
PROXYES = []
VER = "3.6.125"
REGION = 2
OS = "android"
SUB_DISTRIBUTOR = 0
APPID = "globle"
COUNTRY = "RU"
UID = 0
APP_TYPE = 1
USER_AGENT = "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 4 XL Build/TP1A.221005.002.B2)"
DEVICE_TOKEN = "firebase_dDVeHsMJQkqPNn{h}SNJ_BBH:APA9{k}bGnFaaS1Ck8ix0GJNPOndrWfjEFdahmlBv2Fg69VzFpWCFILyWL_tDAiqN46dm8hgDCjJQrRMNkSoO5KEjsCoUnCZWCOhngNrFNJn6OwJXSW4G794{j}P5o_xHG9fj3zqinXgbz{i}_"


LANG = "ru"
OPERATING_COMPANY = "MTS RUS" + str(random.randint(0, 9))
PLATFORM_TYPE = ("2",)


def md5(data):
    if type(data) == str:
        data = data.encode()
        
    return hashlib.md5(data).hexdigest()


def callRegisterPhp(username, password, imei=None, prox=""):
    params = {
        "username": username,
        "password": md5(md5(password)),
        "distributor": "1",
        "sub_distributor": SUB_DISTRIBUTOR,
        "country": COUNTRY,
        "appid": APPID,
        "os": OS,
        "imei": imei,
        "clientvar": VER,
        "android_id": imei,
        "region": REGION,
        "app_type": APP_TYPE,
    }
    headers = {
        "User-Agent": "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 4 XL Build/TP1A.221005.002.B2)",
        "Host": "www.pppoker.club",
    }

    response = prox_sess.get(
        "https://www.pppoker.club/poker/api/register.php", params=params, headers=headers, proxies={'http': "socks5://"+prox, 'https': "socks5://"+prox}
    )
    try:
        return response.json()
    except:
        logger.log(response.text)


def register(a, prox=""):
    imei = "020000%02x%02x%02x" % (
        random.randint(0, 255),
        random.randint(0, 255),
        random.randint(0, 255),
    )
    DEVICE_TOKEN.replace("{i}", str(random.randint(0, 9)))
    DEVICE_TOKEN.replace("{j}", str(random.randint(0, 9)))
    DEVICE_TOKEN.replace("{k}", str(random.randint(0, 9)))
    DEVICE_TOKEN.replace("{h}", str(random.randint(0, 9)))
    try:
        res = callRegisterPhp(a, a, imei, prox)
        logger.log(res)
        if res["code"] == -2:
            logger.log("Достигнут максимум регистраций")
        return res["code"]
    except requests.exceptions.ConnectionError as e:
        return -200


def gen_login():
    letters = "1234567890qwertyuiopasdfghjklzxcvbmn"
    lens = random.randint(6, 10)
    login = ""
    for _ in range(lens):
        login += random.choice(letters)
    return login

def add(login, password="", tables="", proxy1="", proxy2="", folder="bot1"):
    with open("/data/bot1/proxy_list.txt", "r") as f:
        proxy_list = f.read().split("\n")

    if not tables:
        tables = "50K 500K 1M 5M 20M"

    if not password:
        password = login

    if not proxy1:
        proxy1 = random.choice(proxy_list)
        proxy_list.remove(proxy1)
        logger.log(f"Выбран случайный прокси: {proxy1}")

    if not proxy2:
        proxy2 = random.choice(proxy_list)
        logger.log(f"Выбран случайный резервный прокси: {proxy2}")

    # chack if os is linux
    if os.name != "nt":
        print("{"+login+"}")

        with open(f"/data/{folder}/start_{login}.sh", "w") as f:
            f.write(f"""python main.py "{login}" "{password}" "{tables}" "{proxy1}" "{proxy2}"
""")
    else:
        file_content = f"""
python main.py "{login}" "{password}" "{tables}" "{proxy1}" "{proxy2}"
@echo off
set /p "id=Press Enter to close"
"""
        with open(f"start_{login}.bat", "w") as f:
            f.write(file_content)
def get_proxyes(index=1):
    result = []
    res = prox_sess.get(
        # "https://freeproxylist.ru/protocol/socks?page=" + str(index),
        "https://advanced.name/ru/freeproxy?type=socks5"
    )
    # hosts = re.findall(r"class=\"w-30 tblport\">([^<]+)", res.text)
    # ports = re.findall(r"class=\"w-10 tblport\">([^<]+)", res.text)
    hosts_bs64 = re.findall(r"data-ip=\"([^:\"]+)", res.text)
    ports_bs64 = re.findall(r"data-port=\"([^:\"]+)", res.text)
    hosts = [base64.b64decode(host_bs64).decode() for host_bs64 in hosts_bs64]
    ports = [base64.b64decode(port_bs64).decode() for port_bs64 in ports_bs64]
    for i in range(len(hosts)):
        result.append(hosts[i]+":"+ports[i])
    logger.log("aviable proxyes: " + str(len(result)))
    return result

def __get_proxyes():
    with open("/data/bot1/proxy_list.txt", "r") as f:
        proxy_list = f.read().split("\n")
    res = []
    for prox in proxy_list:
        prox = prox.split(":")
        res.append(f"{prox[2]}:{prox[3]}@{prox[0]}:{prox[1]}")
    return res

def _get_proxyes():
    proxs = prox_sess.get("https://advanced.name/freeproxy/65c1c9eb9ae92").text.split()
    return proxs    

def auto_reg(prefix, tables, folder):
    try:
        p = PROXYES.pop(random.randint(0, len(PROXYES) - 1))
        try:
            prox_sess.get("https://2ip.ru", timeout=10 , proxies={"https": f"socks5://{p}"})
        except (requests.exceptions.ConnectionError, requests.exceptions.ConnectTimeout) as e:
            logger.log("Proxy {} is bad".format(p))
            return auto_reg(prefix,tables, folder)
        login = prefix+gen_login()
        reg_res = register(login, p)
        if reg_res == 0:
            add(login, tables=tables, folder=folder)
            return 0
        else:
            return auto_reg(prefix,tables, folder)
    except Exception as e:
        return e
    

if __name__ == "__main__":
    print(sys.argv)
    PROXYES = __get_proxyes()
    if len(sys.argv) == 1:
        logger.log(auto_reg("test", None, "bot4"))
    elif len(sys.argv) >= 4:
        if sys.argv == 4:
            sys.argv.append("")
        prefix, count, folder, tables = sys.argv[1:]
        count = int(count)
        i = 0
        while i < count:
            res = auto_reg(prefix, tables, folder)
            if res == 0:
                i+=1
        sys.path.insert(0, f"/data/{folder}")
        os.chdir(f"/data/{folder}")
        import startall
                
        
    
