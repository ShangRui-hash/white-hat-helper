import json
import time
import logging

from django.shortcuts import render
from django.views.decorators.http import require_http_methods
from scanhosts.models import Company
from lib.utils.response import resp_success, resp_fail

logger = logging.getLogger("django")


# Create your views here.
@require_http_methods(["POST", "GET", "DELETE", "PUT"])
def company_handler(request):
    if request.method == "GET":
        return get_company_handler(request)
    elif request.method == "POST":
        return add_company_handler(request)
    elif request.method == "DELETE":
        return del_company_handler(request)
    elif request.method == "PUT":
        return update_company_handler(request)


def get_company_handler(request):
    """获取所有公司"""
    #1.接收参数 count,offset
    count = int(request.GET.get('count'))
    offset = int(request.GET.get('offset'))
    #2.业务逻辑
    rows = Company.objects.all().order_by("created_at")[offset:count]
    companys = [{
        "id": i.id,
        "name": i.name,
        "created_at": i.created_at.strftime('%Y-%m-%d %H:%M:%S')
    } for i in rows]
    return resp_success(msg="获取公司列表成功", data=companys)


def add_company_handler(request):
    """添加公司"""
    #1.接收传参,后端效验
    name = request.POST.get('name')
    if not name:
        return resp_fail(msg="公司名称不能为空")
    name = name.strip()
    #2.业务逻辑
    # 2.1 查重
    exists = Company.objects.filter(name=name).exists()
    if exists:
        return resp_fail(msg="该公司已存在")
    # 2.2 添加新公司
    res = Company.objects.create(name=name)
    #3.返回响应
    return resp_success(msg="添加成功", data={"id": res.id})


def del_company_handler(request):
    """删除公司"""
    #1.接收传参,后端效验
    try:
        params = json.loads(request.body.decode("utf-8"))
    except:
        return resp_fail(msg="参数错误")
    id = params.get('id')
    if not id or not isinstance(id, int):
        return resp_fail(msg="id不能为空")
    #2.业务逻辑
    # 2.1 检查是否存在
    exists = Company.objects.filter(id=id).exists()
    if not exists:
        return resp_fail(msg="该公司不存在")
    # 2.2 删除
    Company.objects.filter(id=id).delete()
    #3.返回响应
    return resp_success(msg="删除成功")


def update_company_handler(request):
    """更新公司"""
    #1.接收传参,后端效验
    try:
        params = json.loads(request.body.decode("utf-8"))
    except:
        return resp_fail(msg="参数错误")
    id = params.get('id')
    name = params.get('name')
    if not id or not isinstance(id, int):
        return resp_fail(msg="id不能为空,且必须为数字类型")
    if not name or not isinstance(name, str):
        return resp_fail(msg="公司名称不能为空,且必须为字符串类型")
    name = name.strip()
    #2.业务逻辑
    # 2.1 检查是否存在
    exists = Company.objects.filter(id=id).exists()
    if not exists:
        return resp_fail(msg="该公司不存在")
    # 2.2 检查公司名是否重名
    exists = Company.objects.filter(name=name).exists()
    if exists:
        return resp_fail(msg="该公司名已被使用")
    # 2.2 更新
    Company.objects.filter(id=id).update(name=name)
    #3.返回响应
    return resp_success(msg="更新成功")