import pexpect
import queue
import nmap3
import re

from lib.scanner.web_service_saver import WebServiceSaver
from lib.scanner.port_scanner import PortScanner 
from lib.scanner.port_saver import PortSaver 
from lib.scanner.subdomain_scanner import SubdomainScanner
from lib.scanner.host_ip_saver import HostIPSaver 

def run(company_id, domains=[]):
    print(company_id,domains)
    domain_queue = queue.Queue()
    hosts_queue = queue.Queue()
    unique_hosts_queue = queue.Queue()
    ports_queue = queue.Queue()
    web_service_queue = queue.Queue()
    for domain in domains:
        domain_queue.put(domain)
    #1.启动消费者: web服务存储器
    web_service_saver = WebServiceSaver(company_id,web_service_queue)
    web_service_saver.start()
    #1.启动消费者：端口和服务版本存储器
    port_saver = PortSaver(company_id,ports_queue)
    port_saver.start()
    #1.启动消费者: 端口和服务版本扫描器
    port_scanner_list =[]
    for i in range(15): 
        port_scanner = PortScanner(unique_hosts_queue, ports_queue)
        port_scanner.start()
        port_scanner_list.append(port_scanner)

    print("port_scanner 启动成功")
    #1.启动消费者：主机存储器
    host_ip_saver = HostIPSaver(company_id,hosts_queue,unique_hosts_queue)
    host_ip_saver.start()
    print("saver 启动成功")
    #2.启动生产者：子域名枚举器
    subdomain_scanner = SubdomainScanner(domain_queue,hosts_queue,web_service_queue)
    subdomain_scanner.start()

    # #fake data
    # for i in ['10.100.102.123', '123.123.123.123']:
    #     hosts_queue.put(i)


if __name__ == "__main__":
    domains = ["lenovo.net", "baiying.cn"]
    run(domains)
