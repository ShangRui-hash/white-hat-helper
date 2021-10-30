/*任务页面*/
import React, {Component} from 'react';
import {Button, Card,Table,Space,Progress} from "antd";

const columns = [
    {
        title: <h1 className='title-font'>任务名</h1>,
        dataIndex: 'missionName',
        key: 'missionName',
        render: text => <p className='content-font'>{text}</p>,
    },
    {
        title: <h1 className='title-font'>公司名</h1>,
        dataIndex: 'companyName',
        key: 'companyName',
        render: text => <p className='content-font'>{text}</p>,
    },
    {
        title: <h1 className='title-font'>进度</h1>,
        key: 'rate',
        dataIndex: 'rate',
        render: text => <Progress percent={text} status="active"/>,
    },
    {
        title: <h1 className='title-font'>操作</h1>,
        key: 'action',
        render: () => (
            <Space size="middle">
                <Button type='link' style={{color:'#FFBB00'}}>暂停</Button>
                <Button type='link' danger>删除</Button>
            </Space>
        ),
    },
];

const data = [
    {
        missionName: '任务 1',
        companyName: '腾讯',
        rate:80
    },
    {
        missionName: '任务 2',
        companyName: '腾讯',
        rate:80
    },
    {
        missionName: '任务 3',
        companyName: '腾讯',
        rate:80
    },
    {
        missionName: '任务 4',
        companyName: '腾讯',
        rate:80
    },
];

class MissionPage extends Component {
    render() {
        return (
            <Card
                className='data-card'
                title={<p className='title-font'>任务列表</p>}
                extra={<Button type='primary'>添加扫描任务</Button>}
                bordered={false}
            >
                <Table
                    columns={columns}
                    dataSource={data}
                />
            </Card>
        );
    }
}

export default MissionPage;
