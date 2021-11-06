/*资产页面*/
import React, {Component} from 'react';
import {Button, Card, Table} from "antd";
import {Link} from "react-router-dom";

const columns = [
    {
        title: <h1 className='title-font'>IP</h1>,
        dataIndex: 'IP',
        key: 'IP',
        render: text =>(
            <Link
                to={'/home/CompanyPage/AssetsPage/details'}
            >
                <Button className='my-font' type='link'>{text}</Button>
            </Link>
        ),
    },
    {
        title: <h1 className='title-font'>Port</h1>,
        dataIndex: 'Port',
        key: 'Port',
        render: text => <p className='content-font'>{text}</p>,
    },
    {
        title: <h1 className='title-font'>Protocol</h1>,
        dataIndex: 'Protocol',
        key: 'Protocol',
        render: text => <p className='content-font'>{text}</p>,
    },
    {
        title: <h1 className='title-font'>Component</h1>,
        key: 'Component',
        dataIndex: 'Component',
        render: text => <p className='content-font'>{text}</p>,
    },
];

const data = [
    {
        key: '1',
        IP: 'John Brown',
        Port: 32,
        Protocol: 'New York No. 1 Lake Park',
        Component: ['nice', 'developer'],
    },
    {
        key: '2',
        IP: 'Jim Green',
        Port: 42,
        Protocol: 'London No. 1 Lake Park',
        Component: ['loser'],
    },
    {
        key: '3',
        IP: 'Joe Black',
        Port: 32,
        Protocol: 'Sidney No. 1 Lake Park',
        Component: ['cool', 'teacher'],
    },
];

class AssetsPage extends Component {
    render() {
        return (
            <Card
                className='data-card'
                title={<p className='title-font'>公司资产</p>}
                bordered={false}
            >
                <Table columns={columns} dataSource={data}/>
            </Card>
        );
    }
}

export default AssetsPage;
