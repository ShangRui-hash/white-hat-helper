/*公司页面*/
import React, {Component} from 'react';
import { Card, Button, Table, Space} from "antd";
import {Link} from "react-router-dom";

const columns = [
    {
        title: <h1 className='title-font'>公司名</h1>,
        dataIndex: 'companyName',
        key: 'companyName',
        render: text =>(
            <Link
                to={'/home/CompanyPage/AssetsPage'}
            >
                <Button className='my-font' type='link'>{text}</Button>
            </Link>
        ),
        width:'50%'
    },
    {
        title: <h1 className='title-font'>操作</h1>,
        key: 'action',
        render: () => (
            <Space size="middle">
                <Button type='link'>修改</Button>
                <Button type='link' danger>删除</Button>
            </Space>
        ),
    },
];


const data = [
    {
        companyName: '公司 1',
    },
    {
        companyName: '公司 2',
    },
    {
        companyName: '公司 3',
    },
    {
        companyName: '公司 4',
    },
];

class CompanyPage extends Component {
    render() {
        return (
            <Card
                className='data-card'
                title={<p className='title-font'>公司列表</p>}
                extra={<Button type='primary'>添加</Button>}
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

export default CompanyPage;
