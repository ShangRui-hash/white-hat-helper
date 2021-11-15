import threading
from lib.scanner.oneforall_scan import oneforall_scan


class SubdomainScanner(threading.Thread):
    def __init__(self, domain_queue, ip_queue, web_service_queue):
        threading.Thread.__init__(self)
        self.domain_queue = domain_queue
        self.ip_queue = ip_queue
        self.web_service_queue = web_service_queue

    def run(self):
        while True:
            item = self.domain_queue.get()
            self.do_work(item)
            if self.domain_queue.empty():
                break
        self.ip_queue.put("done")

    def do_work(self, item):
        #1.查询子域名
        results = oneforall_scan(item)
        print("开始遍历oneforall扫描结果",len(results))
        for result in results:
            #2.主机信息入管道
            if result['ip'].find(","):
                port = result['port']
                url = result['url']
                title = result['title']
                for ip in result['ip'].split(","):
                    self.ip_queue.put(ip)
                    web_service_item = {
                        "ip": ip,
                        "port": port,
                        "url": url,
                        "title": title
                    }
                    # print(web_service_item)
                    self.web_service_queue.put(web_service_item)
            else:
                ip = result['ip']
                self.ip_queue.put(ip)
                #3.web服务信息入管道
                web_service_item = {
                    "ip": ip,
                    "port": result['port'],
                    "url": result['url'],
                    "title": result['title']
                }
                # print(web_service_item)
                self.web_service_queue.put(web_service_item)
