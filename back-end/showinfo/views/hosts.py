from django.views.decorators.http import require_http_methods
from lib.redis.host import HostSet
from lib.redis.port import HostPortHash
from lib.utils.response import resp_fail, resp_success

@require_http_methods(["GET", "DELETE"])
def hosts_handler(request):
    if request.method == "GET":
        return get_hosts_handler(request)
    elif request.method == "DELTE":
        return del_hosts_handler(request)


def get_hosts_handler(request):
    """获取指定公司的所有主机"""
    #1.接收传参,后端校验
    company_id = request.GET.get("company_id")
    offset = request.GET.get("offset")
    count = request.GET.get("count")
    if not company_id or not company_id.isdigit():
        return resp_fail(msg="company_id 不能为空且必须为数字")
    if not offset or not offset.isdigit():
        return resp_fail("offset 不能为空")
    if not count or not count.isdigit():
        return resp_fail("count 不能为空")
    offset = int(offset)
    count = int(count)
    company_id = int(company_id)
    #2.业务逻辑
    results = {}
    #2.1 获取所有主机
    host_set = HostSet(company_id)
    hosts = host_set.smembers(offset, count)
    #2.2 获取主机的所有端口和服务信息
    host_port_hash = HostPortHash(company_id)
    for host in hosts: 
        results[host]=host_port_hash.get(host)
    #3.返回响应
    print(results)
    return resp_success("获取资产成功", results)
