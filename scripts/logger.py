import socket

def log(text:str):
    text = str(text)
    soc = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    soc.connect(("127.0.0.1", 12225))
    soc.send("server_log".encode() + b"\0")
    soc.send(b"INFO "+text.encode() + b"\0")
    soc.close()
    
if __name__ == "__main__":
    import time
    for i in range(15):
        log(f"hello world {i}")
        time.sleep(0.1)