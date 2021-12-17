import request from '@/utils/request'

//添加任务
export function addTask(data) {
    return request({
        url:"task",
        method:"post",
        data
    })
}
//修改任务
export function editTask(data) {
    return request({
        url:"task",
        method:"put",
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
//开始任务
export function startTask(id){
    return request({
        url:`task/${id}/start`,
        method:"put",
    })
}

//停止任务
export function stopTask(id){
    return request({
        url:`task/${id}/stop`,
        method:"put",
    })
}