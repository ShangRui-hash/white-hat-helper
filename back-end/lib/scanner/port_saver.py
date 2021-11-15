import re
import threading 

from lib.redis.port import HostPortHash 

class PortSaver(threading.Thread): 
    """端口与服务版本存储器""" 
    def __init__(self,company_id,port_queue): 
        threading.Thread.__init__(self)
        self.port_queue = port_queue 
        self.host_port_hash = HostPortHash(company_id) 
    def run(self): 
        while True: 
            port_sevice_info = self.port_queue.get() 
            if port_sevice_info == "done": 
                break 
            self.do_work(port_sevice_info)
    def check_ip(self,ipAddr):
        compile_ip=re.compile('^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$')
        if compile_ip.match(ipAddr):
            return True    
        else:    
            return False
    def do_work(self,port_service_info): 
        for (key,value) in port_service_info.items(): 
            if not self.check_ip(key): 
                continue
            self.host_port_hash.add(key,value)