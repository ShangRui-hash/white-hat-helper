import json
from django_redis import get_redis_connection 

class HostPortHash: 
    """存储host的port和服务信息的hash列表"""
    def __init__(self,company_id): 
        self.client=get_redis_connection("default")
        self.key='company::port::%d'%company_id
    def add(self,host_ip,port_service_info):
        self.client.hset(self.key,host_ip,json.dumps(port_service_info))
    def get(self,host_ip): 
        info = self.client.hget(self.key,host_ip)
        if info is None: 
            return {}
        return json.loads(info.decode("utf-8")) 
