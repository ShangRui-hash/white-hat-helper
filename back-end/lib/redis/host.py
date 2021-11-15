from django_redis import get_redis_connection 

class HostSet: 
    def __init__(self,company_id): 
        self.client=get_redis_connection("default")
        self.key='company::%d'%company_id
    def add(self,host_ip):
        self.client.sadd(self.key,host_ip)
    def smembers(self,offset=0,count=10): 
        hosts_set = [ i.decode("utf-8") for i in list(self.client.smembers(self.key))[offset:offset+count] ] 
        return hosts_set
    def is_member(self,host_ip):
        return self.client.sismember(self.key,host_ip)
        
if __name__ == "__main__": 
    hs=HostSet(1)
    hs.add("123.123.123.123")