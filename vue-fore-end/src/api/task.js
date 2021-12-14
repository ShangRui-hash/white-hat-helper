import request from '@/utils/request'

//添加任务
export function addTask(data) {
    return request({
        url:"task",
        method:"post",
        data
    })
}

//获取任务列表
export function getTaskList(params) {
    return request({
        url:"task",
        method:"get",
        params
    })
}

//删除任务
export function deleteTask(id) {
    return request({
        url:"task",
        method:"delete",
        data:{
            id
        }
    })
}

export function startTask(id){
    return request({
        url:`task/${id}/start`,
        method:"put",
    })
}