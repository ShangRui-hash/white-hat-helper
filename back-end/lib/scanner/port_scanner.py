import threading 
import nmap3 

class PortScanner(threading.Thread):
    """端口扫描与服务版本识别""" 
    def __init__(self, ip_queue, port_queue):
        threading.Thread.__init__(self)
        self.ip_queue = ip_queue
        self.port_queue = port_queue
        self.nmap = nmap3.Nmap()
    def run(self): 
        while True:
            ip = self.ip_queue.get()
            if ip == "done": 
                break
            self.do_work(ip)
        self.port_queue.put("done")
           
    def do_work(self, ip):
        print("scan port:", ip)
        results = self.nmap.nmap_version_detection(ip)
        self.port_queue.put(results)