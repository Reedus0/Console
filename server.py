import socket
import json
import os
from threading import Thread

host = "127.0.0.1"

logs = []


def main():
    tcp = Thread(target=tcp_server)
    http = Thread(target=http_server)
    tcp.start()
    http.start()

def tcp_server():
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server.bind((host, 8888))
    server.listen()
    try:
        while True:
            conn, addr = server.accept()
            client = Thread(target=handle_client, args=[conn])
            client.run()
            conn.close()
    except (KeyboardInterrupt):
        socket.close()


def handle_client(conn):
    data = conn.recv(4096)
    decoded_data = data.decode()
    logs.append(decoded_data.split("\x00"))
    print("Recieved logs:", decoded_data.split("\x00"))


def http_server():
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server.bind((host, 9999))
    server.listen()
    try:
        while True:
            conn, addr = server.accept()
            request = conn.recv(1024)
            request_method = request.decode().split(" ")[0]
            if (request_method == "GET"):
                response(conn)
            else:
                reset(conn, request.decode())
            conn.close()
    except (KeyboardInterrupt):
        socket.close()


def response(conn):
    global logs
    data = json.dumps(logs)
    packet = f"""HTTP/1.1 200 OK
Content-Length: {len(data)}
Content-Type: application/json
Access-Control-Allow-Origin: http://127.0.0.1:8001

{data}
"""
    print("Sent logs:", logs)
    conn.send(packet.encode())

    logs = []


def reset(conn, data):
    bot = data.split("\n")[-1]
    os.system(f"sudo systemctl restart {bot}.service")
    packet = """HTTP/1.1 200 OK
Content-Length: 
Content-Type: application/json
Access-Control-Allow-Origin: http://127.0.0.1:8001

{"status": "ok"}
"""
    conn.send(packet.encode())


if __name__ == "__main__":
    main()
