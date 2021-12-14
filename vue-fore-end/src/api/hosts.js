import request from '@/utils/request'

export function getHostList(params) {
  return request({
    url: '/host/list',
    method: 'get',
    params
  })
}

export function getHostDetail(params) {
  return request({
    url: '/host/detail',
    method: 'get',
    params
  })
}
