import threading
from lib.redis.host import HostSet

class HostIPSaver(threading.Thread):
    """存储IP到redis中，并去重输出"""
    def __init__(self,company_id, inq,outq):
        threading.Thread.__init__(self)
        self._inq = inq
        self._outq = outq
        self.host_set = HostSet(company_id)
        self._set = set()

    def run(self):
        while True:
            host_ip = self._inq.get()
            if host_ip == "done": 
                break
            self.do_work(host_ip)
        self._outq.put("done")

    def do_work(self, host_ip):
        #入库
        self.host_set.add(host_ip)
        #去重输出
        if host_ip in self._set: 
            return 
        self._set.add(host_ip)
        self._outq.put(host_ip)