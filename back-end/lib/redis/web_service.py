import json
from django_redis import get_redis_connection 

class HostWebServiceHash: 
    """存储host对应的web服务的相关信息""" 
    def __init__(self,company_id): 
        self.client =get_redis_connection("default")
        self.key = 'company::websevice::%d'%company_id 
    def get_sub_key(self,host_ip,port): 
        return '%s::%d'%(host_ip,port)
    def add(self,host_ip,port,url,title): 
        sub_key = self.get_sub_key(host_ip,port)
        data = {
            "url":url,
            "title":title
        }
        self.client.hset(self.key,sub_key,json.dumps(data))
    def get(self,host_ip,port): 
        sub_key = self.get_sub_key(host_ip,port)
        data = self.client.hget(self.key,sub_key)
        if data is None : 
            return {} 
        return json.loads(data)