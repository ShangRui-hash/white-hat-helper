import json

from django.http import HttpResponse

def resp_success(msg, data={}):
    """成功时返回的响应格式"""
    resp = {"status": "success", "msg": msg, "data": data}
    return HttpResponse(json.dumps(resp),
                        content_type='application/json;charset:utf-8',
                        status=200)


def resp_fail(msg):
    """失败时返回的响应格式"""
    resp = {'status': 'fail', 'message': msg}
    return HttpResponse(json.dumps(resp),
                        content_type='application/json;charset:utf-8',
                        status=400)
