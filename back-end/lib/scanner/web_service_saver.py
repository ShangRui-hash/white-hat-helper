import threading 
from lib.redis.web_service import HostWebServiceHash 


class WebServiceSaver(threading.Thread): 
    """web服务存储器""" 
    def __init__(self,company_id,web_sevice_queue): 
        threading.Thread.__init__(self) 
        self.web_sevice_queue = web_sevice_queue 
        self.host_web_service_hash = HostWebServiceHash(company_id) 
    def run(self): 
        while True: 
            web_service_info = self.web_sevice_queue.get() 
            if web_service_info == "done": 
                break 
            self.do_work(web_service_info)
    def do_work(self,web_service_info): 
        self.host_web_service_hash.add(web_service_info['ip'],web_service_info['port'],web_service_info['url'],web_service_info['title'])