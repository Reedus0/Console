import socket
import sys

HOST = "0.0.0.0"  
PORT = 8888  

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    name = sys.argv[1]
    s.connect((HOST, PORT))
    s.send(name.encode() + b'\0')
    s.send(b"INFO <script>alert(1)</script>\0")
    s.send(b"""INFO Lorem ipsum dolor sit amet consectetur adipisicing elit. Deleniti numquam sint a incidunt veniam quas optio, blanditiis vero officia quaerat rem nobis recusandae ipsam deserunt fuga nisi commodi totam vel.\0""")
    s.send(b"""ERROR Lorem ipsum dolor sit amet consectetur adipisicing elit. Deleniti numquam sint a incidunt veniam quas optio, blanditiis vero officia quaerat rem nobis recusandae ipsam deserunt fuga nisi commodi totam vel.\0""")