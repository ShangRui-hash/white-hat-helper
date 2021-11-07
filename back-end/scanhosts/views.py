import datetime

from django.views.decorators.http import require_http_methods
from scanhosts.models import Task, Company
from lib.utils.response import *
from lib.scanner import scanner


@require_http_methods(["GET", "POST", "PUT", "DELETE"])
def task_handler(request):
    if request.method == "GET":
        return get_task_list(request)
    elif request.method == "POST":
        return add_task(request)


def add_task(request):
    """添加任务"""
    #1.接收参数，后端校验
    try:
        params = json.loads(request.body.decode("utf-8"))
    except:
        return resp_fail("参数错误")
    targets = params.get('targets', "")
    company_id = params.get('company_id')
    if not targets:
        return resp_fail("targets不能为空")
    if not company_id or not isinstance(company_id, int):
        return resp_fail("company_id不能为空且必须为数字")
    if not Company.objects.filter(id=company_id).exists():
        return resp_fail("company_id不存在")
    #2.业务逻辑
    res = Task.objects.create(targets=",".join(targets), company_id=company_id)
    #3.返回响应
    return resp_success("添加成功", {"id": res.id})


def get_task_list(request):
    """获取任务列表"""
    #1.接收参数，后端校验
    offset = int(request.GET.get('offset', 0))
    count = int(request.GET.get('count', 10))
    #2.业务逻辑
    res = Task.objects.all()[offset:offset + count]
    tasks = [{
        "id": row.id,
        "company_id": row.company_id,
        "company": row.company.name,
        "targets": row.targets.split(","),
        "status": row.status,
        "created_at": row.created_at.strftime('%Y-%m-%d %H:%M:%S')
    } for row in res]
    data = {"count": len(tasks), "tasks": tasks}
    #3.返回响应
    return resp_success("获取成功", data)


@require_http_methods(["GET"])
def run_handler(request):
    """运行一个任务"""
    #1.接收参数,后端校验
    task_id = request.GET.get('task_id', "")
    if not task_id or not task_id.isdigit():
        return resp_fail("非法的任务id")
    if not Task.objects.filter(id=task_id).exists():
        return resp_fail("任务不存在")
    # if Task.objects.filter(id=task_id, status="running").exists():
    #     return resp_fail("该任务已经在运行中")
    #2.业务逻辑
    task = Task.objects.get(id=task_id)
    #2.1 启动任务
    scanner.run(task.company_id, task.targets.split(','))
    #2.2 修改任务状态
    Task.objects.filter(id=task_id).update(status="running")
    #3.返回响应
    return resp_success("启动成功")