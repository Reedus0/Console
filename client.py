import socket
import sys
import random

HOST = "deb"
HOST = "127.0.0.1"
PORT = 12225
letters = "qwertyuiopasdfghjklzxcvbnmйцукенгшщзхъфывапролджэячсмитьбюQWERTYUIOPASDFGHJKLZXCVBNMЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ 1234567890!@#$%^&*()"
def test(name):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        s.send(name.encode() + b"\0")
        while True:
            lenght = random.randint(100,250)
            dat = b""
            for _ in range(lenght):
                dat += random.choice(letters).encode()
            s.send(b"INFO "+dat+b"\0")

if __name__ == "__main__":
    test(sys.argv[0])