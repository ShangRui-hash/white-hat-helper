import request from '@/utils/request' 

//获取公司列表
export function getCompanyList(offset,count) {
  return request({
    url: `company?count=${count}&offset=${offset}`,
    method: 'get'
  })
}

//添加公司
export function addCompany(data) {
  return request({
    url: `company`,
    method: 'post',
    data
  })
}

//删除公司
export function deleteCompany(id) {
  return request({
    url: `company`,
    method: 'delete',
    data:{
      id
    }
  })
}

//修改公司信息
export function updateCompany(data) {
  return request({
    url: `company`,
    method: 'put',
    data
  })
}